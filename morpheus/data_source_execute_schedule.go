package morpheus

import (
	"context"
	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMorpheusExecuteSchedule() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus execute schedule data source.",
		ReadContext: dataSourceMorpheusExecuteScheduleRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"name"},
				Computed:      true,
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the execute schedule",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
		},
	}
}

func dataSourceMorpheusExecuteScheduleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	id := d.Get("id").(int)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == 0 && name != "" {
		resp, err = client.FindExecuteScheduleByName(name)
	} else if id != 0 {
		resp, err = client.GetExecuteSchedule(int64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Execute schedule cannot be read without name or id")
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
	result := resp.Result.(*morpheus.GetExecuteScheduleResult)
	executeSchedule := result.ExecuteSchedule
	if executeSchedule != nil {
		d.SetId(int64ToString(executeSchedule.ID))
		d.Set("name", executeSchedule.Name)
	} else {
		return diag.Errorf("Execute schedule not found in response data.") // should not happen
	}
	return diags
}
