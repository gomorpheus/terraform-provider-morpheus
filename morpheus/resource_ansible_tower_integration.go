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

func resourceAnsibleTowerIntegration() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides an Ansible Tower integration resource",
		CreateContext: resourceAnsibleTowerIntegrationCreate,
		ReadContext:   resourceAnsibleTowerIntegrationRead,
		UpdateContext: resourceAnsibleTowerIntegrationUpdate,
		DeleteContext: resourceAnsibleTowerIntegrationDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the Ansible Tower integration",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Ansible Tower integration",
				Required:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the Ansible Tower integration is enabled",
				Optional:    true,
				Computed:    true,
			},
			"url": {
				Type:        schema.TypeString,
				Description: "The url of the Ansible Tower instance",
				Required:    true,
			},
			"username": {
				Type:        schema.TypeString,
				Description: "The username of the account used to connect to Ansible Tower",
				Required:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "The password of the account used to connect to Ansible Tower",
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
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceAnsibleTowerIntegrationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	resourceAnsibleTowerIntegrationRead(ctx, d, meta)
	return diags
}

func resourceAnsibleTowerIntegrationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	return diags
}

func resourceAnsibleTowerIntegrationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	return resourceAnsibleTowerIntegrationRead(ctx, d, meta)
}

func resourceAnsibleTowerIntegrationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
