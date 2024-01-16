package morpheus

import (
	"context"
	"strings"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEmailTask() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus email task resource",
		CreateContext: resourceEmailTaskCreate,
		ReadContext:   resourceEmailTaskRead,
		UpdateContext: resourceEmailTaskUpdate,
		DeleteContext: resourceEmailTaskDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the email task",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the email task",
				Required:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the email task",
				Optional:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "The organization labels associated with the task (Only supported on Morpheus 5.5.3 or higher)",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"email_address": {
				Type:        schema.TypeString,
				Description: "Email addresses can be entered literally or Morpheus automation variables can be injected, such as <%=instance.createdByEmail%>",
				Required:    true,
			},
			"subject": {
				Type:        schema.TypeString,
				Description: "The subject line of the email, Morpheus automation variables can be injected into the subject field",
				Required:    true,
			},
			"source": {
				Type:        schema.TypeString,
				Description: "Choose local to draft or paste the email directly into the Task. Choose Repository or URL to bring in a template from a Git repository or another outside source (local, repository, url)",
				Optional:    true,
				Default:     "local",
			},
			"content_url": {
				Type:        schema.TypeString,
				Description: "The URL of the template used for the email task, used with a source type of url",
				Optional:    true,
				Computed:    true,
			},
			"content_path": {
				Type:        schema.TypeString,
				Description: "The file path of the template used for the email task, used with a source type of repository",
				Optional:    true,
				Computed:    true,
			},
			"repository_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the git repository to fetch the email template",
				Optional:    true,
				Computed:    true,
			},
			"version_ref": {
				Type:        schema.TypeString,
				Description: "The git reference of the repository to pull (main, master, etc.)",
				Optional:    true,
				Computed:    true,
			},
			"content": {
				Type:        schema.TypeString,
				Description: "The body of the email is HTML. Morpheus automation variables can be injected into the email body when needed. Used with a source type of local",
				Optional:    true,
				Computed:    true,
				StateFunc: func(val interface{}) string {
					return strings.TrimSuffix(val.(string), "\n")
				},
			},
			"skip_wrapped_email_template": {
				Type:        schema.TypeBool,
				Description: "Whether to ignore the Morpheus-styled email template",
				Optional:    true,
				Default:     false,
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
				Description: "Custom configuration data to pass during the execution of the email task",
				Optional:    true,
				Default:     false,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceEmailTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	contentConfig := make(map[string]interface{})

	switch d.Get("source") {
	case "local":
		contentConfig["sourceType"] = "local"
		contentConfig["content"] = d.Get("content").(string)
	case "url":
		contentConfig["sourceType"] = "url"
		contentConfig["contentPath"] = d.Get("content_url").(string)
	case "repository":
		contentConfig["sourceType"] = "repository"
		repository := make(map[string]interface{})
		repository["id"] = d.Get("repository_id")
		contentConfig["contentPath"] = d.Get("content_path")
		if d.Get("version_ref") != "" {
			contentConfig["contentRef"] = d.Get("version_ref")
		}
		contentConfig["repository"] = repository
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
				"name":   name,
				"code":   d.Get("code").(string),
				"labels": labelsPayload,
				"taskType": map[string]interface{}{
					"code": "email",
				},
				"taskOptions": map[string]interface{}{
					"emailAddress":      d.Get("email_address"),
					"emailSubject":      d.Get("subject"),
					"emailSkipTemplate": d.Get("skip_wrapped_email_template"),
				},
				"file":              contentConfig,
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

	resourceEmailTaskRead(ctx, d, meta)
	return diags
}

func resourceEmailTaskRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	emailTask := result.Task

	d.SetId(int64ToString(emailTask.ID))
	d.Set("name", emailTask.Name)
	d.Set("code", emailTask.Code)
	d.Set("labels", emailTask.Labels)
	d.Set("email_address", emailTask.TaskOptions.EmailAddress)
	d.Set("subject", emailTask.TaskOptions.EmailSubject)
	d.Set("source", emailTask.File.SourceType)
	if emailTask.File.SourceType == "url" {
		d.Set("content_url", emailTask.File.ContentPath)
	}
	if emailTask.File.SourceType == "repository" {
		d.Set("content_path", emailTask.File.ContentPath)
		d.Set("repository_id", emailTask.File.Repository.ID)
		d.Set("version_ref", emailTask.File.ContentRef)
	}
	if emailTask.File.SourceType == "local" {
		d.Set("content", emailTask.File.Content)
	}
	if emailTask.TaskOptions.EmailSkipTemplate == "on" {
		d.Set("skip_wrapped_email_template", true)
	} else {
		d.Set("skip_wrapped_email_template", false)
	}
	d.Set("retryable", emailTask.Retryable)
	d.Set("retry_count", emailTask.RetryCount)
	d.Set("retry_delay_seconds", emailTask.RetryDelaySeconds)
	d.Set("allow_custom_config", emailTask.AllowCustomConfig)
	return diags
}

func resourceEmailTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	name := d.Get("name").(string)
	contentConfig := make(map[string]interface{})

	switch d.Get("source") {
	case "local":
		contentConfig["sourceType"] = "local"
		contentConfig["content"] = d.Get("content").(string)
	case "url":
		contentConfig["sourceType"] = "url"
		contentConfig["contentPath"] = d.Get("content_url").(string)
	case "repository":
		contentConfig["sourceType"] = "repository"
		repository := make(map[string]interface{})
		repository["id"] = d.Get("repository_id")
		contentConfig["contentPath"] = d.Get("content_path")
		if d.HasChange("version_ref") {
			contentConfig["contentRef"] = d.Get("version_ref")
		}
		contentConfig["repository"] = repository
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
				"name":   name,
				"code":   d.Get("code").(string),
				"labels": labelsPayload,
				"taskType": map[string]interface{}{
					"code": "email",
				},
				"taskOptions": map[string]interface{}{
					"emailAddress":      d.Get("email_address"),
					"emailSubject":      d.Get("subject"),
					"emailSkipTemplate": d.Get("skip_wrapped_email_template"),
				},
				"file":              contentConfig,
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
	emailTask := result.Task
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(emailTask.ID))
	return resourceEmailTaskRead(ctx, d, meta)
}

func resourceEmailTaskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
