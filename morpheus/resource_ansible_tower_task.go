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

func resourceAnsibleTowerTask() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus vRealize Orchestrator (vRO) task resource",
		CreateContext: resourceAnsibleTowerTaskCreate,
		ReadContext:   resourceAnsibleTowerTaskRead,
		UpdateContext: resourceAnsibleTowerTaskUpdate,
		DeleteContext: resourceAnsibleTowerTaskDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the ansible tower task",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the ansible tower task",
				Required:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the ansible tower task",
				Optional:    true,
				Computed:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "The organization labels associated with the ansible tower task (Only supported on Morpheus 5.5.3 or higher)",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"ansible_tower_integration_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the ansible tower integration",
				Required:    true,
			},
			"ansible_tower_inventory_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the ansible tower inventory",
				Required:    true,
			},
			"group": {
				Type:        schema.TypeString,
				Description: "The name of a new or existing group in the inventory",
				Optional:    true,
				Computed:    true,
			},
			"job_template_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the ansible tower job template",
				Required:    true,
			},
			"scm_override": {
				Type:        schema.TypeString,
				Description: "The git reference override",
				Optional:    true,
				Computed:    true,
			},
			"execute_mode": {
				Type:         schema.TypeString,
				Description:  "The ansible tower execution mode (executeHost, executeGroup, executeAll, off)",
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"executeHost", "executeGroup", "executeAll", "off"}, false),
			},
			"execute_target": {
				Type:         schema.TypeString,
				Description:  "The target that the ansible tower job will be executed on (local, remote, resource)",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"local", "remote", "resource"}, false),
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
				Description: "Custom configuration data to pass during the execution of the ansible tower job task",
				Optional:    true,
				Default:     false,
			},
			"visibility": {
				Type:         schema.TypeString,
				Description:  "The visibility of the ansible tower task (public or private)",
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"public", "private"}, false),
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceAnsibleTowerTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)

	taskType := make(map[string]interface{})
	taskType["code"] = "ansibleTowerTask"

	taskOptions := make(map[string]interface{})
	taskOptions["ansibleTowerIntegrationId"] = d.Get("ansible_tower_integration_id")
	taskOptions["ansibleTowerInventoryId"] = d.Get("ansible_tower_inventory_id")
	taskOptions["ansibleGroup"] = d.Get("group")
	taskOptions["ansibleTowerJobTemplateId"] = d.Get("job_template_id")
	taskOptions["ansibleTowerExecuteMode"] = d.Get("execute_mode")
	taskOptions["ansibleTowerGitRef"] = d.Get("scm_override")

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
				"visibility":        d.Get("visibility").(string),
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

	resourceAnsibleTowerTaskRead(ctx, d, meta)
	return diags
}

func resourceAnsibleTowerTaskRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	result := resp.Result.(*morpheus.GetTaskResult)
	ansibleTowerTask := result.Task

	d.SetId(int64ToString(ansibleTowerTask.ID))
	d.Set("name", ansibleTowerTask.Name)
	d.Set("code", ansibleTowerTask.Code)
	d.Set("labels", ansibleTowerTask.Labels)
	integrationId, err := strconv.Atoi(ansibleTowerTask.TaskOptions.AnsibleTowerIntegrationId)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("ansible_tower_integration_id", integrationId)
	inventoryId, err := strconv.Atoi(ansibleTowerTask.TaskOptions.AnsibleTowerInventoryId)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("ansible_tower_inventory_id", inventoryId)
	d.Set("group", ansibleTowerTask.TaskOptions.AnsibleGroup)
	d.Set("scm_override", ansibleTowerTask.TaskOptions.AnsibleTowerGitRef)
	d.Set("execute_mode", ansibleTowerTask.TaskOptions.AnsibleTowerExecuteMode)
	jobTemplateId, err := strconv.Atoi(ansibleTowerTask.TaskOptions.AnsibleTowerJobTemplateId)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("job_template_id", jobTemplateId)
	d.Set("execute_target", ansibleTowerTask.ExecuteTarget)
	d.Set("retryable", ansibleTowerTask.Retryable)
	d.Set("retry_count", ansibleTowerTask.RetryCount)
	d.Set("retry_delay_seconds", ansibleTowerTask.RetryDelaySeconds)
	d.Set("allow_custom_config", ansibleTowerTask.AllowCustomConfig)
	d.Set("visibility", ansibleTowerTask.Visibility)
	return diags
}

func resourceAnsibleTowerTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	name := d.Get("name").(string)

	taskType := make(map[string]interface{})
	taskType["code"] = "ansibleTowerTask"

	taskOptions := make(map[string]interface{})
	if d.HasChange("ansible_tower_integration_id") {
		taskOptions["ansibleTowerIntegrationId"] = d.Get("ansible_tower_integration_id")
	}
	if d.HasChange("ansible_tower_inventory_id") {
		taskOptions["ansibleTowerInventoryId"] = d.Get("ansible_tower_inventory_id")
	}
	if d.HasChange("group") {
		taskOptions["ansibleGroup"] = d.Get("group")
	}
	if d.HasChange("job_template_id") {
		taskOptions["ansibleTowerJobTemplateId"] = d.Get("job_template_id")
	}
	if d.HasChange("execute_mode") {
		taskOptions["ansibleTowerExecuteMode"] = d.Get("execute_mode")
	}
	if d.HasChange("scm_override") {
		taskOptions["ansibleTowerGitRef"] = d.Get("scm_override")
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
	return resourceAnsibleTowerTaskRead(ctx, d, meta)
}

func resourceAnsibleTowerTaskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
