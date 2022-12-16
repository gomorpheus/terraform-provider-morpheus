package morpheus

import (
	"context"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceInstanceNamePolicy() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus instance name policy resource",
		CreateContext: resourceInstanceNamePolicyCreate,
		ReadContext:   resourceInstanceNamePolicyRead,
		UpdateContext: resourceInstanceNamePolicyUpdate,
		DeleteContext: resourceInstanceNamePolicyDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the instance naming policy",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the instance naming policy",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the instance naming policy",
				Optional:    true,
				Computed:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the policy is enabled",
				Optional:    true,
				Default:     true,
			},
			"enforcement_type": {
				Type:         schema.TypeString,
				Description:  "The policy enforcement type (fixed or user)",
				ValidateFunc: validation.StringInSlice([]string{"fixed", "user"}, false),
				Required:     true,
			},
			"naming_pattern": {
				Type:        schema.TypeString,
				Description: "The instance name naming pattern",
				Required:    true,
			},
			"auto_resolve_conflicts": {
				Type:        schema.TypeBool,
				Description: "Whether to automatically resolve naming conflicts",
				Required:    true,
			},
			"scope": {
				Type:         schema.TypeString,
				Description:  "The filter or scope that the policy is applied to (global, group, cloud, user, role)",
				ValidateFunc: validation.StringInSlice([]string{"global", "group", "cloud", "user", "role"}, false),
				Required:     true,
				ForceNew:     true,
			},
			"group_id": {
				Type:          schema.TypeInt,
				Description:   "The id of the group associated with the group scoped filter",
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"cloud_id", "user_id", "role_id"},
			},
			"cloud_id": {
				Type:          schema.TypeInt,
				Description:   "The id of the cloud associated with the cloud scoped filter",
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"group_id", "user_id", "role_id"},
			},
			"user_id": {
				Type:          schema.TypeInt,
				Description:   "The id of the user associated with the user scoped filter",
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"cloud_id", "group_id", "role_id"},
			},
			"role_id": {
				Type:          schema.TypeInt,
				Description:   "The id of the role associated with the role scoped filter",
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"cloud_id", "user_id", "group_id"},
			},
			"apply_to_each_user": {
				Type:          schema.TypeBool,
				Description:   "Whether to assign the policy at the individual user level to all users assigned the associated role",
				Optional:      true,
				ConflictsWith: []string{"cloud_id", "user_id", "group_id"},
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

func resourceInstanceNamePolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	policy := make(map[string]interface{})

	policy["name"] = d.Get("name").(string)
	policy["description"] = d.Get("description").(string)
	policy["enabled"] = d.Get("enabled").(bool)
	var nameConflictStatus string
	if d.Get("auto_resolve_conflicts").(bool) {
		nameConflictStatus = "on"
	} else {
		nameConflictStatus = ""
	}
	policy["config"] = map[string]interface{}{
		"namingConflict": nameConflictStatus,
		"namingType":     d.Get("enforcement_type").(string),
		"namingPattern":  d.Get("naming_pattern").(string),
	}
	policy["policyType"] = map[string]interface{}{
		"code": "naming",
		"name": "Instance Name",
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

	resourceInstanceNamePolicyRead(ctx, d, meta)
	return diags
}

func resourceInstanceNamePolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	instanceNamePolicy := result.Policy

	d.SetId(int64ToString(instanceNamePolicy.ID))
	d.Set("name", instanceNamePolicy.Name)
	d.Set("description", instanceNamePolicy.Description)
	d.Set("enabled", instanceNamePolicy.Enabled)
	d.Set("enforcement_type", instanceNamePolicy.Config.NamingType)
	d.Set("auto_resolve_conflicts", instanceNamePolicy.Config.NamingConflict)
	d.Set("naming_pattern", instanceNamePolicy.Config.NamingPattern)
	switch instanceNamePolicy.RefType {
	case "ComputeSite":
		d.Set("scope", "group")
		d.Set("group_id", instanceNamePolicy.Site.ID)
	case "ComputeZone":
		d.Set("scope", "cloud")
		d.Set("cloud_id", instanceNamePolicy.Zone.ID)
	case "User":
		d.Set("scope", "user")
		d.Set("user_id", instanceNamePolicy.User.ID)
	case "Role":
		d.Set("scope", "role")
		d.Set("role_id", instanceNamePolicy.Role.ID)
		d.Set("apply_to_each_user", instanceNamePolicy.EachUser)
	default:
		d.Set("scope", "global")
	}

	var tenantIds []int64
	if instanceNamePolicy.Accounts != nil {
		// iterate over the array of accounts
		for _, account := range instanceNamePolicy.Accounts {
			tenantIds = append(tenantIds, account.ID)
		}
	}
	d.Set("tenant_ids", tenantIds)

	return diags
}

func resourceInstanceNamePolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	policy := make(map[string]interface{})

	policy["name"] = d.Get("name").(string)
	policy["description"] = d.Get("description").(string)
	policy["enabled"] = d.Get("enabled").(bool)
	var nameConflictStatus string
	if d.Get("auto_resolve_conflicts").(bool) {
		nameConflictStatus = "on"
	} else {
		nameConflictStatus = ""
	}
	policy["config"] = map[string]interface{}{
		"namingConflict": nameConflictStatus,
		"namingType":     d.Get("enforcement_type").(string),
		"namingPattern":  d.Get("naming_pattern").(string),
	}
	policy["policyType"] = map[string]interface{}{
		"code": "naming",
		"name": "Instance Name",
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
	return resourceInstanceNamePolicyRead(ctx, d, meta)
}

func resourceInstanceNamePolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
