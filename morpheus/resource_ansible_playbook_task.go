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

func resourceAnsiblePlaybookTask() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus ansible playbook task resource",
		CreateContext: resourceAnsiblePlaybookTaskCreate,
		ReadContext:   resourceAnsiblePlaybookTaskRead,
		UpdateContext: resourceAnsiblePlaybookTaskUpdate,
		DeleteContext: resourceAnsiblePlaybookTaskDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the ansible playbook task",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the ansible playbook task",
				Required:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the ansible playbook task",
				Optional:    true,
			},
			"ansible_repo_id": {
				Type:        schema.TypeString,
				Description: "The id of the ansible repo",
				Optional:    true,
			},
			"git_ref": {
				Type:        schema.TypeString,
				Description: "The git reference of the ansible repo to pull (main, master, etc.)",
				Optional:    true,
			},
			"playbook": {
				Type:        schema.TypeString,
				Description: "The name of the ansible playbook to execute",
				Required:    true,
			},
			"tags": {
				Type:        schema.TypeString,
				Description: "The tags to specify during execution of the ansible playbook",
				Optional:    true,
			},
			"skip_tags": {
				Type:        schema.TypeString,
				Description: "The tags to skip during execution of the ansible playbook",
				Optional:    true,
			},
			"command_options": {
				Type:        schema.TypeString,
				Description: "Additional commands options to pass during the execution of the ansible playbook",
				Optional:    true,
			},
			"execute_target": {
				Type:        schema.TypeString,
				Description: "The target that the ansible playbook will be executed on",
				Optional:    true,
			},
			"retryable": {
				Type:        schema.TypeBool,
				Description: "Whether to retry the task if there is a failure",
				Optional:    true,
			},
			"retry_count": {
				Type:        schema.TypeInt,
				Description: "The number of times to retry the task if there is a failure",
				Optional:    true,
				Default:     false,
			},
			"retry_delay_seconds": {
				Type:        schema.TypeInt,
				Description: "The number of seconds to wait between retry attempts",
				Optional:    true,
			},
			"allow_custom_config": {
				Type:        schema.TypeBool,
				Description: "Custom configuration data to pass during the execution of the ansible playbook",
				Optional:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceAnsiblePlaybookTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)

	taskOptions := make(map[string]interface{})
	taskOptions["ansibleGitId"] = d.Get("ansible_repo_id")
	taskOptions["ansibleGitRef"] = d.Get("git_ref")
	taskOptions["ansiblePlaybook"] = d.Get("playbook")
	taskOptions["ansibleTags"] = d.Get("tags")
	taskOptions["ansibleSkipTags"] = d.Get("skip_tags")
	taskOptions["ansibleOptions"] = d.Get("command_options")

	taskType := make(map[string]interface{})
	taskType["code"] = "ansibleTask"

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"task": map[string]interface{}{
				"name":              name,
				"code":              d.Get("code").(string),
				"taskType":          taskType,
				"taskOptions":       taskOptions,
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

	resourceAnsiblePlaybookTaskRead(ctx, d, meta)
	return diags
}

func resourceAnsiblePlaybookTaskRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	var ansiblePlaybookTask AnsiblePlaybook
	json.Unmarshal(resp.Body, &ansiblePlaybookTask)
	d.SetId(intToString(ansiblePlaybookTask.Task.ID))
	d.Set("name", ansiblePlaybookTask.Task.Name)
	d.Set("code", ansiblePlaybookTask.Task.Code)
	d.Set("ansible_repo_id", ansiblePlaybookTask.Task.Taskoptions.Ansiblegitid)
	d.Set("git_ref", ansiblePlaybookTask.Task.Taskoptions.Ansiblegitref)
	d.Set("playbook", ansiblePlaybookTask.Task.Taskoptions.Ansibleplaybook)
	d.Set("tags", ansiblePlaybookTask.Task.Taskoptions.Ansibletags)
	d.Set("skip_tags", ansiblePlaybookTask.Task.Taskoptions.Ansibleskiptags)
	d.Set("command_options", ansiblePlaybookTask.Task.Taskoptions.Ansibleoptions)
	d.Set("execute_target", ansiblePlaybookTask.Task.Executetarget)
	d.Set("retryable", ansiblePlaybookTask.Task.Retryable)
	d.Set("retry_count", ansiblePlaybookTask.Task.Retrycount)
	d.Set("retry_delay_seconds", ansiblePlaybookTask.Task.Retrydelayseconds)
	d.Set("allow_custom_config", ansiblePlaybookTask.Task.Allowcustomconfig)
	return diags
}

func resourceAnsiblePlaybookTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	name := d.Get("name").(string)

	taskOptions := make(map[string]interface{})
	taskOptions["ansibleGitId"] = d.Get("ansible_repo_id")
	taskOptions["ansibleGitRef"] = d.Get("git_ref")
	taskOptions["ansiblePlaybook"] = d.Get("playbook")
	taskOptions["ansibleTags"] = d.Get("tags")
	taskOptions["ansibleSkipTags"] = d.Get("skip_tags")
	taskOptions["ansibleOptions"] = d.Get("command_options")

	taskType := make(map[string]interface{})
	taskType["code"] = "ansibleTask"

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"task": map[string]interface{}{
				"name":              name,
				"code":              d.Get("code").(string),
				"taskType":          taskType,
				"taskOptions":       taskOptions,
				"executeTarget":     d.Get("execute_target"),
				"retryable":         d.Get("retryable").(bool),
				"retryCount":        d.Get("retry_count"),
				"retryDelaySeconds": d.Get("retry_delay_seconds"),
				"allowCustomConfig": d.Get("allow_custom_config").(bool),
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
	ansiblePlaybookTask := result.Task
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(ansiblePlaybookTask.ID))
	return resourceAnsiblePlaybookTaskRead(ctx, d, meta)
}

func resourceAnsiblePlaybookTaskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

type AnsiblePlaybook struct {
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
			Ansibleoptions  string `json:"ansibleOptions"`
			Ansibletags     string `json:"ansibleTags"`
			Ansibleplaybook string `json:"ansiblePlaybook"`
			Ansiblegitref   string `json:"ansibleGitRef"`
			Ansibleskiptags string `json:"ansibleSkipTags"`
			Ansiblegitid    string `json:"ansibleGitId"`
		} `json:"taskOptions"`
		File              interface{} `json:"file"`
		Resulttype        interface{} `json:"resultType"`
		Executetarget     string      `json:"executeTarget"`
		Retryable         bool        `json:"retryable"`
		Retrycount        int         `json:"retryCount"`
		Retrydelayseconds int         `json:"retryDelaySeconds"`
		Allowcustomconfig bool        `json:"allowCustomConfig"`
		Datecreated       time.Time   `json:"dateCreated"`
		Lastupdated       time.Time   `json:"lastUpdated"`
	} `json:"task"`
	Success bool `json:"success"`
}
