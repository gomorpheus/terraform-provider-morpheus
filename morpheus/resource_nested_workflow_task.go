package morpheus

import (
	"context"
	"encoding/json"
	"time"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNestedWorkflowTask() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus nested workflow task resource",
		CreateContext: resourceNestedWorkflowTaskCreate,
		ReadContext:   resourceNestedWorkflowTaskRead,
		UpdateContext: resourceNestedWorkflowTaskUpdate,
		DeleteContext: resourceNestedWorkflowTaskDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the nested workflow task",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the nested workflow task",
				Required:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the nested workflow task",
				Optional:    true,
				Computed:    true,
			},
			"operational_workflow_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the operational workflow",
				Optional:    true,
				Computed:    true,
			},
			"operational_workflow_name": {
				Type:        schema.TypeString,
				Description: "The name of the operational workflow",
				Optional:    true,
				Computed:    true,
			},
			"retryable": {
				Type:        schema.TypeBool,
				Description: "Whether to retry the task if there is a failure",
				Optional:    true,
				Computed:    true,
			},
			"retry_count": {
				Type:        schema.TypeInt,
				Description: "The number of times to retry the task if there is a failure",
				Optional:    true,
				Default:     5,
			},
			"retry_delay_seconds": {
				Type:        schema.TypeInt,
				Description: "The number of seconds to wait between retry attempts",
				Optional:    true,
				Default:     10,
			},
			"allow_custom_config": {
				Type:        schema.TypeBool,
				Description: "Custom configuration data to pass during the execution of the shell script",
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceNestedWorkflowTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)

	taskType := map[string]interface{}{
		"code": "nestedWorkflow",
	}

	taskOptions := map[string]interface{}{
		"operationalWorkflowId":   d.Get("operational_workflow_id").(int),
		"operationalWorkflowName": d.Get("operational_workflow_name").(string),
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"task": map[string]interface{}{
				"name":              name,
				"code":              d.Get("code").(string),
				"taskType":          taskType,
				"taskOptions":       taskOptions,
				"executeTarget":     "local",
				"retryable":         d.Get("retryable"),
				"retryCount":        d.Get("retry_count"),
				"retryDelaySeconds": d.Get("retry_delay_seconds"),
				"allowCustomConfig": d.Get("allow_custom_config"),
			},
		},
	}
	resp, err := client.CreateTask(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateTaskResult)
	task := result.Task
	// Successfully created resource, now set id
	d.SetId(int64ToString(task.ID))
	log.Printf("Task ID: %s", int64ToString(task.ID))

	resourceShellScriptTaskRead(ctx, d, meta)
	return diags
}

func resourceNestedWorkflowTaskRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindTaskByName(name)
	} else if id != "" {
		resp, err = client.GetTask(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Task cannot be read without name or id")
	}

	if err != nil {
		// 404 is ok?
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
	var nestedWorkflowTask NestedWorkflow
	json.Unmarshal(resp.Body, &nestedWorkflowTask)
	d.SetId(intToString(nestedWorkflowTask.Task.ID))
	d.Set("name", nestedWorkflowTask.Task.Name)
	d.Set("code", nestedWorkflowTask.Task.Code)
	d.Set("execute_target", nestedWorkflowTask.Task.Executetarget)
	d.Set("operational_workflow_id", nestedWorkflowTask.Task.Taskoptions.OperationalWorkflowId)
	d.Set("operational_workflow_name", nestedWorkflowTask.Task.Taskoptions.OperationalWorkflowName)
	d.Set("retryable", nestedWorkflowTask.Task.Retryable)
	d.Set("retry_count", nestedWorkflowTask.Task.Retrycount)
	d.Set("retry_delay_seconds", nestedWorkflowTask.Task.Retrydelayseconds)
	d.Set("allow_custom_config", nestedWorkflowTask.Task.Allowcustomconfig)
	return diags
}

func resourceNestedWorkflowTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	name := d.Get("name").(string)

	taskType := map[string]interface{}{
		"code": "nestedWorkflow",
	}

	taskOptions := map[string]interface{}{
		"operationalWorkflowId":   d.Get("operational_workflow_id").(int),
		"operationalWorkflowName": d.Get("operational_workflow_name").(string),
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"task": map[string]interface{}{
				"name":              name,
				"code":              d.Get("code").(string),
				"taskType":          taskType,
				"taskOptions":       taskOptions,
				"executeTarget":     "local",
				"retryable":         d.Get("retryable"),
				"retryCount":        d.Get("retry_count"),
				"retryDelaySeconds": d.Get("retry_delay_seconds"),
				"allowCustomConfig": d.Get("allow_custom_config"),
			},
		},
	}
	resp, err := client.UpdateTask(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateTaskResult)
	shellScriptTask := result.Task
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(shellScriptTask.ID))
	return resourceShellScriptTaskRead(ctx, d, meta)
}

func resourceNestedWorkflowTaskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteTask(toInt64(id), req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
			return nil
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)
			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)
	d.SetId("")
	return diags
}

type NestedWorkflow struct {
	Task struct {
		ID        int    `json:"id"`
		Accountid int    `json:"accountId"`
		Name      string `json:"name"`
		Code      string `json:"code"`
		Tasktype  struct {
			ID   int    `json:"id"`
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"taskType"`
		Taskoptions struct {
			OperationalWorkflowId   int    `json:"operational_workflow_id"`
			OperationalWorkflowName string `json:"operational_workflow_name"`
		}
		Executetarget     string    `json:"executeTarget"`
		Retryable         bool      `json:"retryable"`
		Retrycount        int       `json:"retryCount"`
		Retrydelayseconds int       `json:"retryDelaySeconds"`
		Allowcustomconfig bool      `json:"allowCustomConfig"`
		Datecreated       time.Time `json:"dateCreated"`
		Lastupdated       time.Time `json:"lastUpdated"`
	} `json:"task"`
}
