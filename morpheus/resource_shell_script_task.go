package morpheus

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceShellScriptTask() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus shell script task resource",
		CreateContext: resourceShellScriptTaskCreate,
		ReadContext:   resourceShellScriptTaskRead,
		UpdateContext: resourceShellScriptTaskUpdate,
		DeleteContext: resourceShellScriptTaskDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the shell script task",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the shell script task",
				Required:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the shell script task",
				Optional:    true,
				Computed:    true,
			},
			"result_type": {
				Type:         schema.TypeString,
				Description:  "The expected result type (value, keyValue, json)",
				ValidateFunc: validation.StringInSlice([]string{"value", "keyValue", "json"}, false),
				Optional:     true,
				Computed:     true,
			},
			"sudo": {
				Type:        schema.TypeBool,
				Description: "Whether to run the script as sudo",
				Optional:    true,
				Computed:    true,
			},
			"source_type": {
				Type:         schema.TypeString,
				Description:  "The source of the shell script (local, url or repository)",
				ValidateFunc: validation.StringInSlice([]string{"local", "url", "repository"}, false),
				Required:     true,
			},
			"script_content": {
				Type:        schema.TypeString,
				Description: "The content of the shell script. Used when the local source type is specified",
				Optional:    true,
				Computed:    true,
			},
			"script_path": {
				Type:        schema.TypeString,
				Description: "The path of the shell script, either the url or the path in the repository",
				Optional:    true,
				Computed:    true,
			},
			"repository_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the git repository integration",
				Optional:    true,
				Computed:    true,
			},
			"version_ref": {
				Type:        schema.TypeString,
				Description: "The git reference of the repository to pull (main, master, etc.)",
				Optional:    true,
				Computed:    true,
			},
			"execute_target": {
				Type:         schema.TypeString,
				Description:  "The source of the shell script (local, url or repository)",
				ValidateFunc: validation.StringInSlice([]string{"local", "remote", "resource"}, false),
				Optional:     true,
				Computed:     true,
			},
			"remote_target_host": {
				Type:        schema.TypeString,
				Description: "The hostname or ip address of the remote target",
				Optional:    true,
				Computed:    true,
			},
			"remote_target_port": {
				Type:        schema.TypeString,
				Description: "The port used to connect to the remote target",
				Optional:    true,
				Computed:    true,
			},
			"remote_target_username": {
				Type:        schema.TypeString,
				Description: "The username of the user account used to authenticate to the remote target",
				Optional:    true,
				Computed:    true,
			},
			"remote_target_password": {
				Type:        schema.TypeString,
				Description: "The password of the user account used to authenticate to the remote target",
				Optional:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256_hash := hex.EncodeToString(h.Sum(nil))
					return strings.EqualFold(old, sha256_hash)
					//return strings.ToLower(old) == strings.ToLower(sha256_hash)
				},
				Computed: true,
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

func resourceShellScriptTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	taskType["code"] = "script"

	taskOptions := make(map[string]interface{})
	if d.Get("sudo").(bool) {
		taskOptions["shell.sudo"] = "on"
	} else {
		taskOptions["shell.sudo"] = nil
	}
	if d.Get("remote_target_host") != "" {
		taskOptions["host"] = d.Get("remote_target_host")
	}
	if d.Get("remote_target_port") != "" {
		taskOptions["port"] = d.Get("remote_target_port")
	}
	if d.Get("remote_target_username") != "" {
		taskOptions["username"] = d.Get("remote_target_username")
	}
	if d.Get("remote_target_password") != "" {
		taskOptions["password"] = d.Get("remote_target_password")
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"task": map[string]interface{}{
				"name":              name,
				"code":              d.Get("code").(string),
				"file":              sourceOptions,
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
	log.Printf("Task ID: %s", int64ToString(task.ID))

	resourceShellScriptTaskRead(ctx, d, meta)
	return diags
}

func resourceShellScriptTaskRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	shellScriptTask := result.Task

	d.SetId(int64ToString(shellScriptTask.ID))
	d.Set("name", shellScriptTask.Name)
	d.Set("code", shellScriptTask.Code)
	d.Set("result_type", shellScriptTask.ResultType)
	d.Set("source_type", shellScriptTask.File.SourceType)
	d.Set("script_content", shellScriptTask.File.Content)
	d.Set("script_path", shellScriptTask.File.ContentPath)
	d.Set("version_ref", shellScriptTask.File.ContentRef)
	d.Set("execute_target", shellScriptTask.ExecuteTarget)
	d.Set("repository_id", shellScriptTask.File.Repository.ID)
	if shellScriptTask.TaskOptions.ShellSudo == "on" {
		d.Set("sudo", true)
	} else {
		d.Set("sudo", false)
	}
	d.Set("remote_target_host", shellScriptTask.TaskOptions.Host)
	d.Set("remote_target_port", shellScriptTask.TaskOptions.Port)
	d.Set("remote_target_username", shellScriptTask.TaskOptions.Username)
	d.Set("remote_target_password", shellScriptTask.TaskOptions.PasswordHash)
	d.Set("retryable", shellScriptTask.Retryable)
	d.Set("retry_count", shellScriptTask.RetryCount)
	d.Set("retry_delay_seconds", shellScriptTask.RetryDelaySeconds)
	d.Set("allow_custom_config", shellScriptTask.AllowCustomConfig)
	return diags
}

func resourceShellScriptTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	taskType["code"] = "script"

	taskOptions := make(map[string]interface{})
	if d.Get("sudo").(bool) {
		taskOptions["shell.sudo"] = "on"
	} else {
		taskOptions["shell.sudo"] = nil
	}
	if d.HasChange("remote_target_host") {
		taskOptions["host"] = d.Get("remote_target_host")
	}
	if d.HasChange("remote_target_port") {
		taskOptions["port"] = d.Get("remote_target_port")
	}
	if d.HasChange("remote_target_username") {
		taskOptions["username"] = d.Get("remote_target_username")
	}
	if d.HasChange("remote_target_password") {
		taskOptions["password"] = d.Get("remote_target_password")
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"task": map[string]interface{}{
				"name":              name,
				"code":              d.Get("code").(string),
				"file":              sourceOptions,
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
	return resourceShellScriptTaskRead(ctx, d, meta)
}

func resourceShellScriptTaskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
