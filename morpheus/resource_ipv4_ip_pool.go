package morpheus

import (
	"context"
	"sort"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIPv4IPPool() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus IPv4 ip pool resource",
		CreateContext: resourceIPv4IPPoolCreate,
		ReadContext:   resourceIPv4IPPoolRead,
		UpdateContext: resourceIPv4IPPoolUpdate,
		DeleteContext: resourceIPv4IPPoolDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the IPv4 IP address pool",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the IPv4 IP address pool",
				Required:    true,
			},
			"ip_range": {
				Type:        schema.TypeList,
				Description: "The IPv4 IP address pool IP ranges",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"starting_address": {
							Type:        schema.TypeString,
							Description: "The starting address of the IPv4 IP address pool IP range",
							Required:    true,
						},
						"ending_address": {
							Type:        schema.TypeString,
							Description: "The ending address of the IPv4 IP address pool IP range",
							Required:    true,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceIPv4IPPoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"networkPool": map[string]interface{}{
				"name":     d.Get("name").(string),
				"type":     "morpheus",
				"ipRanges": parseIPPoolRanges(d.Get("ip_range").([]interface{})),
			},
		},
	}
	resp, err := client.CreateNetworkPool(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateNetworkPoolResult)
	pool := result.NetworkPool
	// Successfully created resource, now set id
	d.SetId(int64ToString(pool.ID))
	resourceIPv4IPPoolRead(ctx, d, meta)
	return diags
}

func resourceIPv4IPPoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindNetworkPoolByName(name)
	} else if id != "" {
		resp, err = client.GetNetworkPool(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Pool cannot be read without name or id")
	}

	if err != nil {
		// 404 is ok?
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
			return diag.FromErr(err)
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)
			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)

	// store resource data
	result := resp.Result.(*morpheus.GetNetworkPoolResult)
	pool := result.NetworkPool
	d.SetId(int64ToString(pool.ID))
	d.Set("name", pool.Name)
	var ipRanges []map[string]interface{}
	var unsortedRanges []IPRange
	if pool.IpRanges != nil {
		for _, iprange := range pool.IpRanges {
			var IPR IPRange
			IPR.ID = iprange.ID
			IPR.EndAddress = iprange.EndAddress
			IPR.StartAddress = iprange.StartAddress
			unsortedRanges = append(unsortedRanges, IPR)
		}
	}
	sort.Slice(unsortedRanges, func(i, j int) bool { return unsortedRanges[i].ID < unsortedRanges[j].ID })

	// iterate over the array of IP ranges
	for i := 0; i < len(unsortedRanges); i++ {
		ipRange := unsortedRanges[i]
		rangePayload := make(map[string]interface{})
		rangePayload["ending_address"] = ipRange.EndAddress
		rangePayload["starting_address"] = ipRange.StartAddress
		ipRanges = append(ipRanges, rangePayload)
	}
	d.Set("ip_range", ipRanges)
	return diags
}

func resourceIPv4IPPoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	req := &morpheus.Request{
		Body: map[string]interface{}{
			"networkPool": map[string]interface{}{
				"name":     d.Get("name").(string),
				"type":     "morpheus",
				"ipRanges": parseIPPoolRanges(d.Get("ip_range").([]interface{})),
			},
		},
	}
	resp, err := client.UpdateNetworkPool(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateNetworkPoolResult)
	pool := result.NetworkPool
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(pool.ID))
	return resourceIPv4IPPoolRead(ctx, d, meta)
}

func resourceIPv4IPPoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteNetworkPool(toInt64(id), req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
			return nil
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)
			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)
	d.SetId("")
	return diags
}

func parseIPPoolRanges(variables []interface{}) []map[string]interface{} {
	var poolRanges []map[string]interface{}
	// iterate over the array of poolRanges
	for i := 0; i < len(variables); i++ {
		row := make(map[string]interface{})
		ippoolconfig := variables[i].(map[string]interface{})
		for k, v := range ippoolconfig {
			switch k {
			case "starting_address":
				row["startAddress"] = v.(string)
			case "ending_address":
				row["endAddress"] = v.(string)
			}
		}
		poolRanges = append(poolRanges, row)
	}
	return poolRanges
}

type IPRange struct {
	ID           int64  `json:"id"`
	StartAddress string `json:"startAddress"`
	EndAddress   string `json:"endAddress"`
}
