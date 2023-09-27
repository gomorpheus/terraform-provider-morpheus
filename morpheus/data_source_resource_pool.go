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
				Description: "The id of the Morpheus cloud to search for the resource pool.",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Morpheus resource pool.",
				Optional:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "Optional code for use with policies",
				Computed:    true,
			},
			"active": {
				Type:        schema.TypeBool,
				Description: "Whether the resource pool is enabled or not",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the resource pool",
				Computed:    true,
			},
			"id": {
				Type:        schema.TypeInt,
				Description: "The id of the resource pool",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func dataSourceMorpheusResourcePoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Get("id").(int)
	name := d.Get("name").(string)
	cloud_id := d.Get("cloud_id").(int)

	// Ensure that either the id or name is provided
	if id == 0 && name == "" {
		return diag.Errorf("Either 'id' or 'name' must be provided to search for the resource pool")
	}

	var resp *morpheus.Response
	var err error

	if id != 0 {
		resp, err = client.GetResourcePool(int64(cloud_id), int64(id), &morpheus.Request{})
	} else {
		resp, err = client.FindResourcePoolByName(int64(cloud_id), name)
	}

	if err != nil {
		errorPrefix := "API FAILURE"
		if resp != nil && resp.StatusCode == 404 {
			errorPrefix = "API 404"
		}
		log.Printf("%s: %s - %v", errorPrefix, resp, err)
		return diag.FromErr(err)
	}

	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.GetResourcePoolResult)
	resourcePool := result.ResourcePool

	if resourcePool == nil {
		return diag.Errorf("Resource pool not found in response data.")
	}

	d.SetId(int64ToString(resourcePool.ID))
	d.Set("name", resourcePool.Name)
	d.Set("active", resourcePool.Active)
	d.Set("type", resourcePool.Type)
	d.Set("description", resourcePool.Description)

	return diags
}
