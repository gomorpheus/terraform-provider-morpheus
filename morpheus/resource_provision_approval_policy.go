package morpheus

import (
	"context"
	"strconv"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceProvisionApprovalPolicy() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus provision approval policy resource",
		CreateContext: resourceProvisionApprovalPolicyCreate,
		ReadContext:   resourceProvisionApprovalPolicyRead,
		UpdateContext: resourceProvisionApprovalPolicyUpdate,
		DeleteContext: resourceProvisionApprovalPolicyDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the provision approval policy",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the provision approval policy",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the provision approval policy",
				Optional:    true,
				Computed:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the policy is enabled",
				Optional:    true,
				Default:     true,
			},
			"use_internal_approvals": {
				Type:        schema.TypeBool,
				Description: "Whether the internal Morpheus approval engine is used for approvals",
				Optional:    true,
			},
			"integration_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the approval integration used for approvals",
				Optional:    true,
			},
			"workflow_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the approval workflow used for approvals",
				Optional:    true,
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

func resourceProvisionApprovalPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	policy := make(map[string]interface{})

	policy["name"] = d.Get("name").(string)
	policy["description"] = d.Get("description").(string)
	policy["enabled"] = d.Get("enabled").(bool)
	config := make(map[string]interface{})
	if d.Get("use_internal_approvals").(bool) {
		config["accountIntegrationId"] = strconv.Itoa(-100)
	} else {
		config["accountIntegrationId"] = strconv.Itoa(d.Get("integration_id").(int))
		config["workflowId"] = strconv.Itoa(d.Get("workflow_id").(int))
	}

	policy["config"] = config
	policy["policyType"] = map[string]interface{}{
		"code": "provisionApproval",
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

	resourceProvisionApprovalPolicyRead(ctx, d, meta)
	return diags
}

func resourceProvisionApprovalPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	provisionApprovalPolicy := result.Policy

	d.SetId(int64ToString(provisionApprovalPolicy.ID))
	d.Set("name", provisionApprovalPolicy.Name)
	d.Set("description", provisionApprovalPolicy.Description)
	d.Set("enabled", provisionApprovalPolicy.Enabled)
	if provisionApprovalPolicy.Config.AccountIntegrationId == "-100" {
		d.Set("use_internal_approvals", true)
	} else {
		//	d.Set("use_internal_approvals", false)
		integration_number, err := strconv.Atoi(provisionApprovalPolicy.Config.AccountIntegrationId)
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("integration_id", integration_number)

		workflow_number, err := strconv.Atoi(provisionApprovalPolicy.Config.WorkflowID.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("workflow_id", workflow_number)
	}

	switch provisionApprovalPolicy.RefType {
	case "ComputeSite":
		d.Set("scope", "group")
		d.Set("group_id", provisionApprovalPolicy.Site.ID)
	case "ComputeZone":
		d.Set("scope", "cloud")
		d.Set("cloud_id", provisionApprovalPolicy.Zone.ID)
	case "User":
		d.Set("scope", "user")
		d.Set("user_id", provisionApprovalPolicy.User.ID)
	case "Role":
		d.Set("scope", "role")
		d.Set("role_id", provisionApprovalPolicy.Role.ID)
		d.Set("apply_to_each_user", provisionApprovalPolicy.EachUser)
	default:
		d.Set("scope", "global")
	}

	var tenantIds []int64
	if provisionApprovalPolicy.Accounts != nil {
		// iterate over the array of accounts
		for _, account := range provisionApprovalPolicy.Accounts {
			tenantIds = append(tenantIds, account.ID)
		}
	}
	d.Set("tenant_ids", tenantIds)

	return diags
}

func resourceProvisionApprovalPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	policy := make(map[string]interface{})

	policy["name"] = d.Get("name").(string)
	policy["description"] = d.Get("description").(string)
	policy["enabled"] = d.Get("enabled").(bool)
	config := make(map[string]interface{})
	if d.Get("use_internal_approvals").(bool) {
		config["accountIntegrationId"] = strconv.Itoa(-100)
	} else {
		config["accountIntegrationId"] = strconv.Itoa(d.Get("integration_id").(int))
		config["workflowId"] = strconv.Itoa(d.Get("workflow_id").(int))
	}

	policy["config"] = config
	policy["policyType"] = map[string]interface{}{
		"code": "provisionApproval",
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
	return resourceWorkflowPolicyRead(ctx, d, meta)
}

func resourceProvisionApprovalPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
