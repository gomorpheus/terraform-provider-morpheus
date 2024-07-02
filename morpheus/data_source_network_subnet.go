package morpheus

import (
	"context"
	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMorpheusNetworkSubnet() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus network subnet data source.",
		ReadContext: dataSourceMorpheusNetworkSubnetRead,
		Schema: map[string]*schema.Schema{
			"network_id": {
				Type:        schema.TypeInt,
				Description: "The id of the Morpheus network to search for the subnet.",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Morpheus network subnet.",
				Optional:    true,
			},
			//"display_name": {
			//	Type:        schema.TypeString,
			//	Description: "The user friendly name of the network subnet",
			//	Computed:    true,
			//},
			"id": {
				Type:        schema.TypeInt,
				Description: "The id of the network subnet",
				Optional:    true,
				Computed:    true,
			},
			"external_id": {
				Type:        schema.TypeString,
				Description: "The external id of the network subnet",
				Computed:    true,
			},
			"cidr": {
				Type:        schema.TypeString,
				Description: "The cidr of the network subnet",
				Computed:    true,
			},
			"netmask": {
				Type:        schema.TypeString,
				Description: "The netmask of the network subnet",
				Computed:    true,
			},
			"visibility": {
				Type:        schema.TypeString,
				Description: "The visibility of the network subnet",
				Computed:    true,
			},
		},
	}
}

func dataSourceMorpheusNetworkSubnetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Get("id").(int)
	name := d.Get("name").(string)
	network_id := d.Get("network_id").(int)

	// Ensure that either the id or name is provided
	if id == 0 && name == "" {
		return diag.Errorf("Either 'id' or 'name' must be provided to search for the network subnet")
	}

	var resp *morpheus.Response
	var err error
	var networkSubnet *morpheus.NetworkSubnet

	if id != 0 {
		resp, err = client.GetNetworkSubnet(int64(id), &morpheus.Request{})

		if err != nil {
			errorPrefix := "API FAILURE"
			if resp != nil && resp.StatusCode == 404 {
				errorPrefix = "API 404"
			}
			log.Printf("%s: %s - %v", errorPrefix, resp, err)
			return diag.FromErr(err)
		}

		log.Printf("API RESPONSE: %s", resp)

		result := resp.Result.(*morpheus.GetNetworkSubnetResult)
		networkSubnet = result.NetworkSubnet

	} else {
		resp, err = client.ListNetworkSubnetsByNetwork(int64(network_id), &morpheus.Request{
			QueryParams: map[string]string{
				"name": name,
			},
		})
		if err != nil {
			log.Printf("API FAILURE: %s - %v", resp, err)
			return diag.FromErr(err)
		}

		listResult := resp.Result.(*morpheus.ListNetworkSubnetsByNetworkResult)
		networkSubnetsCount := len(*listResult.NetworkSubnets)
		if networkSubnetsCount != 1 {
			return diag.Errorf("found %d Network Subnets for %v", networkSubnetsCount, name)
		}
		firstRecord := (*listResult.NetworkSubnets)[0]
		networkSubnetID := firstRecord.ID
		resp, err = client.GetNetworkSubnet(networkSubnetID, &morpheus.Request{})
		if err != nil {
			errorPrefix := "API FAILURE"
			if resp != nil && resp.StatusCode == 404 {
				errorPrefix = "API 404"
			}
			log.Printf("%s: %s - %v", errorPrefix, resp, err)
			return diag.FromErr(err)
		}

		log.Printf("API RESPONSE: %s", resp)
		result := resp.Result.(*morpheus.GetNetworkSubnetResult)
		networkSubnet = result.NetworkSubnet
	}

	if networkSubnet == nil {
		return diag.Errorf("Network subnet not found in response data.")
	}

	d.SetId(int64ToString(networkSubnet.ID))
	d.Set("name", networkSubnet.Name)
	d.Set("external_id", networkSubnet.ExternalId)
	d.Set("cidr", networkSubnet.Cidr)
	d.Set("netmask", networkSubnet.Netmask)
	d.Set("visibility", networkSubnet.Visibility)
	return diags
}
