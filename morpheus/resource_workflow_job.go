package morpheus

import (
	"context"
	"strconv"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceWorkflowJob() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a workflow job resource",
		CreateContext: resourceWorkflowJobCreate,
		ReadContext:   resourceWorkflowJobRead,
		UpdateContext: resourceWorkflowJobUpdate,
		DeleteContext: resourceWorkflowJobDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the workflow job",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the workflow job",
				Required:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "The organization labels associated with the workflow job (Only supported on Morpheus 5.5.3 or higher)",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the workflow job is enabled",
				Optional:    true,
				Default:     true,
			},
			"workflow_id": {
				Type:        schema.TypeInt,
				Description: "The id of the workflow associated with the job",
				Required:    true,
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
				ConflictsWith: []string{"execution_schedule_id"},
			},
			"execution_schedule_id": {
				Type:        schema.TypeInt,
				Description: "The id of the execution schedule associated with the job",
				Optional:    true,
			},
			"context_type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"appliance", "server", "instance", "instance-label", "server-label"}, false),
				Description:  "The context that the job should run as (appliance, server, instance, instance-label, server-label)",
				Required:     true,
			},
			"server_ids": {
				Type:          schema.TypeList,
				Description:   "A list of server ids to associate with the job",
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeInt},
				ConflictsWith: []string{"instance_ids", "instance_label", "server_label"},
			},
			"server_label": {
				Type:          schema.TypeString,
				Description:   "The server label used for dynamic automation targeting",
				Optional:      true,
				ConflictsWith: []string{"instance_ids", "server_ids", "instance_label"},
			},
			"instance_ids": {
				Type:          schema.TypeList,
				Description:   "A list of instance ids to associate with the job",
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeInt},
				ConflictsWith: []string{"server_ids", "instance_label", "server_label"},
			},
			"instance_label": {
				Type:          schema.TypeString,
				Description:   "The instance label used for dynamic automation targeting",
				Optional:      true,
				ConflictsWith: []string{"instance_ids", "server_ids", "server_label"},
			},
			"custom_options": {
				Description: "Custom options to pass to the workflow",
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceWorkflowJobCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	job := make(map[string]interface{})

	job["name"] = d.Get("name").(string)
	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		for _, s := range attr.(*schema.Set).List() {
			labelsPayload = append(labelsPayload, s.(string))
		}
	}
	job["labels"] = labelsPayload
	job["enabled"] = d.Get("enabled").(bool)
	job["workflow"] = map[string]int{
		"id": d.Get("workflow_id").(int),
	}

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

	job["targetType"] = d.Get("context_type").(string)
	if d.Get("context_type").(string) == "instance-label" {
		job["instanceLabel"] = d.Get("instance_label").(string)
	}
	if d.Get("context_type").(string) == "server-label" {
		job["serverLabel"] = d.Get("server_label").(string)
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

	// Custom Options
	if d.Get("custom_options") != nil {
		customOptionsInput := d.Get("custom_options").(map[string]interface{})
		customOptions := make(map[string]interface{})
		for key, value := range customOptionsInput {
			customOptions[key] = value.(string)
		}
		job["customOptions"] = customOptions
	}

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

	result := resp.Result.(*morpheus.CreateJobResult)
	jobResult := result.Job

	// Successfully created resource, now set id
	d.SetId(int64ToString(jobResult.ID))

	resourceWorkflowJobRead(ctx, d, meta)
	return diags
}

func resourceWorkflowJobRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	result := resp.Result.(*morpheus.GetJobResult)
	workflowJob := result.Job

	d.SetId(int64ToString(workflowJob.ID))
	d.Set("name", workflowJob.Name)
	if len(workflowJob.Labels) > 0 {
		d.Set("labels", workflowJob.Labels)
	}
	d.Set("enabled", workflowJob.Enabled)
	d.Set("workflow_id", workflowJob.Workflow.ID)
	d.Set("context_type", workflowJob.TargetType)
	switch workflowJob.ScheduleMode {
	case "manual":
		d.Set("schedule_mode", "manual")
	case "dateTime":
		d.Set("schedule_mode", "date_and_time")
		d.Set("scheduled_date_and_time", workflowJob.DateTime)
		// Execute schedule uses the schedule mode field for storing the execute schedule id
	default:
		d.Set("schedule_mode", "scheduled")
		intVar, err := strconv.Atoi(workflowJob.ScheduleMode)
		if err != nil {
			log.Printf("String Conversion Failure: %s", err)
		}
		d.Set("execution_schedule_id", intVar)
	}

	switch workflowJob.TargetType {
	case "instance":
		// instances
		var instanceIds []int64
		if workflowJob.Targets != nil {
			// iterate over the array of targets
			for i := 0; i < len(workflowJob.Targets); i++ {
				instance := workflowJob.Targets[i]
				instanceIds = append(instanceIds, int64(instance.RefId))
			}
		}
		d.Set("instance_ids", instanceIds)
	case "server":
		// servers
		var serverIds []int64
		if workflowJob.Targets != nil {
			// iterate over the array of targets
			for i := 0; i < len(workflowJob.Targets); i++ {
				server := workflowJob.Targets[i]
				serverIds = append(serverIds, int64(server.RefId))
			}
		}
		d.Set("server_ids", serverIds)
	case "instance-label":
		d.Set("instance_label", workflowJob.Targets[0].Name)
	case "server-label":
		d.Set("server_label", workflowJob.Targets[0].Name)
	}
	d.Set("custom_options", workflowJob.CustomOptions)

	return diags
}

func resourceWorkflowJobUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	job := make(map[string]interface{})

	job["name"] = d.Get("name").(string)
	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		for _, s := range attr.(*schema.Set).List() {
			labelsPayload = append(labelsPayload, s.(string))
		}
	}
	job["labels"] = labelsPayload
	job["enabled"] = d.Get("enabled").(bool)
	job["workflow"] = map[string]int{
		"id": d.Get("workflow_id").(int),
	}
	job["targetType"] = d.Get("context_type").(string)
	if d.Get("context_type").(string) == "instance-label" {
		job["instanceLabel"] = d.Get("instance_label").(string)
	}
	if d.Get("context_type").(string) == "server-label" {
		job["serverLabel"] = d.Get("server_label").(string)
	}

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

	// Custom Options
	if d.Get("custom_options") != nil {
		customOptionsInput := d.Get("custom_options").(map[string]interface{})
		customOptions := make(map[string]interface{})
		for key, value := range customOptionsInput {
			customOptions[key] = value.(string)
		}
		job["customOptions"] = customOptions
	}

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
	result := resp.Result.(*morpheus.UpdateJobResult)
	jobResult := result.Job

	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(jobResult.ID))
	return resourceWorkflowJobRead(ctx, d, meta)
}

func resourceWorkflowJobDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
