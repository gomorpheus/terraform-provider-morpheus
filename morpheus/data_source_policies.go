package morpheus

import (
	"context"
	"log"
	"strconv"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceMorpheusPolicies() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus policies data source.",
		ReadContext: dataSourceMorpheusPoliciesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"sort_ascending": {
				Type:        schema.TypeBool,
				Description: "Whether to sort the IDs in ascending order. Defaults to true",
				Default:     true,
				Optional:    true,
			},
			"filter": {
				Type:        schema.TypeSet,
				Description: "Custom filter block as described below.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Description:  "The name of the filter. Filter names are case-sensitive. Valid names are (name, type)",
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"name", "type"}, false),
						},
						"values": {
							Type:        schema.TypeSet,
							Description: "The filter values. Filter values are case-sensitive. Filters values support the use of Golang regex and can be tested at https://regex101.com/",
							Required:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceMorpheusPoliciesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var resp *morpheus.Response
	var err error
	var sortOrder string
	var policyTypes []string
	var names []string

	if len(d.Get("filter").(*schema.Set).List()) > 0 {
		filters := d.Get("filter").(*schema.Set).List()
		for _, filter := range filters {
			test := filter.(map[string]interface{})
			if test["name"].(string) == "type" {
				for _, item := range test["values"].(*schema.Set).List() {
					policyTypes = append(policyTypes, item.(string))
				}
			}

			if test["name"].(string) == "name" {
				for _, item := range test["values"].(*schema.Set).List() {
					names = append(names, item.(string))
				}
			}
		}
	}

	if len(policyTypes) == 0 {
		policyTypes = append(policyTypes, "$")
	}

	if len(names) == 0 {
		names = append(names, "$")
	}

	// Sort environments in ascending or descending order
	if d.Get("sort_ascending").(bool) {
		sortOrder = "asc"
	} else {
		sortOrder = "desc"
	}

	resp, err = client.ListPolicies(&morpheus.Request{
		QueryParams: map[string]string{
			"max":       "100",
			"sort":      "id",
			"direction": sortOrder,
		},
	})

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

	var policyIDs []string

	// store resource data
	result := resp.Result.(*morpheus.ListPoliciesResult)
	policies := result.Policies
	for _, policy := range *policies {
		if regexCheck(policyTypes, policy.PolicyType.Name) && regexCheck(names, policy.Name) {
			policyIDs = append(policyIDs, strconv.Itoa(int(policy.ID)))
		}
	}
	d.SetId("1")
	d.Set("ids", policyIDs)
	return diags
}
