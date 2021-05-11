package morpheus

import (
	"context"
	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMorpheusInstanceLayout() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus instance layout data source.",
		ReadContext: dataSourceMorpheusInstanceLayoutRead,
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

func dataSourceMorpheusInstanceLayoutRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindInstanceLayoutByName(name)
	} else if id != "" {
		resp, err = client.GetInstanceLayout(toInt64(id), &morpheus.Request{})
		// todo: ignore 404 errors...
	} else {
		return diag.Errorf("Instance type cannot be read without name or id")
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
	result := resp.Result.(*morpheus.GetInstanceLayoutResult)
	instanceType := result.InstanceLayout
	if instanceType != nil {
		d.SetId(int64ToString(instanceType.ID))
		d.Set("name", instanceType.Name)
		d.Set("code", instanceType.Code)
		d.Set("description", instanceType.Description)
	} else {
		return diag.Errorf("Instance type not found in response data.") // should not happen
	}
	return diags
}
