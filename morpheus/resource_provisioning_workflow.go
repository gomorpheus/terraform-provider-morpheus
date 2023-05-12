package morpheus

import (
	"context"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceProvisioningWorkflow() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus provisioning workflow resource.",
		CreateContext: resourceProvisioningWorkflowCreate,
		ReadContext:   resourceProvisioningWorkflowRead,
		UpdateContext: resourceProvisioningWorkflowUpdate,
		DeleteContext: resourceProvisioningWorkflowDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the provisioning workflow",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the provisioning workflow",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the provisioning workflow",
				Optional:    true,
				Computed:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "The organization labels associated with the workflow (Only supported on Morpheus 5.5.3 or higher)",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"platform": {
				Type:         schema.TypeString,
				Description:  "The operating system platforms the provisioning workflow is supported on (all, linux, macos, windows)",
				ValidateFunc: validation.StringInSlice([]string{"all", "linux", "macos", "windows"}, false),
				Optional:     true,
				Computed:     true,
			},
			"visibility": {
				Type:         schema.TypeString,
				Description:  "Whether the provisioning workflow is visible in sub-tenants or not",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"private", "public"}, false),
				Default:      "private",
			},
			"task": {
				Type:        schema.TypeList,
				Description: "A list of tasks associated with the provisioning workflow",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_id": {
							Type:        schema.TypeInt,
							Description: "The ID of the task to associate with the provisioning workflow",
							Required:    true,
						},
						"task_phase": {
							Type:         schema.TypeString,
							Description:  "The phase that the task is executed (configure, price, preProvision, provision, postProvision, start, stop, preDeploy, deploy, reconfigure, teardown, shutdown, startup)",
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"configure", "price", "preProvision", "provision", "postProvision", "start", "stop", "preDeploy", "deploy", "reconfigure", "teardown", "shutdown", "startup"}, false),
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceProvisioningWorkflowCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	// tasks
	var tasks []map[string]interface{}
	if d.Get("task") != nil {
		taskList := d.Get("task").([]interface{})
		// iterate over the array of tasks
		for i := 0; i < len(taskList); i++ {
			row := make(map[string]interface{})
			taskconfig := taskList[i].(map[string]interface{})
			row["taskId"] = taskconfig["task_id"]
			row["taskPhase"] = taskconfig["task_phase"]
			tasks = append(tasks, row)
		}
	}

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		for _, s := range attr.(*schema.Set).List() {
			labelsPayload = append(labelsPayload, s.(string))
		}
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"taskSet": map[string]interface{}{
				"name":        name,
				"description": description,
				"labels":      labelsPayload,
				"type":        "provision",
				"visibility":  d.Get("visibility"),
				"platform":    d.Get("platform"),
				"tasks":       tasks,
			},
		},
	}

	resp, err := client.CreateTaskSet(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateTaskSetResult)
	environment := result.TaskSet
	// Successfully created resource, now set id
	d.SetId(int64ToString(environment.ID))

	resourceProvisioningWorkflowRead(ctx, d, meta)
	return diags
}

func resourceProvisioningWorkflowRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindTaskSetByName(name)
	} else if id != "" {
		resp, err = client.GetTaskSet(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("TaskSet cannot be read without name or id")
	}

	if err != nil {
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
	result := resp.Result.(*morpheus.GetTaskSetResult)
	workflow := result.TaskSet

	// Tasks
	var tasks []map[string]interface{}
	if len(workflow.TaskSetTasks) != 0 {
		for _, task := range workflow.TaskSetTasks {
			tag := make(map[string]interface{})
			tag["task_phase"] = task.TaskPhase
			tag["task_id"] = task.Task.ID
			tasks = append(tasks, tag)
		}
	}

	if workflow != nil {
		d.SetId(int64ToString(workflow.ID))
		d.Set("name", workflow.Name)
		d.Set("description", workflow.Description)
		d.Set("labels", workflow.Labels)
		d.Set("visibility", workflow.Visibility)
		if workflow.Platform == "" {
			d.Set("platform", "all")
		} else {
			d.Set("platform", workflow.Platform)
		}
		d.Set("task", tasks)
	} else {
		return diag.Errorf("read operation: workflow not found in response data") // should not happen
	}

	return diags
}

func resourceProvisioningWorkflowUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	// tasks
	var tasks []map[string]interface{}
	if d.Get("task") != nil {
		taskList := d.Get("task").([]interface{})
		// iterate over the array of tasks
		for i := 0; i < len(taskList); i++ {
			row := make(map[string]interface{})
			taskconfig := taskList[i].(map[string]interface{})
			row["taskId"] = taskconfig["task_id"]
			row["taskPhase"] = taskconfig["task_phase"]
			tasks = append(tasks, row)
		}
	}

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		for _, s := range attr.(*schema.Set).List() {
			labelsPayload = append(labelsPayload, s.(string))
		}
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"taskSet": map[string]interface{}{
				"name":        name,
				"description": description,
				"labels":      labelsPayload,
				"visibility":  d.Get("visibility"),
				"platform":    d.Get("platform"),
				"tasks":       tasks,
			},
		},
	}
	resp, err := client.UpdateTaskSet(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateTaskSetResult)
	workflow := result.TaskSet
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(workflow.ID))
	return resourceProvisioningWorkflowRead(ctx, d, meta)
}

func resourceProvisioningWorkflowDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteTaskSet(toInt64(id), req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
			return diag.FromErr(err)
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)
			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)
	d.SetId("")
	return diags
}
