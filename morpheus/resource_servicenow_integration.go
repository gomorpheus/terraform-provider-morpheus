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

func resourceServiceNowIntegration() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a ServiceNow integration resource",
		CreateContext: resourceServiceNowIntegrationCreate,
		ReadContext:   resourceServiceNowIntegrationRead,
		UpdateContext: resourceServiceNowIntegrationUpdate,
		DeleteContext: resourceServiceNowIntegrationDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the ServiceNow integration",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the ServiceNow integration",
				Required:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the SerivceNow integration is enabled",
				Optional:    true,
				Computed:    true,
			},
			"url": {
				Type:        schema.TypeString,
				Description: "The url of the ServiceNow instance",
				Required:    true,
			},
			"username": {
				Type:        schema.TypeString,
				Description: "The username of the account used to connect to ServiceNow",
				Required:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "The password of the account used to connect to ServiceNow",
				Required:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256_hash := hex.EncodeToString(h.Sum(nil))
					return strings.EqualFold(old, sha256_hash)
				},
				DiffSuppressOnRefresh: true,
			},
			"cmdb_custom_mapping": {
				Type:        schema.TypeString,
				Description: "The username of the account used to connect to ServiceNow",
				Optional:    true,
			},
			"default_cmdb_business_class": {
				Type:        schema.TypeString,
				Description: "The username of the account used to connect to ServiceNow",
				Optional:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceServiceNowIntegrationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	integration := make(map[string]interface{})

	integration["name"] = d.Get("name").(string)
	integration["enabled"] = d.Get("enabled").(bool)
	integration["type"] = "ansibleTower"
	integration["version"] = "v2"
	integration["url"] = d.Get("url").(string)
	integration["username"] = d.Get("username").(string)
	integration["password"] = d.Get("password").(string)

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

	resourceServiceNowIntegrationRead(ctx, d, meta)
	return diags
}

func resourceServiceNowIntegrationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
			return diag.FromErr(err)
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
	d.Set("username", integration.Username)
	d.Set("password", integration.PasswordHash)
	integration.Config.ServiceNowCmdbClassMapping
	d.Set("cmdb_custom_mapping", integration.Config.ServiceNowCustomCmdbMapping)
	d.Set("default_cmdb_business_class", integration.Config.ServiceNowCMDBBusinessObject)
	return diags
}

func resourceServiceNowIntegrationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	integration := make(map[string]interface{})

	integration["name"] = d.Get("name").(string)
	integration["enabled"] = d.Get("enabled").(bool)
	integration["type"] = "ansibleTower"
	integration["version"] = "v2"
	integration["url"] = d.Get("url").(string)
	integration["username"] = d.Get("username").(string)
	integration["password"] = d.Get("password").(string)

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
	return resourceServiceNowIntegrationRead(ctx, d, meta)
}

func resourceServiceNowIntegrationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
