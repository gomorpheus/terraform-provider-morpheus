package morpheus

import (
	"context"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceTagPolicy() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus tag policy resource",
		CreateContext: resourceTagPolicyCreate,
		ReadContext:   resourceTagPolicyRead,
		UpdateContext: resourceTagPolicyUpdate,
		DeleteContext: resourceTagPolicyDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the tag policy",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the tag policy",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the tag policy",
				Optional:    true,
				Computed:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the policy is enabled",
				Optional:    true,
				Default:     true,
			},
			"strict_enforcement": {
				Type:        schema.TypeBool,
				Description: "Whether users will be able to provision new workloads if they violate the tag policy",
				Optional:    true,
				Computed:    true,
			},
			"tag_key": {
				Type:        schema.TypeString,
				Description: "The key of the tag to enforce",
				Required:    true,
			},
			"tag_value": {
				Type:        schema.TypeString,
				Description: "The value of the tag to enforce",
				Optional:    true,
				Computed:    true,
			},
			"option_list_id": {
				Type:        schema.TypeInt,
				Description: "The id of the option list associated with the policy",
				Optional:    true,
				Computed:    true,
			},
			"scope": {
				Type:         schema.TypeString,
				Description:  "The filter or scope that the policy is applied to (global, group, cloud, user)",
				ValidateFunc: validation.StringInSlice([]string{"global", "group", "cloud", "user"}, false),
				Required:     true,
				ForceNew:     true,
			},
			"group_id": {
				Type:          schema.TypeInt,
				Description:   "The id of the group associated with the group scoped filter",
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"cloud_id", "user_id"},
			},
			"cloud_id": {
				Type:          schema.TypeInt,
				Description:   "The id of the cloud associated with the cloud scoped filter",
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"group_id", "user_id"},
			},
			"user_id": {
				Type:          schema.TypeInt,
				Description:   "The id of the user associated with the user scoped filter",
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"cloud_id", "group_id"},
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

func resourceTagPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	policy := make(map[string]interface{})

	policy["name"] = d.Get("name").(string)
	policy["description"] = d.Get("description").(string)
	policy["enabled"] = d.Get("enabled").(bool)

	policy["config"] = map[string]interface{}{
		"strict":      d.Get("strict_enforcement").(bool),
		"key":         d.Get("tag_key").(string),
		"valueListId": d.Get("option_list_id").(int),
		"value":       d.Get("tag_value").(string),
	}
	policy["policyType"] = map[string]interface{}{
		"code": "tags",
		"name": "Tags",
	}
	policy["accounts"] = d.Get("tenant_ids")

	switch d.Get("scope") {
	case "group":
		policy["refId"] = d.Get("group_id").(int)
		policy["refType"] = "ComputeSite"
		policy["site"] = map[string]interface{}{
			"id": d.Get("group_id").(int),
		}
	case "cloud":
		policy["refId"] = d.Get("cloud_id").(int)
		policy["refType"] = "ComputeZone"
		policy["zone"] = map[string]interface{}{
			"id": d.Get("cloud_id").(int),
		}
	case "user":
		policy["refId"] = d.Get("user_id").(int)
		policy["refType"] = "User"
		policy["user"] = map[string]interface{}{
			"id": d.Get("user_id").(int),
		}
	}

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

	resourceTagPolicyRead(ctx, d, meta)
	return diags
}

func resourceTagPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	tagPolicy := result.Policy

	d.SetId(int64ToString(tagPolicy.ID))
	d.Set("name", tagPolicy.Name)
	d.Set("description", tagPolicy.Description)
	d.Set("enabled", tagPolicy.Enabled)
	d.Set("strict_enforcement", tagPolicy.Config.Strict)
	d.Set("tag_key", tagPolicy.Config.Key)
	d.Set("tag_value", tagPolicy.Config.Value)
	d.Set("option_list_id", tagPolicy.Config.ValueListId)

	switch tagPolicy.RefType {
	case "ComputeSite":
		d.Set("scope", "group")
		d.Set("group_id", tagPolicy.Site.ID)
	case "ComputeZone":
		d.Set("scope", "cloud")
		d.Set("cloud_id", tagPolicy.Zone.ID)
	case "User":
		d.Set("scope", "user")
		d.Set("user_id", tagPolicy.User.ID)
	default:
		d.Set("scope", "global")
	}

	var tenantIds []int64
	if tagPolicy.Accounts != nil {
		// iterate over the array of accounts
		for _, account := range tagPolicy.Accounts {
			tenantIds = append(tenantIds, account.ID)
		}
	}
	d.Set("tenant_ids", tenantIds)

	return diags
}

func resourceTagPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	policy := make(map[string]interface{})

	policy["name"] = d.Get("name").(string)
	policy["description"] = d.Get("description").(string)
	policy["enabled"] = d.Get("enabled").(bool)

	policy["config"] = map[string]interface{}{
		"strict":      d.Get("strict_enforcement").(bool),
		"key":         d.Get("tag_key").(string),
		"valueListId": d.Get("option_list_id").(int),
		"value":       d.Get("tag_value").(string),
	}
	policy["policyType"] = map[string]interface{}{
		"code": "tags",
		"name": "Tags",
	}
	policy["accounts"] = d.Get("tenant_ids")

	switch d.Get("scope") {
	case "group":
		policy["refId"] = d.Get("group_id").(int)
		policy["refType"] = "ComputeSite"
		policy["site"] = map[string]interface{}{
			"id": d.Get("group_id").(int),
		}
	case "cloud":
		policy["refId"] = d.Get("cloud_id").(int)
		policy["refType"] = "ComputeZone"
		policy["zone"] = map[string]interface{}{
			"id": d.Get("cloud_id").(int),
		}
	case "user":
		policy["refId"] = d.Get("user_id").(int)
		policy["refType"] = "User"
		policy["user"] = map[string]interface{}{
			"id": d.Get("user_id").(int),
		}
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"policy": policy,
		},
	}

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
	return resourceTagPolicyRead(ctx, d, meta)
}

func resourceTagPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
