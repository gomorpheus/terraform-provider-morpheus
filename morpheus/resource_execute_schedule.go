package morpheus

import (
	"context"
	"encoding/json"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceExecuteSchedule() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides an execution schedule resource",
		CreateContext: resourceExecuteScheduleCreate,
		ReadContext:   resourceExecuteScheduleRead,
		UpdateContext: resourceExecuteScheduleUpdate,
		DeleteContext: resourceExecuteScheduleDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the execute schedule",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the execute schedule",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the execute schedule",
				Optional:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the execute schedule is enabled",
				Optional:    true,
				Default:     true,
			},
			"time_zone": {
				Type:        schema.TypeString,
				Description: "The time zone used for scheduling",
				Required:    true,
			},
			"schedule": {
				Type:        schema.TypeString,
				Description: "The cron style syntax for the scheduled execution",
				Required:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceExecuteScheduleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	schedule := make(map[string]interface{})

	schedule["name"] = d.Get("name").(string)
	schedule["description"] = d.Get("description").(string)
	schedule["enabled"] = d.Get("enabled").(bool)
	schedule["scheduleType"] = "execute"
	schedule["scheduleTimezone"] = d.Get("time_zone").(string)
	schedule["cron"] = d.Get("schedule").(string)

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"schedule": schedule,
		},
	}
	resp, err := client.CreateExecuteSchedule(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateExecuteScheduleResult)
	executeScheduleResult := result.ExecuteSchedule
	// Successfully created resource, now set id
	d.SetId(int64ToString(executeScheduleResult.ID))

	resourceExecuteScheduleRead(ctx, d, meta)
	return diags
}

func resourceExecuteScheduleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindExecuteScheduleByName(name)
	} else if id != "" {
		resp, err = client.GetExecuteSchedule(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Execute schedule cannot be read without name or id")
	}

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
			log.Printf("Forcing recreation of resource")
			d.SetId("")
			return diags
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)
			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)

	// store resource data
	var executeSchedule ExecuteSchedule
	if err := json.Unmarshal(resp.Body, &executeSchedule); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(intToString(executeSchedule.Schedule.ID))
	d.Set("name", executeSchedule.Schedule.Name)
	d.Set("description", executeSchedule.Schedule.Description)
	d.Set("enabled", executeSchedule.Schedule.Enabled)
	d.Set("time_zone", executeSchedule.Schedule.Scheduletimezone)
	d.Set("schedule", executeSchedule.Schedule.Cron)

	return diags
}

func resourceExecuteScheduleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	schedule := make(map[string]interface{})

	schedule["name"] = d.Get("name").(string)
	schedule["description"] = d.Get("description").(string)
	schedule["enabled"] = d.Get("enabled").(bool)
	schedule["scheduleType"] = "execute"
	schedule["scheduleTimezone"] = d.Get("time_zone").(string)
	schedule["cron"] = d.Get("schedule").(string)

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"schedule": schedule,
		},
	}

	resp, err := client.UpdateExecuteSchedule(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateExecuteScheduleResult)
	executeSchedule := result.ExecuteSchedule

	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(executeSchedule.ID))
	return resourceExecuteScheduleRead(ctx, d, meta)
}

func resourceExecuteScheduleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteExecuteSchedule(toInt64(id), req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
			return diag.FromErr(err)
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)
			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)
	d.SetId("")
	return diags
}

type ExecuteSchedule struct {
	Schedule struct {
		ID               int    `json:"id"`
		Name             string `json:"name"`
		Description      string `json:"description"`
		Enabled          bool   `json:"enabled"`
		Scheduletype     string `json:"scheduleType"`
		Scheduletimezone string `json:"scheduleTimezone"`
		Cron             string `json:"cron"`
		Datecreated      string `json:"dateCreated"`
		Lastupdated      string `json:"lastUpdated"`
	} `json:"schedule"`
}
