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

func dataSourceMorpheusVrealizeOrchestratorWorkflow() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus vRO workflow data source.",
		ReadContext: dataSourceMorpheusVrealizeOrchestratorWorkflowRead,
		Schema: map[string]*schema.Schema{
			"value": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the option type",
				Optional:      true,
				ConflictsWith: []string{"value"},
			},
		},
	}
}

func dataSourceMorpheusVrealizeOrchestratorWorkflowRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	value := d.Get("value")

	// lookup by name if we do not have an value yet
	var resp *morpheus.Response
	var err error

	resp, err = client.GetOptionSource("vroWorkflow", &morpheus.Request{})
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

	var workflow morpheus.OptionSourceOption
	allWorkflows := *resp.Result.(*morpheus.GetOptionSourceResult).Data
	for i := range allWorkflows {
		if value == 0 && name != "" {
			if strings.EqualFold(allWorkflows[i].Name, name) {
				workflow = allWorkflows[i]
				break
			}
		} else if value != 0 {
			if value == allWorkflows[i].Value {
				workflow = allWorkflows[i]
				break
			}
		} else {
			return diag.Errorf("vRO workflow cannot be read without name or value")
		}
	}

	// store resource data
	d.SetId(fmt.Sprintf("%g", workflow.Value.(float64)))
	d.Set("value", workflow.Value)
	d.Set("name", workflow.Name)

	return diags
}
