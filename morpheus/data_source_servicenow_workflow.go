package morpheus

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMorpheusServiceNowWorkflow() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus ServiceNow workflow data source.",
		ReadContext: dataSourceMorpheusServiceNowWorkflowRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Description:   "The ID of the ServiceNow integration",
				Optional:      true,
				ConflictsWith: []string{"name"},
				Computed:      true,
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the ServiceNow workflow",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
			"integration_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the ServiceNow integration",
				Required:    true,
			},
		},
	}
}

func dataSourceMorpheusServiceNowWorkflowRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)

	resp, err := client.Execute(&morpheus.Request{
		Method:      "GET",
		Path:        fmt.Sprintf("/api/options/deleteApprovalServiceNowWorkflows?config.accountIntegrationId=%d", d.Get("integration_id").(int)),
		QueryParams: map[string]string{},
	})
	if err != nil {
		log.Println("API ERROR: ", err)
	}
	log.Println("API RESPONSE:", resp)

	var itemResponsePayload CodeRepositories
	json.Unmarshal(resp.Body, &itemResponsePayload)
	foundWorkflow := false
	for _, v := range itemResponsePayload.Data {
		if v.Name == name {
			foundWorkflow = true
			d.SetId(intToString(v.Value))
		}
	}
	if !foundWorkflow {
		return diag.Errorf("Workflow named %s not found", name) // should not happen
	}
	return diags
}
