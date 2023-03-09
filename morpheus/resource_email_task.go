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
				Description: "Choose local to draft or paste the email directly into the Task. Choose Repository or URL to bring in a template from a Git repository or another outside source",
				Optional:    true,
				Default:     "local",
			},
			"content": {
				Type:        schema.TypeString,
				Description: "The body of the email is HTML. Morpheus automation variables can be injected into the email body when needed",
				Optional:    true,
				Default:     "",
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

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"task": map[string]interface{}{
				"name": name,
				"code": d.Get("code").(string),
				"taskType": map[string]interface{}{
					"code": "email",
				},
				"taskOptions": map[string]interface{}{
					"emailAddress":      d.Get("email_address"),
					"emailSubject":      d.Get("subject"),
					"emailSkipTemplate": d.Get("skip_wrapped_email_template"),
				},
				"file": map[string]interface{}{
					"content": d.Get("content"),
				},
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

	resourceRestartTaskRead(ctx, d, meta)
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
			return diag.FromErr(err)
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)
			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)

	// store resource data
	var emailTask Email
	json.Unmarshal(resp.Body, &emailTask)
	d.SetId(intToString(emailTask.Task.ID))
	d.Set("name", emailTask.Task.Name)
	d.Set("code", emailTask.Task.Code)
	d.Set("email_address", emailTask.Task.Taskoptions.Emailaddress)
	d.Set("subject", emailTask.Task.Taskoptions.Emailsubject)
	d.Set("source", emailTask.Task.File.Sourcetype)
	d.Set("content", emailTask.Task.File.Content)
	d.Set("skip_wrapped_email_template", emailTask.Task.Taskoptions.Emailskiptemplate)
	d.Set("retryable", emailTask.Task.Retryable)
	d.Set("retry_count", emailTask.Task.Retrycount)
	d.Set("retry_delay_seconds", emailTask.Task.Retrydelayseconds)
	d.Set("allow_custom_config", emailTask.Task.Allowcustomconfig)
	return diags
}

func resourceEmailTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	name := d.Get("name").(string)
	taskType := make(map[string]interface{})
	taskType["code"] = "email"

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"task": map[string]interface{}{
				"name": name,
				"code": d.Get("code").(string),
				"taskType": map[string]interface{}{
					"code": "email",
				},
				"taskOptions": map[string]interface{}{
					"emailAddress":      d.Get("email_address"),
					"emailSubject":      d.Get("subject"),
					"emailSkipTemplate": d.Get("skip_wrapped_email_template"),
				},
				"file": map[string]interface{}{
					"content": d.Get("content"),
				},
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

type Email struct {
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
			Username          interface{} `json:"username"`
			Host              interface{} `json:"host"`
			Localscriptgitref interface{} `json:"localScriptGitRef"`
			Password          interface{} `json:"password"`
			Passwordhash      interface{} `json:"passwordHash"`
			Port              interface{} `json:"port"`
			Emailaddress      interface{} `json:"emailAddress"`
			Emailsubject      interface{} `json:"emailSubject"`
			Emailskiptemplate interface{} `json:"emailSkipTemplate"`
		} `json:"taskOptions"`
		File struct {
			Id          interface{} `json:"id"`
			Sourcetype  interface{} `json:"sourceType"`
			Contentref  interface{} `json:"contentRef"`
			Contentpath interface{} `json:"contentPath"`
			Repository  interface{} `json:"repository"`
			Content     interface{} `json:"content"`
		} `json:"file"`
		Resulttype        interface{} `json:"resultType"`
		Executetarget     string      `json:"executeTarget"`
		Retryable         bool        `json:"retryable"`
		Retrycount        int         `json:"retryCount"`
		Retrydelayseconds int         `json:"retryDelaySeconds"`
		Allowcustomconfig bool        `json:"allowCustomConfig"`
		Credential        struct {
			Type string `json:"type"`
		} `json:"credential"`
		Datecreated time.Time `json:"dateCreated"`
		Lastupdated time.Time `json:"lastUpdated"`
	} `json:"task"`
}
