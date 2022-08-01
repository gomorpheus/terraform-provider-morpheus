package morpheus

import (
	"context"
	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMorpheusPowerSchedule() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus power schedule data source.",
		ReadContext: dataSourceMorpheusPowerScheduleRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Description:   "The ID of the power schedule",
				Optional:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the power schedule",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
		},
	}
}

func dataSourceMorpheusPowerScheduleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	id := d.Get("id").(int)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == 0 && name != "" {
		resp, err = client.FindPowerScheduleByName(name)
	} else if id != 0 {
		resp, err = client.GetPowerSchedule(int64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Power schedule cannot be read without name or id")
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
	result := resp.Result.(*morpheus.GetPowerScheduleResult)
	powerSchedule := result.PowerSchedule
	if powerSchedule != nil {
		d.SetId(int64ToString(powerSchedule.ID))
		d.Set("name", powerSchedule.Name)
	} else {
		return diag.Errorf("Power schedule not found in response data.") // should not happen
	}
	return diags
}
