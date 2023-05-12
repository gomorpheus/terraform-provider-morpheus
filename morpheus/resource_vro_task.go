package morpheus

import (
	"context"
	"encoding/json"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceVrealizeOrchestratorTask() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus vRealize Orchestrator (vRO) task resource",
		CreateContext: resourceVrealizeOrchestratorTaskCreate,
		ReadContext:   resourceVrealizeOrchestratorTaskRead,
		UpdateContext: resourceVrealizeOrchestratorTaskUpdate,
		DeleteContext: resourceVrealizeOrchestratorTaskDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the vRO workflow task",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the vRO workflow task",
				Required:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the vRO workflow task",
				Optional:    true,
				Computed:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "The organization labels associated with the task (Only supported on Morpheus 5.5.3 or higher)",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"result_type": {
				Type:         schema.TypeString,
				Description:  "The expected result type (value, keyValue, json)",
				ValidateFunc: validation.StringInSlice([]string{"value", "keyValue", "json"}, false),
				Optional:     true,
				Computed:     true,
			},
			"vro_integration_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the vRO integration",
				Required:    true,
			},
			"vro_workflow_value": {
				Type:        schema.TypeInt,
				Description: "The value of the vRO workflow",
				Required:    true,
			},
			"body": {
				Type:             schema.TypeString,
				Description:      "The JSON body to send to vRO",
				Optional:         true,
				DiffSuppressFunc: suppressEquivalentJsonDiffs,
			},
			"execute_target": {
				Type:        schema.TypeString,
				Description: "The target that the vRO workflow will be executed on",
				Optional:    true,
			},
			"retryable": {
				Type:        schema.TypeBool,
				Description: "Whether to retry the task if there is a failure",
				Optional:    true,
				Default:     false,
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
				Description: "Custom configuration data to pass during the execution of the vRO workflow task",
				Optional:    true,
				Default:     false,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceVrealizeOrchestratorTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)

	taskType := make(map[string]interface{})
	taskType["code"] = "vro"

	taskOptions := make(map[string]interface{})
	taskOptions["vroIntegrationId"] = d.Get("vro_integration_id")
	taskOptions["vroWorkflow"] = d.Get("vro_workflow_value")
	taskOptions["vroBody"] = d.Get("body")

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		for _, s := range attr.(*schema.Set).List() {
			labelsPayload = append(labelsPayload, s.(string))
		}
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"task": map[string]interface{}{
				"name":              name,
				"code":              d.Get("code").(string),
				"labels":            labelsPayload,
				"taskType":          taskType,
				"taskOptions":       taskOptions,
				"resultType":        d.Get("result_type"),
				"executeTarget":     d.Get("execute_target").(string),
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

	resourceVrealizeOrchestratorTaskRead(ctx, d, meta)
	return diags
}

func resourceVrealizeOrchestratorTaskRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	result := resp.Result.(*morpheus.GetTaskResult)
	workflowTask := result.Task

	json.Unmarshal(resp.Body, &workflowTask)
	d.SetId(int64ToString(workflowTask.ID))
	d.Set("name", workflowTask.Name)
	d.Set("code", workflowTask.Code)
	d.Set("labels", workflowTask.Labels)
	d.Set("result_type", workflowTask.ResultType)
	d.Set("vro_integration_id", workflowTask.TaskOptions.VroIntegrationId)
	d.Set("vro_workflow_value", workflowTask.TaskOptions.VroWorkflow)
	d.Set("body", workflowTask.TaskOptions.VroBody)
	d.Set("execute_target", workflowTask.ExecuteTarget)
	d.Set("retryable", workflowTask.Retryable)
	d.Set("retry_count", workflowTask.RetryCount)
	d.Set("retry_delay_seconds", workflowTask.RetryDelaySeconds)
	d.Set("allow_custom_config", workflowTask.AllowCustomConfig)
	return diags
}

func resourceVrealizeOrchestratorTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	name := d.Get("name").(string)

	taskType := make(map[string]interface{})
	taskType["code"] = "vro"

	taskOptions := make(map[string]interface{})
	if d.HasChange("vro_integration_id") {
		taskOptions["vroIntegrationId"] = d.Get("vro_integration_id")
	}
	if d.HasChange("vro_workflow_value") {
		taskOptions["vroWorkflow"] = d.Get("vro_workflow_value")
	}
	if d.HasChange("body") {
		taskOptions["vroBody"] = d.Get("body")
	}

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		for _, s := range attr.(*schema.Set).List() {
			labelsPayload = append(labelsPayload, s.(string))
		}
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"task": map[string]interface{}{
				"name":              name,
				"code":              d.Get("code").(string),
				"labels":            labelsPayload,
				"taskType":          taskType,
				"taskOptions":       taskOptions,
				"resultType":        d.Get("result_type"),
				"executeTarget":     d.Get("execute_target").(string),
				"retryable":         d.Get("retryable"),
				"retryCount":        d.Get("retry_count"),
				"retryDelaySeconds": d.Get("retry_delay_seconds"),
				"allowCustomConfig": d.Get("allow_custom_config"),
			},
		},
	}
	log.Printf("API REQUEST: %s", req)
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
	return resourceVrealizeOrchestratorTaskRead(ctx, d, meta)
}

func resourceVrealizeOrchestratorTaskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
