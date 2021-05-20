package morpheus

import (
	"context"
	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMorpheusNetwork() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus network data source.",
		ReadContext: dataSourceMorpheusNetworkRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Morpheus network.",
				Optional:    true,
			},
			"active": {
				Type:        schema.TypeBool,
				Description: "Optional code for use with policies",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the network",
				Computed:    true,
			},
			"visibility": {
				Type:        schema.TypeString,
				Description: "Whether the network is visible in sub-tenants or not",
				Computed:    true,
			},
		},
	}
}

func dataSourceMorpheusNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindNetworkByName(name)
	} else if id != "" {
		resp, err = client.GetNetwork(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Network cannot be read without name or id")
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
	result := resp.Result.(*morpheus.GetNetworkResult)
	network := result.Network
	if network != nil {
		d.SetId(int64ToString(network.ID))
		d.Set("name", network.Name)
		d.Set("active", network.Active)
		d.Set("description", network.Description)
		d.Set("visibility", network.Visibility)
	} else {
		return diag.Errorf("Network not found in response data.") // should not happen
	}
	return diags
}
