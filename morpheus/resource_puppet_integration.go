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

func resourcePuppetIntegration() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides an puppet integration resource",
		CreateContext: resourcePuppetIntegrationCreate,
		ReadContext:   resourcePuppetIntegrationRead,
		UpdateContext: resourcePuppetIntegrationUpdate,
		DeleteContext: resourcePuppetIntegrationDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the puppet integration",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the puppet integration",
				Required:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the puppet integration is enabled",
				Optional:    true,
				Computed:    true,
			},
			"puppet_master_hostname": {
				Type:        schema.TypeString,
				Description: "The hostname of the puppet server",
				Required:    true,
			},
			"allow_immediate_execution": {
				Type:        schema.TypeBool,
				Description: "Whether to trigger the immediate execution of a puppet agent run",
				Optional:    true,
				Computed:    true,
			},
			"puppet_master_ssh_username": {
				Type:        schema.TypeString,
				Description: "The username of the account on the puppet server used to trigger the immediate execution of a puppet agent run",
				Optional:    true,
				Computed:    true,
			},
			"puppet_master_ssh_password": {
				Type:        schema.TypeString,
				Description: "The password of the account on the puppet server used to trigger the immediate execution of a puppet agent run",
				Optional:    true,
				Computed:    true,
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

func resourcePuppetIntegrationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	integration := make(map[string]interface{})

	integration["name"] = d.Get("name").(string)
	integration["enabled"] = d.Get("enabled").(bool)
	integration["type"] = "puppet"

	config := make(map[string]interface{})
	config["puppetMaster"] = d.Get("puppet_master_hostname").(string)
	if d.Get("allow_immediate_execution").(bool) {
		config["puppetFireNow"] = "true"
	} else {
		config["puppetFireNow"] = "false"
	}
	config["puppetSshUser"] = d.Get("puppet_master_ssh_username").(string)
	config["puppetSshPassword"] = d.Get("puppet_master_ssh_password").(string)

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

	resourcePuppetIntegrationRead(ctx, d, meta)
	return diags
}

func resourcePuppetIntegrationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	d.Set("puppet_master_hostname", integration.Config.PuppetMaster)
	if integration.Config.PuppetFireNow == "true" {
		d.Set("allow_immediate_execution", true)
	} else {
		d.Set("allow_immediate_execution", false)
	}
	d.Set("puppet_master_ssh_username", integration.Config.PuppetSshUser)
	d.Set("puppet_master_ssh_password", integration.Config.PuppetSshPasswordHash)

	return diags
}

func resourcePuppetIntegrationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	integration := make(map[string]interface{})

	integration["name"] = d.Get("name").(string)
	integration["enabled"] = d.Get("enabled").(bool)
	integration["type"] = "puppet"

	config := make(map[string]interface{})
	config["puppetMaster"] = d.Get("puppet_master_hostname").(string)
	if d.Get("allow_immediate_execution").(bool) {
		config["puppetFireNow"] = "true"
	} else {
		config["puppetFireNow"] = "false"
	}
	config["puppetSshUser"] = d.Get("puppet_master_ssh_username").(string)
	config["puppetSshPassword"] = d.Get("puppet_master_ssh_password").(string)

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
	return resourcePuppetIntegrationRead(ctx, d, meta)
}

func resourcePuppetIntegrationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
