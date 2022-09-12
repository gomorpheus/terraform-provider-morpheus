package morpheus

import (
	"context"
	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMorpheusPriceSet() *schema.Resource {
	return &schema.Resource{
		Description: "The Price Set data source allows details of a Price Set to be retrieved by its name.",
		ReadContext: dataSourceMorpheusPriceSetRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"name"},
				Computed:      true,
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the Morpheus price set.",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the Morpheus price set",
				Computed:    true,
			},
		},
	}
}

func dataSourceMorpheusPriceSetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	id := d.Get("id").(int)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == 0 && name != "" {
		resp, err = client.FindPriceSetByName(name)
	} else if id != 0 {
		resp, err = client.GetPriceSet(int64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Price set cannot be read without name or id")
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
	result := resp.Result.(*morpheus.GetPriceSetResult)
	priceSet := result.PriceSet
	if priceSet != nil {
		d.SetId(int64ToString(priceSet.ID))
		d.Set("name", priceSet.Name)
		d.Set("code", priceSet.Code)
	} else {
		return diag.Errorf("Price set not found in response data.") // should not happen
	}
	return diags
}
