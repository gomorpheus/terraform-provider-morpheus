package morpheus

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMorpheusPolicy() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus policy data source.",
		ReadContext: dataSourceMorpheusPolicyRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Description:   "The ID of the Morpheus policy.",
				Optional:      true,
				ConflictsWith: []string{"name"},
				Computed:      true,
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the Morpheus policy.",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the policy",
				Computed:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the policy is enabled",
				Computed:    true,
			},
			"scope": {
				Type:        schema.TypeString,
				Description: "The filter scope of the policy",
				Computed:    true,
			},
			"policy_type_name": {
				Type:        schema.TypeString,
				Description: "The name of the policy type",
				Computed:    true,
			},
			"policy_type_code": {
				Type:        schema.TypeString,
				Description: "The code of the policy type",
				Computed:    true,
			},
			"tenant_ids": {
				Type:        schema.TypeList,
				Description: "Tenants the policy is assigned to",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func dataSourceMorpheusPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	id := d.Get("id").(int)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == 0 && name != "" {
		resp, err = client.FindPolicyByName(name)
	} else if id != 0 {
		resp, err = client.GetPolicy(int64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Policy cannot be read without name or id")
	}
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %v", resp, err)
			return nil
		} else {
			log.Printf("API FAILURE: %s - %v", resp, err)
			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)

	// store resource data
	result := resp.Result.(*morpheus.GetPolicyResult)
	policy := result.Policy

	// store resource data
	var policyPayload PolicyPayload
	json.Unmarshal(resp.Body, &policyPayload)

	if policy != nil {
		d.SetId(int64ToString(int64(policyPayload.Policy.ID)))
		d.Set("name", policyPayload.Policy.Name)
		d.Set("enabled", policyPayload.Policy.Enabled)
		d.Set("description", policyPayload.Policy.Description)

		switch policyPayload.Policy.Reftype {
		case "ComputeSite":
			d.Set("scope", "group")
		case "ComputeZone":
			d.Set("scope", "cloud")
		case "User":
			d.Set("scope", "user")
		case "Role":
			d.Set("scope", "role")
		default:
			d.Set("scope", "global")
		}
		d.Set("policy_type_name", policyPayload.Policy.Policytype.Name)
		d.Set("policy_type_code", policyPayload.Policy.Policytype.Code)

		// tenants
		var tenants []int64
		if policyPayload.Policy.Accounts != nil {
			// iterate over the array of tasks
			for i := 0; i < len(policyPayload.Policy.Accounts); i++ {
				tenant := policyPayload.Policy.Accounts[i].(map[string]interface{})
				tenantID := int64(tenant["id"].(float64))
				tenants = append(tenants, tenantID)
			}
		}
		d.Set("tenant_ids", tenants)
	} else {
		return diag.Errorf("Policy not found in response data.") // should not happen
	}
	return diags
}

type PolicyPayload struct {
	Policy struct {
		Accounts    []interface{} `json:"accounts"`
		Config      []interface{} `json:"config"`
		Description string        `json:"description"`
		Eachuser    bool          `json:"eachUser"`
		Enabled     bool          `json:"enabled"`
		ID          int           `json:"id"`
		Name        string        `json:"name"`
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
