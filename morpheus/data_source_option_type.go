package morpheus

import (
	"context"
	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMorpheusOptionType() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus option type data source.",
		ReadContext: dataSourceMorpheusOptionTypeRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the option type",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
		},
	}
}

func dataSourceMorpheusOptionTypeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	id := d.Get("id").(int)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == 0 && name != "" {
		resp, err = client.FindOptionTypeByName(name)
	} else if id != 0 {
		resp, err = client.GetOptionType(int64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Option type cannot be read without name or id")
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
	result := resp.Result.(*morpheus.GetOptionTypeResult)
	optionType := result.OptionType
	if optionType != nil {
		d.SetId(int64ToString(optionType.ID))
		d.Set("name", optionType.Name)
	} else {
		return diag.Errorf("Option type not found in response data.") // should not happen
	}
	return diags
}
