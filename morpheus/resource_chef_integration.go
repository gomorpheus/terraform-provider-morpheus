package morpheus

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceChefIntegration() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Chef integration resource",
		CreateContext: resourceChefIntegrationCreate,
		ReadContext:   resourceChefIntegrationRead,
		UpdateContext: resourceChefIntegrationUpdate,
		DeleteContext: resourceChefIntegrationDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the Chef integration",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Chef integration",
				Required:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the Chef integration is enabled",
				Optional:    true,
				Computed:    true,
			},
			"url": {
				Type:        schema.TypeString,
				Description: "The url of the Chef server",
				Required:    true,
			},
			"version": {
				Type:        schema.TypeString,
				Description: "The version of the Chef server",
				Optional:    true,
			},
			"windows_version": {
				Type:        schema.TypeString,
				Description: "The Windows agent version",
				Optional:    true,
			},
			"windows_msi_install_url": {
				Type:        schema.TypeString,
				Description: "The URL for the Windows MSI installation package",
				Optional:    true,
			},
			"organization": {
				Type:        schema.TypeString,
				Description: "The chef organization",
				Optional:    true,
			},
			"use_fqdn_node_name": {
				Type:        schema.TypeBool,
				Description: "Whether to use the FQDN of the node instead of the instance name",
				Optional:    true,
				Default:     false,
			},
			"username": {
				Type:          schema.TypeString,
				Description:   "The username of the account used to connect to the Chef server",
				Optional:      true,
				ConflictsWith: []string{"credential_id"},
			},
			"private_key": {
				Type:        schema.TypeString,
				Description: "The private key of the account used to connect to the Chef server",
				Optional:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256_hash := hex.EncodeToString(h.Sum(nil))
					return strings.EqualFold(old, sha256_hash)
				},
				ConflictsWith: []string{"credential_id"},
			},
			"credential_id": {
				Description:   "The ID of the credential store entry used for authentication",
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"username", "private_key"},
			},
			"organization_validator_key": {
				Type:        schema.TypeString,
				Description: "The organization validator key used to connect to the Chef server",
				Optional:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256_hash := hex.EncodeToString(h.Sum(nil))
					return strings.EqualFold(old, sha256_hash)
				},
			},
			/* AWAITING API SUPPORT
			"databags": {
				Description: "",
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			*/
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceChefIntegrationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	integration := make(map[string]interface{})

	integration["name"] = d.Get("name").(string)
	integration["enabled"] = d.Get("enabled").(bool)
	integration["type"] = "chef"
	integration["serviceUrl"] = d.Get("url").(string)
	integration["serviceVersion"] = d.Get("version").(string)
	integration["serviceWindowsVersion"] = d.Get("windows_version").(string)

	config := make(map[string]interface{})
	config["windowsInstallUrl"] = d.Get("windows_msi_install_url").(string)
	config["useFqdn"] = d.Get("use_fqdn_node_name").(bool)
	config["org"] = d.Get("organization").(string)

	if d.Get("credential_id").(int) != 0 {
		credential := make(map[string]interface{})
		credential["type"] = "username-keypair"
		credential["id"] = d.Get("credential_id").(int)
		credential["credential"] = credential
	} else {
		credential := make(map[string]interface{})
		credential["type"] = "local"
		integration["credential"] = credential
		config["chefUser"] = d.Get("username").(string)
		config["userKey"] = d.Get("private_key").(string)
	}
	config["orgKey"] = d.Get("organization_validator_key").(string)

	// databags
	/* AWAITING API SUPPORT
	if d.Get("databags") != nil {
		databagsInput := d.Get("databags").(map[string]interface{})
		var databags []map[string]interface{}
		for key, value := range databagsInput {
			databag := make(map[string]interface{})
			databag["name"] = key
			databag["value"] = value.(string)
			databags = append(databags, databag)
		}
		config["databag"] = databags
	}
	*/
	integration["config"] = config

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"integration": integration,
		},
	}

	resp, err := client.CreateIntegration(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateIntegrationResult)
	integrationResult := result.Integration
	// Successfully created resource, now set id
	d.SetId(int64ToString(integrationResult.ID))

	resourceChefIntegrationRead(ctx, d, meta)
	return diags
}

func resourceChefIntegrationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindIntegrationByName(name)
	} else if id != "" {
		resp, err = client.GetIntegration(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Integration cannot be read without name or id")
	}

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
			log.Printf("Forcing recreation of resource")
			d.SetId("")
			return diags
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)
			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)

	// store resource data
	result := resp.Result.(*morpheus.GetIntegrationResult)
	integration := result.Integration
	d.SetId(int64ToString(integration.ID))
	d.Set("name", integration.Name)
	d.Set("enabled", integration.Enabled)
	d.Set("url", integration.URL)
	d.Set("version", integration.Version)
	d.Set("windows_version", integration.WindowsVersion)
	d.Set("windows_msi_install_url", integration.Config.WindowsInstallURL)
	d.Set("organization", integration.Config.Org)
	d.Set("use_fqdn_node_name", integration.Config.UseFqdn)
	if integration.Credential.ID == 0 {
		d.Set("username", integration.Config.ChefUser)
		d.Set("private_key", integration.Config.UserKeyHash)
	} else {
		d.Set("credential_id", integration.Credential.ID)
	}
	d.Set("organization_validator_key", integration.Config.OrgKeyHash)

	// databags
	/* AWAITING API SUPPORT
	databags := make(map[string]interface{})
	if integration.Config.Databags != nil {
		output := integration.Config.Databags
		databagList := output
		// iterate over the array of databags
		for i := 0; i < len(databagList); i++ {
			databag := databagList[i]
			databagName := databag.Path
			databags[databagName] = databag.Key
		}
	}
	d.Set("databags", databags)
	*/
	return diags
}

func resourceChefIntegrationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	integration := make(map[string]interface{})

	integration["name"] = d.Get("name").(string)
	integration["enabled"] = d.Get("enabled").(bool)
	integration["type"] = "chef"
	integration["serviceUrl"] = d.Get("url").(string)
	integration["serviceVersion"] = d.Get("version").(string)
	integration["serviceWindowsVersion"] = d.Get("windows_version").(string)

	config := make(map[string]interface{})
	config["windowsInstallUrl"] = d.Get("windows_msi_install_url").(string)
	config["useFqdn"] = d.Get("use_fqdn_node_name").(bool)
	config["org"] = d.Get("organization").(string)

	if d.Get("credential_id").(int) != 0 {
		credential := make(map[string]interface{})
		credential["type"] = "username-keypair"
		credential["id"] = d.Get("credential_id").(int)
		integration["credential"] = credential
	} else {
		credential := make(map[string]interface{})
		credential["type"] = "local"
		integration["credential"] = credential
		if d.HasChange("username") {
			config["chefUser"] = d.Get("username").(string)
		}
		if d.HasChange("private_key") {
			config["userKey"] = d.Get("private_key").(string)
		}
	}

	if d.HasChange("organization_validator_key") {
		config["orgKey"] = d.Get("organization_validator_key").(string)
	}

	integration["config"] = config

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"integration": integration,
		},
	}

	resp, err := client.UpdateIntegration(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateIntegrationResult)
	integrationResult := result.Integration

	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(integrationResult.ID))
	return resourceChefIntegrationRead(ctx, d, meta)
}

func resourceChefIntegrationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteIntegration(toInt64(id), req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
			return diag.FromErr(err)
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)
			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)
	d.SetId("")
	return diags
}
