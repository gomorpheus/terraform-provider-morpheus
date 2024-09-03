package morpheus

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMorpheusChefServer() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus chef server data source.",
		ReadContext: dataSourceMorpheusChefServerRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the chef server",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
		},
	}
}

func dataSourceMorpheusChefServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	value := d.Get("id")

	// lookup by name if we do not have an value yet
	var resp *morpheus.Response
	var err error

	resp, err = client.GetOptionSource("chefServer", &morpheus.Request{})
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

	var server morpheus.OptionSourceOption
	chefServers := *resp.Result.(*morpheus.GetOptionSourceResult).Data
	for i := range chefServers {
		if value == 0 && name != "" {
			if strings.EqualFold(chefServers[i].Name, name) {
				server = chefServers[i]
				break
			}
		} else if value != 0 {
			if value == chefServers[i].Value {
				server = chefServers[i]
				break
			}
		} else {
			return diag.Errorf("Chef server cannot be read without name or value")
		}
	}

	// store resource data
	d.SetId(fmt.Sprintf("%g", server.Value.(float64)))
	d.Set("id", server.Value)
	d.Set("name", server.Name)

	return diags
}
