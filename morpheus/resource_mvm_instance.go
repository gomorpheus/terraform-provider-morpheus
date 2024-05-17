package morpheus

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMVMInstance() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus instance resource.",
		CreateContext: resourceMVMInstanceCreate,
		ReadContext:   resourceMVMInstanceRead,
		UpdateContext: resourceMVMInstanceUpdate,
		DeleteContext: resourceMVMInstanceDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(45 * time.Minute),
			Delete: schema.DefaultTimeout(45 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the instance",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The name of the instance",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"description": {
				Description: "The user friendly description of the instance",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"cloud_id": {
				Description: "The ID of the cloud associated with the instance",
				Type:        schema.TypeInt,
				ForceNew:    true,
				Required:    true,
			},
			"group_id": {
				Description: "The ID of the group associated with the instance",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"instance_type_id": {
				Description: "The type of instance to provision",
				Type:        schema.TypeInt,
				ForceNew:    true,
				Required:    true,
			},
			"instance_layout_id": {
				Description: "The layout to provision the instance from",
				Type:        schema.TypeInt,
				ForceNew:    true,
				Required:    true,
			},
			"plan_id": {
				Description: "The service plan associated with the instance",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"resource_pool_name": {
				Description: "The name of the resource pool (cluster) to provision the instance to",
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
			},
			"domain_id": {
				Description: "The ID of the network domain to provision the instance to",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"environment": {
				Description: "The environment to assign the instance to",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"labels": {
				Type:        schema.TypeList,
				Description: "The list of labels to add to the instance",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"tags": {
				Description: "Tags to assign to the instance",
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"custom_options": {
				Description: "Custom options to pass to the instance",
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"workflow_id": {
				Description:   "The ID of the provisioning workflow to execute (`workflow_name` can be used alternatively, only one is needed)",
				Type:          schema.TypeInt,
				ForceNew:      true,
				Optional:      true,
				ConflictsWith: []string{"workflow_name"},
			},
			"workflow_name": {
				Description:   "The name of the provisioning workflow to execute (`workflow_id` can be used alternatively, only one is needed)",
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				ConflictsWith: []string{"workflow_id"},
			},
			"create_user": {
				Description: "Whether to create a user account on the instance that is associated with the provisioning user account",
				Type:        schema.TypeBool,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Default:     "",
			},
			"user_group_id": {
				Description: "The id of the user group associated with the instance",
				Type:        schema.TypeInt,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
			},
			"image_id": {
				Description: "The id of the image associated with the instance",
				Type:        schema.TypeInt,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
			},
			"asset_tag": {
				Description: "The asset tag associated with the instance",
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
			},
			"qemu_arguments": {
				Description: "The qemu arguments to add to the instance",
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
			},
			"skip_agent_install": {
				Description: "Whether to skip installation of the Morpheus agent",
				Type:        schema.TypeBool,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
			},
			"nested_virtualization": {
				Description: "Whether to skip configuration of nested virtualization",
				Type:        schema.TypeBool,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
			},
			"attach_virtio_drivers": {
				Description: "Whether to attach the virtio drivers to the instance",
				Type:        schema.TypeBool,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
			},
			"evar": {
				Type:        schema.TypeList,
				Description: "The environment variables to create",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "The name of the environment variable",
							Optional:    true,
						},
						"value": {
							Type:        schema.TypeString,
							Description: "The value of the environment variable",
							Optional:    true,
						},
						"export": {
							Type:        schema.TypeBool,
							Description: "Whether the environment variable is exported as an instance tag",
							Optional:    true,
						},
						"masked": {
							Type:        schema.TypeBool,
							Description: "Whether the environment variable is masked for security purposes",
							Optional:    true,
						},
					},
				},
			},
			"volumes": {
				Description: "The instance volumes to create",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"root": {
							Description: "Whether the volume is the root volume of the instance",
							Type:        schema.TypeBool,
							Optional:    true,
						},
						"name": {
							Description: "The name/type of the LV being created",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"size": {
							Description: "The size of the LV being created",
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
						},
						"size_id": {
							Description: "The ID of an existing LV to assign to the instance",
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
						},
						"storage_type": {
							Description: "The ID of the LV type",
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
						},
						"datastore_id": {
							Description: "The ID of the datastore",
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"interfaces": {
				Description: "The instance network interfaces to create",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_id": {
							Description: "The network to assign the network interface to",
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
						},
						"network_group": {
							Description: "Whether the network id provided is for a network group or not",
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
						},
						"ip_address": {
							Description: "",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"ip_mode": {
							Description: "",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"network_interface_type_id": {
							Description: "The network interface type",
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"ip_address": {
				Description: "The instance IP address",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceMVMInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	group := d.Get("group_id").(int)
	cloud := d.Get("cloud_id").(int)
	name := d.Get("name").(string)

	// Service Plan
	planResp, err := client.GetPlan(int64(d.Get("plan_id").(int)), &morpheus.Request{})
	if err != nil {
		diag.FromErr(err)
	}
	planResult := planResp.Result.(*morpheus.GetPlanResult)
	plan := planResult.Plan

	// Instance Type
	instanceTypeResp, err := client.GetInstanceType(int64(d.Get("instance_type_id").(int)), &morpheus.Request{})
	if err != nil {
		diag.FromErr(err)
	}
	instanceTypeResult := instanceTypeResp.Result.(*morpheus.GetInstanceTypeResult)
	instanceTypeCode := instanceTypeResult.InstanceType.Code

	// Instance Layout
	instanceLayoutResp, err := client.GetInstanceLayout(int64(d.Get("instance_layout_id").(int)), &morpheus.Request{})
	if err != nil {
		diag.FromErr(err)
	}
	instanceLayoutResult := instanceLayoutResp.Result.(*morpheus.GetInstanceLayoutResult)
	instanceLayout := instanceLayoutResult.InstanceLayout

	// Config
	config := make(map[string]interface{})

	// Resource Pool
	resourcePoolResp, err := client.Execute(&morpheus.Request{
		Method:      "GET",
		Path:        fmt.Sprintf("/api/options/zonePools?layoutId=%d", d.Get("instance_layout_id").(int)),
		QueryParams: map[string]string{},
	})
	if err != nil {
		diag.FromErr(err)
	}

	var itemResponsePayload ResourcePoolOptions
	json.Unmarshal(resourcePoolResp.Body, &itemResponsePayload)
	var resourcePoolId int
	for _, v := range itemResponsePayload.Data {
		if v.ProviderType == "mvm" {
			resourcePoolId = v.Id
		}
	}

	config["resourcePoolId"] = resourcePoolId
	config["poolProviderType"] = "mvm"

	// Custom Options
	if d.Get("custom_options") != nil {
		customOptionsInput := d.Get("custom_options").(map[string]interface{})
		customOptions := make(map[string]interface{})
		for key, value := range customOptionsInput {
			customOptions[key] = value.(string)
		}
		config["customOptions"] = customOptions
	}

	// Create User
	config["createUser"] = d.Get("create_user").(bool)

	// Asset Tag
	if d.Get("asset_tag") != nil {
		config["smbiosAssetTag"] = d.Get("asset_tag").(string)
	}

	// Image ID
	if d.Get("image_id") != nil {
		config["imageId"] = d.Get("image_id").(int)
	}

	// Skip Agent Install
	config["noAgent"] = d.Get("skip_agent_install").(bool)

	// Nested Virtualization
	config["nestedVirtualization"] = d.Get("nested_virtualization").(bool)

	instancePayload := map[string]interface{}{
		"name": name,
		"type": instanceTypeCode,
		"site": map[string]interface{}{
			"id": group,
		},
		"plan": map[string]interface{}{
			"id":   plan.ID,
			"code": plan.Code,
			"name": plan.Name,
		},
		"layout": map[string]interface{}{
			"id":   instanceLayout.ID,
			"code": instanceLayout.Code,
			"name": instanceLayout.Name,
		},
	}

	// Description
	if d.Get("description") != nil {
		instancePayload["description"] = d.Get("description").(string)
	}

	// Environment
	if d.Get("environment") != nil {
		instancePayload["instanceContext"] = d.Get("environment").(string)
	}

	// User Group ID
	if d.Get("user_group_id") != 0 {
		userGroupPayload := map[string]interface{}{
			"id": d.Get("user_group_id").(int),
		}
		instancePayload["userGroup"] = userGroupPayload
	}

	// Network Domain
	if d.Get("domain_id") != 0 {
		domainConfig := make(map[string]interface{})
		domainConfig["id"] = d.Get("domain_id").(int)
		instancePayload["networkDomain"] = domainConfig
	}

	payload := map[string]interface{}{
		"zoneId":   cloud,
		"instance": instancePayload,
		"config":   config,
	}

	// tags
	if d.Get("tags") != nil {
		tagsInput := d.Get("tags").(map[string]interface{})
		var tags []map[string]interface{}
		for key, value := range tagsInput {
			tag := make(map[string]interface{})
			tag["name"] = key
			tag["value"] = value.(string)
			tags = append(tags, tag)
		}
		payload["tags"] = tags
	}

	// Labels
	if d.Get("labels") != nil {
		payload["labels"] = d.Get("labels")
	}

	// Provisioning Workflow ID
	if d.Get("workflow_id") != nil {
		payload["taskSetId"] = d.Get("workflow_id")
	}

	// Provisioning Workflow Name
	if d.Get("workflow_name") != nil {
		payload["taskSetName"] = d.Get("workflow_name")
	}

	// Environment Variables
	if d.Get("evar") != "" {
		payload["evars"] = parseEnvironmentVariables(d.Get("evar").([]interface{}))
	}

	// Network Interfaces
	if d.Get("interfaces") != nil {
		payload["networkInterfaces"] = parseNetworkInterfaces(d.Get("interfaces").([]interface{}))
	}

	// Volumes
	if d.Get("volumes") != nil {
		payload["volumes"] = parseStorageVolumes(d.Get("volumes").([]interface{}))
	}

	log.Printf("API Request: %s", payload)

	req := &morpheus.Request{Body: payload}
	resp, err := client.CreateInstance(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.CreateInstanceResult)
	instance := result.Instance

	stateConf := &resource.StateChangeConf{
		Pending: []string{"provisioning", "starting", "stopping", "pending"},
		Target:  []string{"running", "failed", "warning", "denied", "cancelled", "suspended"},
		Refresh: func() (interface{}, string, error) {
			instanceDetails, err := client.GetInstance(instance.ID, &morpheus.Request{})
			if err != nil {
				return "", "", err
			}
			result := instanceDetails.Result.(*morpheus.GetInstanceResult)
			instance := result.Instance
			return result, instance.Status, nil
		},
		Timeout:      3 * time.Hour,
		MinTimeout:   1 * time.Minute,
		Delay:        2 * time.Minute,
		PollInterval: 30 * time.Second,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error creating instance: %s", err)
	}

	// Successfully created resource, now set id
	d.SetId(int64ToString(instance.ID))
	resourceMVMInstanceRead(ctx, d, meta)
	return diags
}

func resourceMVMInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindInstanceByName(name)
	} else if id != "" {
		resp, err = client.GetInstance(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Instance cannot be read without name or id")
	}
	if err != nil {
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
	result := resp.Result.(*morpheus.GetInstanceResult)
	instance := result.Instance
	if instance == nil {
		return diag.Errorf("Instance not found in response data.") // should not happen
	}

	d.SetId(int64ToString(instance.ID))
	d.Set("name", instance.Name)
	d.Set("description", instance.Description)
	d.Set("cloud_id", instance.Cloud.ID)
	d.Set("group_id", instance.Group.ID)
	d.Set("instance_type_id", instance.InstanceType.ID)
	d.Set("instance_layout_id", instance.Layout.ID)
	d.Set("plan_id", instance.Plan.ID)

	// Resource Pool
	resourcePoolResp, err := client.Execute(&morpheus.Request{
		Method:      "GET",
		Path:        fmt.Sprintf("/api/options/zonePools?layoutId=%d", instance.Layout.ID),
		QueryParams: map[string]string{},
	})
	if err != nil {
		diag.FromErr(err)
	}

	var itemResponsePayload ResourcePoolOptions
	json.Unmarshal(resourcePoolResp.Body, &itemResponsePayload)
	var resourcePoolName string
	for _, v := range itemResponsePayload.Data {
		if v.ProviderType == "mvm" {
			resourcePoolName = v.Name
		}
	}
	d.Set("resource_pool_name", resourcePoolName)
	d.Set("environment", instance.Environment)
	d.Set("labels", instance.Labels)
	d.Set("evar", instance.EnvironmentVariables)
	// Tags
	tags := make(map[string]interface{})
	if instance.Tags != nil {
		output := instance.Tags
		tagList := output
		// iterate over the array of tags
		for i := 0; i < len(tagList); i++ {
			tag := tagList[i]
			tagName := tag.Name
			tags[tagName] = tag.Value
		}
	}
	d.Set("tags", tags)
	if instance.Config["userGroup"] != nil {
		userGroup := instance.Config["userGroup"].(map[string]interface{})
		d.Set("user_group_id", userGroup["id"])
	}
	d.Set("create_user", instance.Config["createUser"])
	d.Set("asset_tag", instance.Config["smbiosAssetTag"])
	d.Set("skip_agent_install", instance.Config["noAgent"])
	if instance.Config["nestedVirtualization"] == "off" {
		d.Set("nested_virtualization", false)
	} else {
		d.Set("nested_virtualization", true)
	}
	d.Set("custom_options", instance.Config["customOptions"])
	d.Set("domain_id", instance.NetworkDomain.Id)
	d.Set("ip_address", instance.ConnectionInfo[0].Ip)
	return diags
}

func resourceMVMInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	name := d.Get("name").(string)
	description := d.Get("description").(string)

	// Tags
	var tags []map[string]interface{}
	if d.HasChange("tags") {
		tagsInput := d.Get("tags").(map[string]interface{})
		for key, value := range tagsInput {
			tag := make(map[string]interface{})
			tag["name"] = key
			tag["value"] = value.(string)
			tags = append(tags, tag)
		}
	}
	config := make(map[string]interface{})

	// Custom Options
	if d.Get("custom_options") != nil {
		customOptionsInput := d.Get("custom_options").(map[string]interface{})
		customOptions := make(map[string]interface{})
		for key, value := range customOptionsInput {
			customOptions[key] = value.(string)
		}
		config["customOptions"] = customOptions
	}

	instancePayload := map[string]interface{}{
		"name":            name,
		"description":     description,
		"labels":          d.Get("labels"),
		"tags":            tags,
		"instanceContext": d.Get("environment"),
		"config":          config,
	}
	payload := map[string]interface{}{
		"instance": instancePayload,
	}
	req := &morpheus.Request{Body: payload}
	resp, err := client.UpdateInstance(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateInstanceResult)
	instance := result.Instance
	// Successfully updated resource, now set id
	d.SetId(int64ToString(instance.ID))
	return resourceMVMInstanceRead(ctx, d, meta)
}

func resourceMVMInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{
		QueryParams: map[string]string{},
	}
	if USE_FORCE {
		req.QueryParams["force"] = "true"
	}
	resp, err := client.DeleteInstance(toInt64(id), req)
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

type ResourcePoolOptions struct {
	Success bool `json:"success"`
	Data    []struct {
		Id           int    `json:"id"`
		Name         string `json:"name"`
		IsGroup      bool   `json:"isGroup"`
		Group        string `json:"group"`
		IsDefault    bool   `json:"isDefault"`
		Type         string `json:"type"`
		ProviderType string `json:"providerType"`
		Value        string `json:"value"`
	} `json:"data"`
}
