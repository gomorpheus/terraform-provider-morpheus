package morpheus

import (
	"context"

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
				Computed:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "The organization labels associated with the task (Only supported on Morpheus 5.5.3 or higher)",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
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
				Computed:    true,
			},
			"skip_tags": {
				Type:        schema.TypeString,
				Description: "The tags to skip during execution of the ansible playbook",
				Optional:    true,
				Computed:    true,
			},
			"command_options": {
				Type:        schema.TypeString,
				Description: "Additional commands options to pass during the execution of the ansible playbook",
				Optional:    true,
				Computed:    true,
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
				Description: "Custom configuration data to pass during the execution of the ansible playbook",
				Optional:    true,
				Default:     false,
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
	ansiblePlaybookTask := result.Task
	d.SetId(int64ToString(ansiblePlaybookTask.ID))
	d.Set("name", ansiblePlaybookTask.Name)
	d.Set("code", ansiblePlaybookTask.Code)
	d.Set("labels", ansiblePlaybookTask.Labels)
	d.Set("ansible_repo_id", ansiblePlaybookTask.TaskOptions.AnsibleGitId)
	d.Set("git_ref", ansiblePlaybookTask.TaskOptions.AnsibleGitRef)
	d.Set("playbook", ansiblePlaybookTask.TaskOptions.AnsiblePlaybook)
	d.Set("tags", ansiblePlaybookTask.TaskOptions.AnsibleTags)
	d.Set("skip_tags", ansiblePlaybookTask.TaskOptions.AnsibleSkipTags)
	d.Set("command_options", ansiblePlaybookTask.TaskOptions.AnsibleOptions)
	d.Set("execute_target", ansiblePlaybookTask.ExecuteTarget)
	d.Set("retryable", ansiblePlaybookTask.Retryable)
	d.Set("retry_count", ansiblePlaybookTask.RetryCount)
	d.Set("retry_delay_seconds", ansiblePlaybookTask.RetryDelaySeconds)
	d.Set("allow_custom_config", ansiblePlaybookTask.AllowCustomConfig)
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
				"executeTarget":     d.Get("execute_target"),
				"retryable":         d.Get("retryable").(bool),
				"retryCount":        d.Get("retry_count"),
				"retryDelaySeconds": d.Get("retry_delay_seconds"),
				"allowCustomConfig": d.Get("allow_custom_config").(bool),
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
