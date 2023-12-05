package morpheus

import (
	"context"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceCypherAccessPolicy() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus cypher access policy resource",
		CreateContext: resourceCypherAccessPolicyCreate,
		ReadContext:   resourceCypherAccessPolicyRead,
		UpdateContext: resourceCypherAccessPolicyUpdate,
		DeleteContext: resourceCypherAccessPolicyDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the cypher access policy",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the cypher access policy",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the cypher access policy",
				Optional:    true,
				Computed:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the policy is enabled",
				Optional:    true,
				Default:     true,
			},
			"key_path": {
				Type:        schema.TypeString,
				Description: "The key path associated with the cypher access policy",
				Required:    true,
			},
			"read_access": {
				Type:        schema.TypeBool,
				Description: "Whether the policy grants read access",
				Optional:    true,
				Computed:    true,
			},
			"write_access": {
				Type:        schema.TypeBool,
				Description: "Whether the policy grants write access",
				Optional:    true,
				Computed:    true,
			},
			"update_access": {
				Type:        schema.TypeBool,
				Description: "Whether the policy grants update access",
				Optional:    true,
				Computed:    true,
			},
			"delete_access": {
				Type:        schema.TypeBool,
				Description: "Whether the policy grants delete access",
				Optional:    true,
				Computed:    true,
			},
			"list_access": {
				Type:        schema.TypeBool,
				Description: "Whether the policy grants list access",
				Optional:    true,
				Computed:    true,
			},
			"scope": {
				Type:         schema.TypeString,
				Description:  "The filter or scope that the policy is applied to (global, user, role)",
				ValidateFunc: validation.StringInSlice([]string{"global", "user", "role"}, false),
				Required:     true,
				ForceNew:     true,
			},
			"user_id": {
				Type:          schema.TypeInt,
				Description:   "The id of the user associated with the user scoped filter",
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"role_id"},
			},
			"role_id": {
				Type:          schema.TypeInt,
				Description:   "The id of the role associated with the role scoped filter",
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"user_id"},
			},
			"apply_to_each_user": {
				Type:          schema.TypeBool,
				Description:   "Whether to assign the policy at the individual user level to all users assigned the associated role",
				Optional:      true,
				ConflictsWith: []string{"user_id"},
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

func resourceCypherAccessPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	policy := make(map[string]interface{})

	policy["name"] = d.Get("name").(string)
	policy["description"] = d.Get("description").(string)
	policy["enabled"] = d.Get("enabled").(bool)
	policy["config"] = map[string]interface{}{
		"keyPattern": d.Get("key_path").(string),
		"read":       evaluateStringBoolean(d.Get("read_access").(bool)),
		"write":      evaluateStringBoolean(d.Get("write_access").(bool)),
		"update":     evaluateStringBoolean(d.Get("update_access").(bool)),
		"delete":     evaluateStringBoolean(d.Get("delete_access").(bool)),
		"list":       evaluateStringBoolean(d.Get("list_access").(bool)),
	}
	policy["policyType"] = map[string]interface{}{
		"code": "cypher",
		"name": "Cypher Access",
	}

	policy["accounts"] = d.Get("tenant_ids")

	switch d.Get("scope") {
	case "user":
		policy["refId"] = d.Get("user_id").(int)
		policy["refType"] = "User"
		policy["user"] = map[string]interface{}{
			"id": d.Get("user_id").(int),
		}
	case "role":
		policy["refId"] = d.Get("role_id").(int)
		policy["refType"] = "Role"
		policy["eachUser"] = d.Get("apply_to_each_user").(bool)
		policy["role"] = map[string]interface{}{
			"id": d.Get("role_id").(int),
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

	resourceCypherAccessPolicyRead(ctx, d, meta)
	return diags
}

func resourceCypherAccessPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	result := resp.Result.(*morpheus.GetPolicyResult)
	cypherAccessPolicy := result.Policy

	d.SetId(int64ToString(cypherAccessPolicy.ID))
	d.Set("name", cypherAccessPolicy.Name)
	d.Set("description", cypherAccessPolicy.Description)
	d.Set("enabled", cypherAccessPolicy.Enabled)
	d.Set("key_path", cypherAccessPolicy.Config.KeyPattern)
	d.Set("read_access", parseStringBoolean(cypherAccessPolicy.Config.Read))
	d.Set("write_access", parseStringBoolean(cypherAccessPolicy.Config.Write))
	d.Set("update_access", parseStringBoolean(cypherAccessPolicy.Config.Update))
	d.Set("delete_access", parseStringBoolean(cypherAccessPolicy.Config.Delete))
	d.Set("list_access", parseStringBoolean(cypherAccessPolicy.Config.List))

	switch cypherAccessPolicy.RefType {
	case "User":
		d.Set("scope", "user")
		d.Set("user_id", cypherAccessPolicy.User.ID)
	case "Role":
		d.Set("scope", "role")
		d.Set("role_id", cypherAccessPolicy.Role.ID)
		d.Set("apply_to_each_user", cypherAccessPolicy.EachUser)
	default:
		d.Set("scope", "global")
	}

	var tenantIds []int64
	if cypherAccessPolicy.Accounts != nil {
		// iterate over the array of accounts
		for _, account := range cypherAccessPolicy.Accounts {
			tenantIds = append(tenantIds, account.ID)
		}
	}
	d.Set("tenant_ids", tenantIds)

	return diags
}

func resourceCypherAccessPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	policy := make(map[string]interface{})

	policy["name"] = d.Get("name").(string)
	policy["description"] = d.Get("description").(string)
	policy["enabled"] = d.Get("enabled").(bool)
	policy["config"] = map[string]interface{}{
		"keyPattern": d.Get("key_path").(string),
		"read":       evaluateStringBoolean(d.Get("read_access").(bool)),
		"write":      evaluateStringBoolean(d.Get("write_access").(bool)),
		"update":     evaluateStringBoolean(d.Get("update_access").(bool)),
		"delete":     evaluateStringBoolean(d.Get("delete_access").(bool)),
		"list":       evaluateStringBoolean(d.Get("list_access").(bool)),
	}
	policy["policyType"] = map[string]interface{}{
		"code": "cypher",
		"name": "Cypher Access",
	}

	policy["accounts"] = d.Get("tenant_ids")

	switch d.Get("scope") {
	case "user":
		policy["refId"] = d.Get("user_id").(int)
		policy["refType"] = "User"
		policy["user"] = map[string]interface{}{
			"id": d.Get("user_id").(int),
		}
	case "role":
		policy["refId"] = d.Get("role_id").(int)
		policy["refType"] = "Role"
		policy["eachUser"] = d.Get("apply_to_each_user").(bool)
		policy["role"] = map[string]interface{}{
			"id": d.Get("role_id").(int),
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
	return resourceCypherAccessPolicyRead(ctx, d, meta)
}

func resourceCypherAccessPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func evaluateStringBoolean(value bool) (output string) {
	if value {
		return "on"
	} else {
		return ""
	}
}

func parseStringBoolean(value string) (output bool) {
	if value == "on" {
		return true
	} else {
		return false
	}
}
