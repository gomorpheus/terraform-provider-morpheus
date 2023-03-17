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

func dataSourceMorpheusAnsibleTowerJobTemplate() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus ansible tower job template data source.",
		ReadContext: dataSourceMorpheusAnsibleTowerJobTemplateRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the ansible tower job template",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
		},
	}
}

func dataSourceMorpheusAnsibleTowerJobTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	value := d.Get("id")

	// lookup by name if we do not have an value yet
	var resp *morpheus.Response
	var err error

	resp, err = client.GetOptionSource("ansibleTowerJobTemplate", &morpheus.Request{})
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

	var template morpheus.OptionSourceOption
	allTemplates := *resp.Result.(*morpheus.GetOptionSourceResult).Data
	for i := range allTemplates {
		if value == 0 && name != "" {
			if strings.EqualFold(allTemplates[i].Name, name) {
				template = allTemplates[i]
				break
			}
		} else if value != 0 {
			if value == allTemplates[i].Value {
				template = allTemplates[i]
				break
			}
		} else {
			return diag.Errorf("Ansible tower job template cannot be read without name or value")
		}
	}

	// store resource data
	d.SetId(fmt.Sprintf("%g", template.Value.(float64)))
	d.Set("id", template.Value)
	d.Set("name", template.Name)

	return diags
}
