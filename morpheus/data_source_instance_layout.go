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
				Description: "The name of the Morpheus instance layout",
				Optional:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "Optional code for use with policies",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the instance layout",
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
	} else {
		return diag.Errorf("Instance layout cannot be read without name or id")
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
	result := resp.Result.(*morpheus.GetInstanceLayoutResult)
	instanceLayout := result.InstanceLayout
	if instanceLayout != nil {
		d.SetId(int64ToString(instanceLayout.ID))
		d.Set("name", instanceLayout.Name)
		d.Set("code", instanceLayout.Code)
		d.Set("description", instanceLayout.Description)
	} else {
		return diag.Errorf("Instance layout not found in response data.") // should not happen
	}
	return diags
}
