package morpheus

import (
	"context"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceMotdPolicy() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus message of the day policy resource",
		CreateContext: resourceMotdPolicyCreate,
		ReadContext:   resourceMotdPolicyRead,
		UpdateContext: resourceMotdPolicyUpdate,
		DeleteContext: resourceMotdPolicyDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the message of the day policy",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the message of the day policy",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the message of the day policy",
				Optional:    true,
				Computed:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the policy is enabled",
				Optional:    true,
				Default:     true,
			},
			"title": {
				Type:        schema.TypeString,
				Description: "The title of the message of the day",
				Optional:    true,
				Computed:    true,
			},
			"message": {
				Type:        schema.TypeString,
				Description: "The message of the message of the day",
				Required:    true,
			},
			"type": {
				Type:         schema.TypeString,
				Description:  "The message type of the message of the day (info, warning, critical)",
				ValidateFunc: validation.StringInSlice([]string{"info", "warning", "critical"}, false),
				Required:     true,
			},
			"full_page": {
				Type:        schema.TypeString,
				Description: "Whether the message of the day is displayed as a full page or just a notification dialog box",
				Optional:    true,
				Computed:    true,
			},
			"tenant_ids": {
				Type:        schema.TypeList,
				Description: "A list of tenant IDs to assign the policy to",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceMotdPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	policy := make(map[string]interface{})

	policy["name"] = d.Get("name").(string)
	policy["description"] = d.Get("description").(string)
	policy["enabled"] = d.Get("enabled").(bool)

	policy["config"] = map[string]interface{}{
		"motd.title":    d.Get("title").(string),
		"motd.message":  d.Get("message").(string),
		"motd.type":     d.Get("type").(string),
		"motd.fullPage": d.Get("full_page").(string),
	}
	policy["policyType"] = map[string]interface{}{
		"code": "motd",
		"name": "Message of the Day",
	}

	policy["accounts"] = d.Get("tenant_ids")

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"policy": policy,
		},
	}

	resp, err := client.CreatePolicy(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreatePolicyResult)
	policyResult := result.Policy
	// Successfully created resource, now set id
	d.SetId(int64ToString(policyResult.ID))

	resourceMotdPolicyRead(ctx, d, meta)
	return diags
}

func resourceMotdPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindPolicyByName(name)
	} else if id != "" {
		resp, err = client.GetPolicy(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Policy cannot be read without name or id")
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
	result := resp.Result.(*morpheus.GetPolicyResult)
	MotdPolicy := result.Policy

	d.SetId(int64ToString(MotdPolicy.ID))
	d.Set("name", MotdPolicy.Name)
	d.Set("description", MotdPolicy.Description)
	if MotdPolicy.Enabled {
		d.Set("enabled", true)
	} else {
		d.Set("enabled", false)
	}
	d.Set("title", MotdPolicy.Config.MotdTitle)
	d.Set("full_page", MotdPolicy.Config.MotdFullPage)
	d.Set("type", MotdPolicy.Config.MotdType)
	d.Set("message", MotdPolicy.Config.MotdMessage)
	var tenantIds []int64
	if MotdPolicy.Accounts != nil {
		// iterate over the array of accounts
		for _, account := range MotdPolicy.Accounts {
			tenantIds = append(tenantIds, account.ID)
		}
	}
	d.Set("tenant_ids", tenantIds)

	return diags
}

func resourceMotdPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	policy := make(map[string]interface{})

	policy["name"] = d.Get("name").(string)
	policy["description"] = d.Get("description").(string)
	policy["enabled"] = d.Get("enabled").(bool)

	policy["config"] = map[string]interface{}{
		"motd.title":    d.Get("title").(string),
		"motd.message":  d.Get("message").(string),
		"motd.type":     d.Get("type").(string),
		"motd.fullPage": d.Get("full_page").(string),
	}
	policy["policyType"] = map[string]interface{}{
		"code": "motd",
		"name": "Message of the Day",
	}

	policy["accounts"] = d.Get("tenant_ids")

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"policy": policy,
		},
	}

	log.Printf("API REQUEST: %s", req)
	resp, err := client.UpdatePolicy(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdatePolicyResult)
	policyResult := result.Policy

	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(policyResult.ID))
	return resourceMotdPolicyRead(ctx, d, meta)
}

func resourceMotdPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeletePolicy(toInt64(id), req)
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
