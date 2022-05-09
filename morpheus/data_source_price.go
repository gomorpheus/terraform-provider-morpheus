package morpheus

import (
	"context"
	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMorpheusPrice() *schema.Resource {
	return &schema.Resource{
		Description: "The Price data source allows details of a Price to be retrieved by its name.",
		ReadContext: dataSourceMorpheusPriceRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Morpheus price.",
				Optional:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the Morpheus price",
				Computed:    true,
			},
		},
	}
}

func dataSourceMorpheusPriceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindPriceByName(name)
	} else if id != "" {
		resp, err = client.GetPrice(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Price cannot be read without name or id")
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
	result := resp.Result.(*morpheus.GetPriceResult)
	price := result.Price
	if price != nil {
		d.SetId(int64ToString(price.ID))
		d.Set("name", price.Name)
		d.Set("code", price.Code)
	} else {
		return diag.Errorf("Price not found in response data.") // should not happen
	}
	return diags
}
