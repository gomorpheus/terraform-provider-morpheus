package morpheus

import (
	"context"
	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMorpheusEnvironments() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus environments data source.",
		ReadContext: dataSourceMorpheusEnvironmentsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"sort_ascending": {
				Type:        schema.TypeBool,
				Description: "Whether to sort the IDs in ascending order",
				Default:     true,
				Optional:    true,
			},
			/*
				"filter": {
					Type:        schema.TypeSet,
					Description: "The environment variables to create",
					Optional:    true,
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name": {
								Type:        schema.TypeString,
								Description: "The name of the environment variable",
								Required:    true,
							},
							"values": {
								Type:     schema.TypeSet,
								Required: true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
						},
					},
				},
			*/
		},
	}
}

func dataSourceMorpheusEnvironmentsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var resp *morpheus.Response
	var err error
	var sortOrder string
	/*
		payload := d.Get("filter").(*schema.Set).List()[0]
		test := payload.(map[string]interface{})

		log.Printf("FILTER INFO - %s - ", test["name"].(string))
		log.Printf("FILTER INFO - %v -", test["values"].(*schema.Set).List())
	*/
	// Sort environments in ascending or descending order
	if d.Get("sort_ascending").(bool) {
		sortOrder = "asc"
	} else {
		sortOrder = "desc"
	}

	resp, err = client.ListEnvironments(&morpheus.Request{
		QueryParams: map[string]string{
			"max":       "50",
			"sort":      "id",
			"direction": sortOrder,
		},
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

	environmentIDs := []int64{}

	// store resource data
	result := resp.Result.(*morpheus.ListEnvironmentsResult)
	environments := result.Environments
	for _, environment := range *environments {
		environmentIDs = append(environmentIDs, environment.ID)
	}
	d.SetId("1")
	d.Set("ids", environmentIDs)
	return diags
}
