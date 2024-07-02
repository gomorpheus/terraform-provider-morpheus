package morpheus

import (
	"context"
	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMorpheusCloudType() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus cloud type data source.",
		ReadContext: dataSourceMorpheusCloudTypeRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Morpheus cloud type",
				Required:    true,
			},
		},
	}
}

func dataSourceMorpheusCloudTypeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)

	var resp *morpheus.Response
	var err error

	resp, err = client.Execute(&morpheus.Request{
		Method: "GET",
		QueryParams: map[string]string{
			"name": name,
		},
		Path:   "/api/appliance-settings/zone-types",
		Result: &CloudTypes{},
	})

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
	cloudType := resp.Result.(*CloudTypes)

	if cloudType.ZoneTypes != nil {
		for _, cType := range cloudType.ZoneTypes {
			if cType.Name == name {
				d.SetId(int64ToString(cType.ID))
				d.Set("name", cType.Name)
			}
		}
	} else {
		return diag.Errorf("cloud type not found in response data.") // should not happen
	}
	return diags
}

type CloudTypes struct {
	ZoneTypes []CloudType       `json:"zoneTypes"`
	Message   string            `json:"msg"`
	Errors    map[string]string `json:"errors"`
}

type CloudType struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}
