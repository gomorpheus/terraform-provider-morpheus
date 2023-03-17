package morpheus

import (
	"context"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceRubyScriptTask() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus ruby script task resource",
		CreateContext: resourceRubyScriptTaskCreate,
		ReadContext:   resourceRubyScriptTaskRead,
		UpdateContext: resourceRubyScriptTaskUpdate,
		DeleteContext: resourceRubyScriptTaskDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the ruby script task",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the ruby script task",
				Required:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the ruby script task",
				Optional:    true,
			},
			"result_type": {
				Type:         schema.TypeString,
				Description:  "The expected result type (single value, key pairs, json)",
				ValidateFunc: validation.StringInSlice([]string{"value", "keyValue", "json"}, false),
				Optional:     true,
			},
			"source_type": {
				Type:         schema.TypeString,
				Description:  "The source of the ruby script (local, url or repository)",
				ValidateFunc: validation.StringInSlice([]string{"local", "url", "repository"}, false),
				Required:     true,
			},
			"script_content": {
				Type:        schema.TypeString,
				Description: "The content of the ruby script. Used when the local source type is specified",
				Optional:    true,
			},
			"script_path": {
				Type:        schema.TypeString,
				Description: "The path of the ruby script, either the url or the path in the repository",
				Optional:    true,
			},
			"repository_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the git repository integration",
				Optional:    true,
			},
			"version_ref": {
				Type:        schema.TypeString,
				Description: "The git reference of the repository to pull (main, master, etc.)",
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
				Description: "Custom configuration data to pass during the execution of the ruby script",
				Optional:    true,
				Default:     false,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceRubyScriptTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)

	sourceOptions := make(map[string]interface{})
	if d.Get("script_content") != "" {
		sourceOptions["content"] = d.Get("script_content")
	}
	if d.Get("script_path") != "" {
		sourceOptions["contentPath"] = d.Get("script_path")
	}
	sourceOptions["contentRef"] = d.Get("version_ref")
	sourceOptions["repository"] = map[string]interface{}{
		"id": d.Get("repository_id"),
	}
	sourceOptions["sourceType"] = d.Get("source_type")

	taskType := make(map[string]interface{})
	taskType["code"] = "jrubyTask"

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"task": map[string]interface{}{
				"name":              name,
				"code":              d.Get("code").(string),
				"file":              sourceOptions,
				"taskType":          taskType,
				"resultType":        d.Get("result_type"),
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

	resourceRubyScriptTaskRead(ctx, d, meta)
	return diags
}

func resourceRubyScriptTaskRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	rubyScriptTask := result.Task
	d.SetId(int64ToString(rubyScriptTask.ID))
	d.Set("name", rubyScriptTask.Name)
	d.Set("code", rubyScriptTask.Code)
	d.Set("result_type", rubyScriptTask.ResultType)
	d.Set("source_type", rubyScriptTask.File.SourceType)
	d.Set("script_content", rubyScriptTask.File.Content)
	d.Set("script_path", rubyScriptTask.File.ContentPath)
	d.Set("version_ref", rubyScriptTask.File.ContentRef)
	d.Set("repository_id", rubyScriptTask.File.Repository.ID)
	d.Set("retryable", rubyScriptTask.Retryable)
	d.Set("retry_count", rubyScriptTask.RetryCount)
	d.Set("retry_delay_seconds", rubyScriptTask.RetryDelaySeconds)
	d.Set("allow_custom_config", rubyScriptTask.AllowCustomConfig)
	return diags
}

func resourceRubyScriptTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	name := d.Get("name").(string)

	sourceOptions := make(map[string]interface{})
	if d.Get("script_content") != "" {
		sourceOptions["content"] = d.Get("script_content")
	}
	if d.Get("script_path") != "" {
		sourceOptions["contentPath"] = d.Get("script_path")
	}
	sourceOptions["contentRef"] = d.Get("version_ref")
	sourceOptions["repository"] = map[string]interface{}{
		"id": d.Get("repository_id"),
	}
	sourceOptions["sourceType"] = d.Get("source_type")

	taskType := make(map[string]interface{})
	taskType["code"] = "jythonTask"

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"task": map[string]interface{}{
				"name":              name,
				"code":              d.Get("code").(string),
				"file":              sourceOptions,
				"taskType":          taskType,
				"resultType":        d.Get("result_type"),
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
	rubyScriptTask := result.Task
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(rubyScriptTask.ID))
	return resourceRubyScriptTaskRead(ctx, d, meta)
}

func resourceRubyScriptTaskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
