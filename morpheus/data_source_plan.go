package morpheus

import (
	"context"
	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMorpheusPlan() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus plan data source.",
		ReadContext: dataSourceMorpheusPlanRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Morpheus cloud.",
				Optional:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "Optional code for use with policies",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the plan",
				Computed:    true,
			},
		},
	}
}

func dataSourceMorpheusPlanRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindPlanByName(name)
	} else if id != "" {
		resp, err = client.GetPlan(toInt64(id), &morpheus.Request{})
		// todo: ignore 404 errors...
	} else {
		return diag.Errorf("Plan cannot be read without name or id")
	}
	if err != nil {
		// 404 is ok?
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
	result := resp.Result.(*morpheus.GetPlanResult)
	plan := result.Plan
	if plan != nil {
		d.SetId(int64ToString(plan.ID))
		d.Set("name", plan.Name)
		d.Set("code", plan.Code)
		d.Set("description", plan.Description)
		d.Set("visibility", plan.Visibility)
	} else {
		return diag.Errorf("Plan not found in response data.") // should not happen
	}
	return diags
}
