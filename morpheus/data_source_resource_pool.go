package morpheus

import (
	"context"
	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMorpheusResourcePool() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus resource pool data source.",
		ReadContext: dataSourceMorpheusResourcePoolRead,
		Schema: map[string]*schema.Schema{
			"cloud_id": {
				Type:        schema.TypeInt,
				Description: "The name of the Morpheus cloud.",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Morpheus cloud.",
				Optional:    true,
			},
			"type": {
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
		},
	}
}

func dataSourceMorpheusResourcePoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)
	cloud_id := d.Get("cloud_id").(int)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindResourcePoolByName(int64(cloud_id), name)
	} else if id != "" {
		resp, err = client.GetResourcePool(int64(cloud_id), toInt64(id), &morpheus.Request{})
		// todo: ignore 404 errors...
	} else {
		return diag.Errorf("Resource Pool cannot be read without name or id")
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
	result := resp.Result.(*morpheus.GetResourcePoolResult)
	resourcePool := result.ResourcePool
	if resourcePool != nil {
		d.SetId(int64ToString(resourcePool.ID))
		d.Set("name", resourcePool.Name)
		d.Set("active", resourcePool.Active)
		d.Set("type", resourcePool.Type)
		d.Set("description", resourcePool.Description)
		d.Set("visibility", resourcePool.Visibility)
	} else {
		return diag.Errorf("Resource pool not found in response data.") // should not happen
	}
	return diags
}
