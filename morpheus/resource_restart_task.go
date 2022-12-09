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

func resourceRestartTask() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus restart task resource",
		CreateContext: resourceRestartTaskCreate,
		ReadContext:   resourceRestartTaskRead,
		UpdateContext: resourceRestartTaskUpdate,
		DeleteContext: resourceRestartTaskDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the restart task",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the restart task",
				Required:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the restart task",
				Optional:    true,
				Computed:    true,
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
				Description: "Custom configuration data to pass during the execution of the restart task",
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceRestartTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	taskType := make(map[string]interface{})
	taskType["code"] = "restartTask"

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"task": map[string]interface{}{
				"name":              name,
				"code":              d.Get("code").(string),
				"taskType":          taskType,
				"executeTarget":     "resource",
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

func resourceRestartTaskRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	var restartTask Restart
	json.Unmarshal(resp.Body, &restartTask)
	d.SetId(intToString(restartTask.Task.ID))
	d.Set("name", restartTask.Task.Name)
	d.Set("code", restartTask.Task.Code)
	d.Set("retryable", restartTask.Task.Retryable)
	d.Set("retry_count", restartTask.Task.Retrycount)
	d.Set("retry_delay_seconds", restartTask.Task.Retrydelayseconds)
	d.Set("allow_custom_config", restartTask.Task.Allowcustomconfig)
	return diags
}

func resourceRestartTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	name := d.Get("name").(string)
	taskType := make(map[string]interface{})
	taskType["code"] = "restartTask"

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"task": map[string]interface{}{
				"name":              name,
				"code":              d.Get("code").(string),
				"taskType":          taskType,
				"executeTarget":     "resource",
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
	restartTask := result.Task
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(restartTask.ID))
	return resourceRestartTaskRead(ctx, d, meta)
}

func resourceRestartTaskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

type Restart struct {
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
		} `json:"taskOptions"`
		File              interface{} `json:"file"`
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
