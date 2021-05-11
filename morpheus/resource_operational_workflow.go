package morpheus

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceOperationalWorkflow() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus operational workflow resource.",
		CreateContext: resourceOperationalWorkflowCreate,
		ReadContext:   resourceOperationalWorkflowRead,
		UpdateContext: resourceOperationalWorkflowUpdate,
		DeleteContext: resourceOperationalWorkflowDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the operational workflow",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the operational workflow",
				Optional:    true,
			},
			"option_types": {
				Type:        schema.TypeList,
				Description: "The option types associated with the operational workflow",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"platform": {
				Type:         schema.TypeString,
				Description:  "The operating system platforms the operational workflow is supported to run on",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"all", "linux", "macos", "windows", ""}, false),
			},
			"allow_custom_config": {
				Type:        schema.TypeBool,
				Description: "",
				Optional:    true,
				Default:     false,
			},
			"visibility": {
				Type:         schema.TypeString,
				Description:  "",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"private", "public", ""}, false),
				Default:      "private",
			},
			"task": {
				Type:        schema.TypeList,
				Description: "",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_id": {
							Type:        schema.TypeInt,
							Description: "",
							Required:    true,
						},
						"task_phase": {
							Type:        schema.TypeString,
							Description: "",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func resourceOperationalWorkflowCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"taskSet": map[string]interface{}{
				"name":              name,
				"description":       description,
				"type":              "operation",
				"optionTypes":       d.Get("option_types"),
				"visibility":        d.Get("visibility"),
				"platform":          d.Get("platform"),
				"allowCustomConfig": d.Get("allow_custom_config"),
				"tasks":             tasks,
			},
		},
	}
	slcB, _ := json.Marshal(req.Body)
	log.Printf("API JSON REQUEST: %s", string(slcB))
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

	resourceOperationalWorkflowRead(ctx, d, meta)
	return diags
}

func resourceOperationalWorkflowRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	result := resp.Result.(*morpheus.GetTaskSetResult)
	workflow := result.TaskSet
	if workflow != nil {
		d.SetId(int64ToString(workflow.ID))
		d.Set("name", workflow.Name)
		d.Set("description", workflow.Description)
		// option types
		var optionTypes []int64
		if workflow.OptionTypes != nil {
			// iterate over the array of tasks
			for i := 0; i < len(workflow.OptionTypes); i++ {
				option := workflow.OptionTypes[i].(map[string]interface{})
				optionID := int64(option["id"].(float64))
				optionTypes = append(optionTypes, optionID)
			}
		}
		d.Set("option_types", optionTypes)
		d.Set("visibility", workflow.Visibility)
		d.Set("allow_custom_config", workflow.AllowCustomConfig)
		d.Set("platform", workflow.Platform)
	} else {
		log.Println(workflow)
		return diag.Errorf("read operation: workflow not found in response data") // should not happen
	}

	return diags
}

func resourceOperationalWorkflowUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	name := d.Get("name").(string)
	description := d.Get("description").(string)

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"taskSet": map[string]interface{}{
				"name":              name,
				"description":       description,
				"type":              "operation",
				"optionTypes":       d.Get("option_types"),
				"visibility":        d.Get("visibility"),
				"platform":          d.Get("platform"),
				"allowCustomConfig": d.Get("allow_custom_config"),
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
	account := result.TaskSet
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(account.ID))
	return resourceOperationalWorkflowRead(ctx, d, meta)
}

func resourceOperationalWorkflowDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
