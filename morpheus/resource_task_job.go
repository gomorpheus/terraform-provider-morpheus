package morpheus

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceTaskJob() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a task job resource",
		CreateContext: resourceTaskJobCreate,
		ReadContext:   resourceTaskJobRead,
		UpdateContext: resourceTaskJobUpdate,
		DeleteContext: resourceTaskJobDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the task job",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the task job",
				Required:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the task job is enabled",
				Optional:    true,
				Default:     true,
			},
			"task_id": {
				Type:        schema.TypeInt,
				Description: "The id of the task associated with the job",
				Optional:    true,
				Computed:    true,
			},
			"schedule_mode": {
				Type:         schema.TypeString,
				Description:  "The job scheduling type (manual, date_and_time, scheduled)",
				ValidateFunc: validation.StringInSlice([]string{"manual", "date_and_time", "scheduled"}, false),
				Required:     true,
			},
			"scheduled_date_and_time": {
				Type:          schema.TypeString,
				Description:   "The date and time the job will be executed if schedule mode date_and_time is used",
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"execution_schedule_id"},
			},
			"execution_schedule_id": {
				Type:        schema.TypeInt,
				Description: "The id of the execution schedule associated with the job",
				Optional:    true,
				Computed:    true,
			},
			"context_type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"appliance", "server", "instance"}, false),
				Description:  "The context that the job should run as (appliance, server, instance)",
				Required:     true,
			},
			"server_ids": {
				Type:          schema.TypeList,
				Description:   "A list of server ids to associate with the job",
				Optional:      true,
				Computed:      true,
				Elem:          &schema.Schema{Type: schema.TypeInt},
				ConflictsWith: []string{"instance_ids"},
			},
			"instance_ids": {
				Type:          schema.TypeList,
				Description:   "A list of instance ids to associate with the job",
				Optional:      true,
				Computed:      true,
				Elem:          &schema.Schema{Type: schema.TypeInt},
				ConflictsWith: []string{"server_ids"},
			},
			"custom_config": {
				Type:        schema.TypeString,
				Description: "The task custom configuration",
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaskJobCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	job := make(map[string]interface{})

	job["name"] = d.Get("name").(string)
	job["enabled"] = d.Get("enabled").(bool)
	job["task"] = map[string]int{
		"id": d.Get("task_id").(int),
	}
	job["targetType"] = d.Get("context_type").(string)
	job["customConfig"] = d.Get("custom_config")

	// Evaluate different schedululing methods
	switch d.Get("schedule_mode") {
	case "manual":
		job["scheduleMode"] = "manual"
	case "date_and_time":
		job["scheduleMode"] = "dateTime"
		job["dateTime"] = d.Get("scheduled_date_and_time").(string)
	case "scheduled":
		job["scheduleMode"] = d.Get("execution_schedule_id")
	}

	// instances
	var targets []map[string]interface{}
	if d.Get("context_type") == "instance" {
		instanceList := d.Get("instance_ids").([]interface{})
		// iterate over the array of instance ids
		for i := 0; i < len(instanceList); i++ {
			row := make(map[string]interface{})
			row["refId"] = instanceList[i]
			targets = append(targets, row)
		}
	}

	// servers
	if d.Get("context_type") == "server" {
		serverList := d.Get("server_ids").([]interface{})
		// iterate over the array of server ids
		for i := 0; i < len(serverList); i++ {
			row := make(map[string]interface{})
			row["refId"] = serverList[i]
			targets = append(targets, row)
		}
	}

	job["targets"] = targets

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"job": job,
		},
	}
	resp, err := client.CreateJob(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		log.Fatal(err)
	}

	jobId := fmt.Sprintf("%v", result["id"])
	// Successfully created resource, now set id
	d.SetId(jobId)

	resourceTaskJobRead(ctx, d, meta)
	return diags
}

func resourceTaskJobRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindJobByName(name)
	} else if id != "" {
		resp, err = client.GetJob(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Execute schedule cannot be read without name or id")
	}

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

	// store resource data
	var taskJob TaskJob
	json.Unmarshal(resp.Body, &taskJob)

	d.SetId(intToString(taskJob.Job.ID))
	d.Set("name", taskJob.Job.Name)
	d.Set("enabled", taskJob.Job.Enabled)
	d.Set("task_id", taskJob.Job.Task.ID)

	switch taskJob.Job.Schedulemode {
	case "manual":
		d.Set("schedule_mode", "manual")
	case "dateTime":
		d.Set("schedule_mode", "date_and_time")
	case "scheduled":
		d.Set("schedule_mode", "scheduled")
	}
	d.Set("scheduled_date_and_time", taskJob.Job.Datetime)
	d.Set("custom_config", taskJob.Job.Customconfig)

	// instances
	var instanceIds []int64
	if taskJob.Job.Targets != nil {
		// iterate over the array of targets
		for i := 0; i < len(taskJob.Job.Targets); i++ {
			instance := taskJob.Job.Targets[i]
			instanceIds = append(instanceIds, int64(instance.Refid))
		}
	}
	d.Set("instance_ids", instanceIds)

	// servers
	var serverIds []int64
	if taskJob.Job.Targets != nil {
		// iterate over the array of targets
		for i := 0; i < len(taskJob.Job.Targets); i++ {
			server := taskJob.Job.Targets[i]
			serverIds = append(serverIds, int64(server.Refid))
		}
	}
	d.Set("server_ids", serverIds)

	return diags
}

func resourceTaskJobUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	job := make(map[string]interface{})

	job["name"] = d.Get("name").(string)
	job["enabled"] = d.Get("enabled").(bool)
	job["task"] = map[string]int{
		"id": d.Get("task_id").(int),
	}
	job["targetType"] = d.Get("context_type").(string)
	job["customConfig"] = d.Get("custom_config")

	// Evaluate different schedululing methods
	switch d.Get("schedule_mode") {
	case "manual":
		job["scheduleMode"] = "manual"
	case "date_and_time":
		job["scheduleMode"] = "dateTime"
		job["dateTime"] = d.Get("scheduled_date_and_time").(string)
	case "scheduled":
		job["scheduleMode"] = d.Get("execution_schedule_id")
	}

	// instances
	var targets []map[string]interface{}
	if d.Get("context_type") == "instance" {
		instanceList := d.Get("instance_ids").([]interface{})
		// iterate over the array of instance ids
		for i := 0; i < len(instanceList); i++ {
			row := make(map[string]interface{})
			row["refId"] = instanceList[i]
			targets = append(targets, row)
		}
	}

	// servers
	if d.Get("context_type") == "server" {
		serverList := d.Get("server_ids").([]interface{})
		// iterate over the array of instance ids
		for i := 0; i < len(serverList); i++ {
			row := make(map[string]interface{})
			row["refId"] = serverList[i]
			targets = append(targets, row)
		}
	}

	job["targets"] = targets

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"job": job,
		},
	}

	resp, err := client.UpdateJob(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		log.Fatal(err)
	}

	// Successfully updated resource, now set id
	// the task API doesn't return the id so setting
	// to the original id
	d.SetId(id)
	return resourceTaskJobRead(ctx, d, meta)
}

func resourceTaskJobDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteJob(toInt64(id), req)
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

type TaskJob struct {
	Job struct {
		Category  interface{} `json:"category"`
		Createdby struct {
			Displayname string `json:"displayName"`
			ID          int    `json:"id"`
			Username    string `json:"username"`
		} `json:"createdBy"`
		Customconfig  string      `json:"customConfig"`
		Customoptions interface{} `json:"customOptions"`
		Datecreated   time.Time   `json:"dateCreated"`
		Datetime      interface{} `json:"dateTime"`
		Description   interface{} `json:"description"`
		Enabled       bool        `json:"enabled"`
		ID            int         `json:"id"`
		Jobsummary    string      `json:"jobSummary"`
		Lastresult    string      `json:"lastResult"`
		Lastrun       time.Time   `json:"lastRun"`
		Lastupdated   time.Time   `json:"lastUpdated"`
		Name          string      `json:"name"`
		Namespace     interface{} `json:"namespace"`
		Schedulemode  string      `json:"scheduleMode"`
		Status        interface{} `json:"status"`
		Targettype    string      `json:"targetType"`
		Targets       []struct {
			ID         int    `json:"id"`
			Name       string `json:"name"`
			Refid      int    `json:"refId"`
			Targettype string `json:"targetType"`
		} `json:"targets"`
		Task struct {
			ID int `json:"id"`
		} `json:"task"`
		Type struct {
			Code string `json:"code"`
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"type"`
	} `json:"job"`
}
