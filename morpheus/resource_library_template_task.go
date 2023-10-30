package morpheus

import (
	"context"
	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceLibraryTemplateTask() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus library template task resource",
		CreateContext: resourceLibraryTemplateTaskCreate,
		ReadContext:   resourceLibraryTemplateTaskRead,
		UpdateContext: resourceLibraryTemplateTaskUpdate,
		DeleteContext: resourceLibraryTemplateTaskDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the library template task",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the library template task",
				Required:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the library template task",
				Optional:    true,
				Computed:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "The organization labels associated with the library template task (Only supported on Morpheus 5.5.3 or higher)",
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
			"file_template": {
				Type:        schema.TypeString,
				Description: "The library file template in Morpheus",
				Optional:    true,
				Computed:    true,
			},
			"file_template_id": {
				Type:        schema.TypeString,
				Description: "The library file template id in Morpheus",
				Optional:    true,
				Computed:    true,
			},
			"execute_target": {
				Type:         schema.TypeString,
				Description:  "The target for the library template",
				ValidateFunc: validation.StringInSlice([]string{"resource"}, false),
				Optional:     true,
				Computed:     true,
			},
			"retryable": {
				Type:        schema.TypeBool,
				Description: "Whether to retry the library task if there is a failure",
				Optional:    true,
				Computed:    true,
			},
			"retry_count": {
				Type:        schema.TypeInt,
				Description: "The number of times to retry the library task if there is a failure",
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
				Description: "Custom configuration data to pass during the execution of the library template",
				Optional:    true,
				Computed:    true,
			},
			"visibility": {
				Type:         schema.TypeString,
				Description:  "The visibility of the task (private or public)",
				ValidateFunc: validation.StringInSlice([]string{"private", "public"}, false),
				Optional:     true,
				Computed:     true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceLibraryTemplateTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)

	taskType := make(map[string]interface{})
	taskType["code"] = "containerTemplate"

	taskOptions := make(map[string]interface{})
	if d.Get("visibility") != "" {
		taskOptions["visibility"] = d.Get("visibility")
	}

	if d.Get("file_template_id") != "" {
		taskOptions["containerTemplateId"] = d.Get("file_template_id")
	}
	if d.Get("file_template") != "" {
		taskOptions["containerTemplate"] = d.Get("file_template")
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
				"visibility":        d.Get("visibility"),
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

	resourceLibraryTemplateTaskRead(ctx, d, meta)
	return diags
}

func resourceLibraryTemplateTaskRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	libraryTemplateTask := result.Task

	d.SetId(int64ToString(libraryTemplateTask.ID))
	d.Set("name", libraryTemplateTask.Name)
	d.Set("code", libraryTemplateTask.Code)
	d.Set("labels", libraryTemplateTask.Labels)
	d.Set("result_type", libraryTemplateTask.ResultType)
	d.Set("file_template", libraryTemplateTask.TaskOptions.ContainerTemplate)
	d.Set("file_template_id", libraryTemplateTask.TaskOptions.ContainerTemplateId)
	d.Set("execute_target", libraryTemplateTask.ExecuteTarget)
	d.Set("retryable", libraryTemplateTask.Retryable)
	d.Set("retry_count", libraryTemplateTask.RetryCount)
	d.Set("retry_delay_seconds", libraryTemplateTask.RetryDelaySeconds)
	d.Set("allow_custom_config", libraryTemplateTask.AllowCustomConfig)
	d.Set("visibility", libraryTemplateTask.Visibility)
	return diags
}

func resourceLibraryTemplateTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	name := d.Get("name").(string)

	taskType := make(map[string]interface{})
	taskType["code"] = "containerTemplate"

	taskOptions := make(map[string]interface{})
	if d.Get("visibility") != "" {
		taskOptions["visibility"] = d.Get("visibility")
	}

	if d.Get("file_template_id") != "" {
		taskOptions["containerTemplateId"] = d.Get("file_template_id")
	}
	if d.Get("file_template") != "" {
		taskOptions["containerTemplate"] = d.Get("file_template")
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
				"visibility":        d.Get("visibility"),
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
	libraryTemplateTask := result.Task
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(libraryTemplateTask.ID))
	return resourceLibraryTemplateTaskRead(ctx, d, meta)
}

func resourceLibraryTemplateTaskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
