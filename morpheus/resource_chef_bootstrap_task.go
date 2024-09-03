package morpheus

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceChefBootstrapTask() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus chef bootstrap task resource",
		CreateContext: resourceChefBootstrapTaskCreate,
		ReadContext:   resourceChefBootstrapTaskRead,
		UpdateContext: resourceChefBootstrapTaskUpdate,
		DeleteContext: resourceChefBootstrapTaskDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the chef bootstrap task",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the chef bootstrap task",
				Required:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the chef bootstrap task",
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
			"chef_server_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the Chef Server integration",
				Optional:    true,
			},
			"environment": {
				Type:        schema.TypeString,
				Description: "The chef environment",
				Optional:    true,
			},
			"run_list": {
				Type:        schema.TypeString,
				Description: "The chef run list",
				Optional:    true,
			},
			"data_bag_key": {
				Type:        schema.TypeString,
				Description: "The chef databag key",
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256_hash := hex.EncodeToString(h.Sum(nil))
					return strings.EqualFold(old, sha256_hash)
				},
				Optional: true,
			},
			"data_bag_key_path": {
				Type:        schema.TypeString,
				Description: "The chef databag key path",
				Optional:    true,
			},
			"node_name": {
				Type:        schema.TypeString,
				Description: "The chef node name",
				Optional:    true,
			},
			"node_attributes": {
				Type:             schema.TypeString,
				Description:      "The chef node attributes (JSON)",
				Optional:         true,
				DiffSuppressFunc: suppressEquivalentJsonDiffs,
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
				Description: "Custom configuration data to pass during the execution of the chef bootstrap",
				Optional:    true,
				Default:     false,
			},
			"visibility": {
				Type:         schema.TypeString,
				Description:  "Whether the task is visible in sub-tenants or not",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"private", "public"}, false),
				Default:      "private",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceChefBootstrapTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)

	taskOptions := make(map[string]interface{})
	taskOptions["chefServerId"] = d.Get("chef_server_id").(int)
	taskOptions["chefEnv"] = d.Get("environment")
	taskOptions["chefRunList"] = d.Get("run_list")
	taskOptions["chefDataKey"] = d.Get("data_bag_key")
	taskOptions["chefDataKeyPath"] = d.Get("data_bag_key_path")
	taskOptions["chefNodeName"] = d.Get("node_name")
	taskOptions["chefAttributes"] = d.Get("node_attributes")

	taskType := make(map[string]interface{})
	taskType["code"] = "chefTask"

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
				"executeTarget":     "resource",
				"retryable":         d.Get("retryable"),
				"retryCount":        d.Get("retry_count"),
				"retryDelaySeconds": d.Get("retry_delay_seconds"),
				"allowCustomConfig": d.Get("allow_custom_config"),
				"visibility":        d.Get("visibility"),
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

	resourceChefBootstrapTaskRead(ctx, d, meta)
	return diags
}

func resourceChefBootstrapTaskRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	chefBootstrapTask := result.Task
	d.SetId(int64ToString(chefBootstrapTask.ID))
	d.Set("name", chefBootstrapTask.Name)
	d.Set("code", chefBootstrapTask.Code)
	d.Set("labels", chefBootstrapTask.Labels)
	serverId, _ := strconv.Atoi(chefBootstrapTask.TaskOptions.ChefServerId)
	d.Set("chef_server_id", serverId)
	d.Set("environment", chefBootstrapTask.TaskOptions.ChefEnv)
	d.Set("run_list", chefBootstrapTask.TaskOptions.ChefRunList)
	d.Set("data_bag_key", chefBootstrapTask.TaskOptions.ChefDataKeyHash)
	d.Set("data_bag_key_path", chefBootstrapTask.TaskOptions.ChefDataKeyPath)
	d.Set("node_name", chefBootstrapTask.TaskOptions.ChefNodeName)
	d.Set("node_attributes", chefBootstrapTask.TaskOptions.ChefAttributes)
	d.Set("retryable", chefBootstrapTask.Retryable)
	d.Set("retry_count", chefBootstrapTask.RetryCount)
	d.Set("retry_delay_seconds", chefBootstrapTask.RetryDelaySeconds)
	d.Set("allow_custom_config", chefBootstrapTask.AllowCustomConfig)
	d.Set("visibility", chefBootstrapTask.Visibility)
	return diags
}

func resourceChefBootstrapTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	name := d.Get("name").(string)

	taskOptions := make(map[string]interface{})
	taskOptions["chefServerId"] = d.Get("chef_server_id")
	taskOptions["chefEnv"] = d.Get("environment")
	taskOptions["chefRunList"] = d.Get("run_list")
	if d.HasChange("data_bag_key") {
		taskOptions["chefDataKey"] = d.Get("data_bag_key")
	}
	taskOptions["chefDataKeyPath"] = d.Get("data_bag_key_path")
	taskOptions["chefNodeName"] = d.Get("node_name")
	taskOptions["chefAttributes"] = d.Get("node_attributes")

	taskType := make(map[string]interface{})
	taskType["code"] = "chefTask"

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
				"executeTarget":     "resource",
				"retryable":         d.Get("retryable").(bool),
				"retryCount":        d.Get("retry_count"),
				"retryDelaySeconds": d.Get("retry_delay_seconds"),
				"allowCustomConfig": d.Get("allow_custom_config").(bool),
				"visibility":        d.Get("visibility"),
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
	chefBootstrapTask := result.Task
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(chefBootstrapTask.ID))
	return resourceChefBootstrapTaskRead(ctx, d, meta)
}

func resourceChefBootstrapTaskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
