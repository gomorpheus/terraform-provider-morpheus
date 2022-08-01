package morpheus

import (
	"context"
	"encoding/json"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceMaxStoragePolicy() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus max storage policy resource",
		CreateContext: resourceMaxStoragePolicyCreate,
		ReadContext:   resourceMaxStoragePolicyRead,
		UpdateContext: resourceMaxStoragePolicyUpdate,
		DeleteContext: resourceMaxStoragePolicyDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the workflow policy",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the max storage policy",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the max storage policy",
				Optional:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the policy is enabled",
				Optional:    true,
				Default:     true,
			},
			"max_storage": {
				Type:        schema.TypeInt,
				Description: "The maximum storage defined by the policy in GB",
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
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceMaxStoragePolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	policy := make(map[string]interface{})

	policy["name"] = d.Get("name").(string)
	policy["description"] = d.Get("description").(string)
	policy["enabled"] = d.Get("enabled").(bool)
	policy["config"] = map[string]interface{}{
		"maxStorage": d.Get("max_storage").(int),
	}
	policy["policyType"] = map[string]interface{}{
		"code": "maxStorage",
		"name": "Max Storage",
	}

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

	resourceMaxStoragePolicyRead(ctx, d, meta)
	return diags
}

func resourceMaxStoragePolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	var maxStoragePolicy MaxStoragePolicy
	json.Unmarshal(resp.Body, &maxStoragePolicy)

	d.SetId(intToString(maxStoragePolicy.Policy.ID))
	d.Set("name", maxStoragePolicy.Policy.Name)
	d.Set("description", maxStoragePolicy.Policy.Description)
	d.Set("enabled", maxStoragePolicy.Policy.Enabled)
	d.Set("max_storage", maxStoragePolicy.Policy.Config.MaxStorage)

	switch maxStoragePolicy.Policy.Reftype {
	case "ComputeSite":
		d.Set("scope", "group")
		d.Set("group_id", maxStoragePolicy.Policy.Site.ID)
	case "ComputeZone":
		d.Set("scope", "cloud")
		d.Set("cloud_id", maxStoragePolicy.Policy.Zone.ID)
	case "User":
		d.Set("scope", "user")
		d.Set("user_id", maxStoragePolicy.Policy.User.ID)
	case "Role":
		d.Set("scope", "role")
		d.Set("role_id", maxStoragePolicy.Policy.Role.ID)
		d.Set("apply_to_each_user", maxStoragePolicy.Policy.Eachuser)
	default:
		d.Set("scope", "global")
	}

	return diags
}

func resourceMaxStoragePolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	policy := make(map[string]interface{})

	policy["name"] = d.Get("name").(string)
	policy["description"] = d.Get("description").(string)
	policy["enabled"] = d.Get("enabled").(bool)
	policy["config"] = map[string]interface{}{
		"maxStorage": d.Get("max_storage").(int),
	}
	policy["policyType"] = map[string]interface{}{
		"code": "maxStorage",
		"name": "Max Storage",
	}

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
	return resourceWorkflowPolicyRead(ctx, d, meta)
}

func resourceMaxStoragePolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

type MaxStoragePolicy struct {
	Policy struct {
		Accounts []interface{} `json:"accounts"`
		Config   struct {
			MaxStorage string `json:"maxStorage"`
		} `json:"config"`
		Description string `json:"description"`
		Eachuser    bool   `json:"eachUser"`
		Enabled     bool   `json:"enabled"`
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Owner       struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"owner"`
		Policytype struct {
			Code string `json:"code"`
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"policyType"`
		Refid   int    `json:"refId"`
		Reftype string `json:"refType"`
		Role    struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"role"`
		Site struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"site"`
		User struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"user"`
		Zone struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"zone"`
	}
}