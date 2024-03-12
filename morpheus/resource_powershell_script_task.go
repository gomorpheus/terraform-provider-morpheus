package morpheus

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourcePowerShellScriptTask() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus powershell script task resource",
		CreateContext: resourcePowerShellScriptTaskCreate,
		ReadContext:   resourcePowerShellScriptTaskRead,
		UpdateContext: resourcePowerShellScriptTaskUpdate,
		DeleteContext: resourcePowerShellScriptTaskDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the powershell script task",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the powershell script task",
				Required:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the powershell script task",
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
			"elevated_shell": {
				Type:        schema.TypeBool,
				Description: "Run the powershell script with elevated permissions",
				Optional:    true,
				Default:     false,
			},
			"source_type": {
				Type:         schema.TypeString,
				Description:  "The source of the powershell script (local, url or repository)",
				ValidateFunc: validation.StringInSlice([]string{"local", "url", "repository"}, false),
				Required:     true,
			},
			"script_content": {
				Type:        schema.TypeString,
				Description: "The content of the powershell script. Used when the local source type is specified",
				Optional:    true,
				Computed:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					oldPayload := strings.TrimSuffix(old, "\n")
					newPayload := strings.TrimSuffix(new, "\n")
					return oldPayload == newPayload
				},
				StateFunc: func(val interface{}) string {
					return strings.TrimSuffix(val.(string), "\n")
				},
			},
			"script_path": {
				Type:        schema.TypeString,
				Description: "The path of the powershell script, either the url or the path in the repository",
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
				Description:  "The execute target for the powershell script (local, remote or resource)",
				ValidateFunc: validation.StringInSlice([]string{"local", "remote", "resource"}, false),
				Default:      "local",
				Optional:     true,
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
				Description: "Custom configuration data to pass during the execution of the shell script",
				Optional:    true,
				Default:     false,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourcePowerShellScriptTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	taskType["code"] = "winrmTask"

	taskOptions := make(map[string]interface{})
	if d.Get("elevated_shell").(bool) {
		taskOptions["winrm.elevated"] = "on"
	} else {
		taskOptions["winrm.elevated"] = nil
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

	resourcePowerShellScriptTaskRead(ctx, d, meta)
	return diags
}

func resourcePowerShellScriptTaskRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	result := resp.Result.(*morpheus.GetTaskResult)
	powerShellScriptTask := result.Task
	d.SetId(int64ToString(powerShellScriptTask.ID))
	d.Set("name", powerShellScriptTask.Name)
	d.Set("code", powerShellScriptTask.Code)
	d.Set("labels", powerShellScriptTask.Labels)
	d.Set("result_type", powerShellScriptTask.ResultType)
	d.Set("source_type", powerShellScriptTask.File.SourceType)
	d.Set("script_content", powerShellScriptTask.File.Content)
	d.Set("script_path", powerShellScriptTask.File.ContentPath)
	d.Set("version_ref", powerShellScriptTask.File.ContentRef)
	d.Set("execute_target", powerShellScriptTask.ExecuteTarget)
	d.Set("repository_id", powerShellScriptTask.File.Repository.ID)
	if powerShellScriptTask.TaskOptions.WinrmElevated == "on" {
		d.Set("elevated_shell", true)
	} else {
		d.Set("elevated_shell", false)
	}
	d.Set("remote_target_host", powerShellScriptTask.TaskOptions.Host)
	d.Set("remote_target_port", powerShellScriptTask.TaskOptions.Port)
	d.Set("remote_target_username", powerShellScriptTask.TaskOptions.Username)
	d.Set("remote_target_password", powerShellScriptTask.TaskOptions.PasswordHash)
	d.Set("retryable", powerShellScriptTask.Retryable)
	d.Set("retry_count", powerShellScriptTask.RetryCount)
	d.Set("retry_delay_seconds", powerShellScriptTask.RetryDelaySeconds)
	d.Set("allow_custom_config", powerShellScriptTask.AllowCustomConfig)
	return diags
}

func resourcePowerShellScriptTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	taskType["code"] = "winrmTask"

	taskOptions := make(map[string]interface{})
	if d.Get("elevated_shell").(bool) {
		taskOptions["winrm.elevated"] = "on"
	} else {
		taskOptions["winrm.elevated"] = nil
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
	task := result.Task
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(task.ID))
	return resourcePowerShellScriptTaskRead(ctx, d, meta)
}

func resourcePowerShellScriptTaskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

type PowerShellScript struct {
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
			Port              string `json:"port"`
			Host              string `json:"host"`
			Password          string `json:"password"`
			PasswordHash      string `json:"passwordHash"`
			Username          string `json:"username"`
			WinrmElevated     string `json:"winrm.elevated"`
			LocalScriptGitRef string `json:"localScriptGitRef"`
		}
		File struct {
			ID          int    `json:"id"`
			Sourcetype  string `json:"sourceType"`
			Contentref  string `json:"contentRef"`
			Contentpath string `json:"contentPath"`
			Repository  struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"repository"`
			Content interface{} `json:"content"`
		} `json:"file"`
		Resulttype        string    `json:"resultType"`
		Executetarget     string    `json:"executeTarget"`
		Retryable         bool      `json:"retryable"`
		Retrycount        int       `json:"retryCount"`
		Retrydelayseconds int       `json:"retryDelaySeconds"`
		Allowcustomconfig bool      `json:"allowCustomConfig"`
		Datecreated       time.Time `json:"dateCreated"`
		Lastupdated       time.Time `json:"lastUpdated"`
	} `json:"task"`
}
