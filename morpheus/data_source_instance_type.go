package morpheus

import (
	"context"
	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMorpheusInstanceType() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus resource pool data source.",
		ReadContext: dataSourceMorpheusInstanceTypeRead,
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
			"active": {
				Type:        schema.TypeBool,
				Description: "Optional code for use with policies",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the plan",
				Computed:    true,
			},
			"visibility": {
				Type:        schema.TypeString,
				Description: "The description of the plan",
				Computed:    true,
			},
		},
	}
}

func dataSourceMorpheusInstanceTypeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindInstanceTypeByName(name)
	} else if id != "" {
		resp, err = client.GetInstanceType(toInt64(id), &morpheus.Request{})
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
	result := resp.Result.(*morpheus.GetInstanceTypeResult)
	instanceType := result.InstanceType
	if instanceType != nil {
		d.SetId(int64ToString(instanceType.ID))
		d.Set("name", instanceType.Name)
		d.Set("code", instanceType.Code)
		d.Set("active", instanceType.Active)
		d.Set("description", instanceType.Description)
		d.Set("visibility", instanceType.Visibility)
	} else {
		return diag.Errorf("Instance type not found in response data.") // should not happen
	}
	return diags
}
