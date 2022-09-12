package morpheus

import (
	"context"
	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMorpheusClusterType() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus cluster type source.",
		ReadContext: dataSourceMorpheusClusterTypeRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Morpheus cluster type.",
				Required:    true,
			},
		},
	}
}

func dataSourceMorpheusClusterTypeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	resp, err = client.FindClusterTypeByName(name)

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
	result := resp.Result.(*morpheus.ListClusterTypesResult)
	clusterTypesPayload := result.ClusterTypes
	clusterTypes := *clusterTypesPayload
	clusterType := clusterTypes[0]
	if result.Meta.Total > 0 {
		d.SetId(int64ToString(clusterType.ID))
		d.Set("name", clusterType.Name)
	} else {
		return diag.Errorf("cluster type not found in response data.") // should not happen
	}
	return diags
}
