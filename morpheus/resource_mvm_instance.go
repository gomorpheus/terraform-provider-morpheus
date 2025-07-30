package morpheus

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMVMInstance() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus MVM instance resource.",
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
			"display_name": {
				Description: "The display name of the instance",
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
				Description: "The ID of the instance type to provision the instance from",
				Type:        schema.TypeInt,
				ForceNew:    true,
				Required:    true,
			},
			"instance_layout_id": {
				Description: "The ID of the instance layout to provision the instance from",
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
				Required:    true,
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
			},
			"user_group_id": {
				Description: "The ID of the user group associated with the instance",
				Type:        schema.TypeInt,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
			},
			"image_id": {
				Description: "The ID of the image associated with the instance (Only neccessary when using the default MVM instance type that requires specifying a virtual image)",
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
				Description: "Whether to enable nested virtualization",
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
			"storage_volume": {
				Description: "The instance volumes to create",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Description: "The storage volume uuid",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"root": {
							Description: "Whether the volume is the root volume of the instance",
							Type:        schema.TypeBool,
							Required:    true,
						},
						"name": {
							Description: "The name of the volume",
							Type:        schema.TypeString,
							Required:    true,
						},
						"size": {
							Description: "The size of the volume in GB",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"storage_type": {
							Description: "The storage volume type ID",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"datastore_id": {
							Description: "The ID of the datastore",
							Type:        schema.TypeInt,
							Required:    true,
						},
					},
				},
			},
			"network_interface": {
				Description: "The instance network interfaces to create",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_id": {
							Description: "The network to assign the network interface to",
							Type:        schema.TypeInt,
							Required:    true,
						},
						/* AWAITING PLATFORM SUPPORT
						"network_group": {
							Description: "Whether the network id provided is for a network group or not",
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
						},
						*/
						"ip_address": {
							Description: "The IP address to assign to the instance",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"ip_mode": {
							Description: "The IP address assignment mode",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"network_interface_type_id": {
							Description: "The id of the network interface type",
							Type:        schema.TypeInt,
							Required:    true,
						},
					},
				},
			},
			"primary_ip_address": {
				Description: "The instance primary IP address",
				Type:        schema.TypeString,
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
	displayName := d.Get("display_name").(string)

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
	if err := json.Unmarshal(resourcePoolResp.Body, &itemResponsePayload); err != nil {
		return diag.FromErr(err)
	}
	var resourcePoolId int
	resourcePoolFound := false
	for _, v := range itemResponsePayload.Data {
		if v.ProviderType == "mvm" && v.Name == d.Get("resource_pool_name").(string) {
			resourcePoolId = v.Id
			resourcePoolFound = true
		}
	}
	if !resourcePoolFound {
		return diag.Errorf("resource pool with name %s not found with providerType mvm", d.Get("resource_pool_name").(string))
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
		"name":        name,
		"displayName": displayName,
		"type":        instanceTypeCode,
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
	if d.Get("network_interface") != nil {
		payload["networkInterfaces"] = generateMVMInstanceCreatePayload(d.Get("network_interface").([]interface{}))
	}

	// Volumes
	if d.Get("storage_volume") != nil {
		payload["volumes"] = parseStorageVolumes(d.Get("storage_volume").([]interface{}))
	}

	req := &morpheus.Request{Body: payload}
	resp, err := client.CreateInstance(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.CreateInstanceResult)
	instance := result.Instance
	instanceStatus := "provisioning"

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
			instanceStatus = instance.Status
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

	// Fail the instance deployment if the
	// instance status is in a failed state
	if instanceStatus == "failed" {
		return diag.Errorf("error creating instance: failed to create server")
	}
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
	d.Set("environment", instance.Environment)
	d.Set("labels", instance.Labels)

	var evars []map[string]interface{}
	evarMap := make(map[string]string, len(instance.EnvironmentVariables))
	for i := 0; i < len(instance.EnvironmentVariables); i++ {
		evar := instance.EnvironmentVariables[i]
		row := make(map[string]interface{})
		row["name"] = evar.Name
		value := fmt.Sprintf("%v", evar.Value)
		row["value"] = value
		row["export"] = evar.Export
		row["masked"] = evar.Masked
		evarMap[evar.Name] = value
		evars = append(evars, row)
	}

	// If the evar field is set, we need to check if the evars match
	if d.Get("evar") != nil {
		for _, row := range d.Get("evar").([]interface{}) {
			evarData := row.(map[string]interface{})
			evarName := evarData["name"].(string)
			value := evarData["value"].(string)

			if mapValue, exists := evarMap[evarName]; exists {
				if mapValue != value {
					return diag.Errorf("evar %s is missing from returned evar map", evarName)
				}
			}
		}
	}

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
	d.Set("primary_ip_address", instance.ConnectionInfo[0].Ip)

	var volumes []map[string]interface{}
	// iterate over the array of volumes
	for i := 0; i < len(instance.Volumes); i++ {
		row := make(map[string]interface{})
		volume := instance.Volumes[i]
		row["uuid"] = volume.Uuid
		row["root"] = volume.RootVolume.(bool)
		row["name"] = volume.Name
		if volume.Size != nil {
			row["size"] = volume.Size.(float64)
		}
		if volume.StorageType != nil {
			row["storage_type"] = volume.StorageType.(float64)
		}
		if volume.DatastoreId != nil {
			datastoreId, errConv := convertToInt(volume.DatastoreId)
			if errConv != nil {
				log.Printf("Error converting datastore ID to int: %s", errConv)

				return diag.FromErr(errConv)
			}

			row["datastore_id"] = datastoreId
		}
		volumes = append(volumes, row)
	}

	// If the storage_volume field is set, we need to check if the volumes match
	// And it there is a match, we need to set the uuid in the state
	var volumesForState []map[string]interface{}
	if d.Get("storage_volume") != nil {
		for _, row := range d.Get("storage_volume").([]interface{}) {
			volumeData := row.(map[string]interface{})
			volumeName := volumeData["name"]
			datastoreId := volumeData["datastore_id"]

			found := false
			for _, v := range volumes {
				if v["name"] == volumeName && v["datastore_id"] == datastoreId {
					found = true
					volumeData["uuid"] = v["uuid"]
					volumesForState = append(volumesForState, volumeData)

					break
				}
			}
			if !found {
				return diag.Errorf("storage_volume %s is missing from returned storage volume map", volumeName)
			}
		}
	}
	if len(volumesForState) > 0 {
		d.Set("storage_volume", volumesForState)
	}

	var networkInterfaces []map[string]interface{}
	// iterate over the array of svcports
	for i := 0; i < len(instance.Interfaces); i++ {
		row := make(map[string]interface{})
		networkInterface := instance.Interfaces[i]
		row["network_id"] = int(networkInterface.Network.ID)
		row["ip_address"] = networkInterface.IpAddress
		row["ip_mode"] = networkInterface.IpMode
		row["network_interface_type_id"] = networkInterface.NetworkInterfaceTypeId
		networkInterfaces = append(networkInterfaces, row)
	}
	d.Set("network_interface", networkInterfaces)

	return diags
}

// convertToInt converts a value to an int, supporting both int and string types.
func convertToInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, fmt.Errorf("unsupported type %T for conversion to int", value)
	}
}

func resourceMVMInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	instanceGetResp, err := client.GetInstance(toInt64(id), &morpheus.Request{})
	if err != nil {
		log.Printf("API FAILURE: %s - %s", instanceGetResp, err)
		return diag.FromErr(err)
	}
	result := instanceGetResp.Result.(*morpheus.GetInstanceResult)
	fetchInstance := result.Instance

	var instance *morpheus.Instance

	instancePayload := make(map[string]interface{})
	baseInstanceConfigUpdate := false

	if d.HasChange("name") {
		instancePayload["name"] = d.Get("name").(string)
		baseInstanceConfigUpdate = true
	}

	if d.HasChange("description") {
		instancePayload["description"] = d.Get("description").(string)
		baseInstanceConfigUpdate = true
	}

	if d.HasChange("display_name") {
		instancePayload["displayName"] = d.Get("display_name").(string)
		baseInstanceConfigUpdate = true
	}

	if d.HasChange("environment") {
		instancePayload["instanceContext"] = d.Get("environment").(string)
		baseInstanceConfigUpdate = true
	}

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
		baseInstanceConfigUpdate = true
		instancePayload["tags"] = tags
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
		baseInstanceConfigUpdate = true
	}

	if d.HasChanges("labels") {
		instancePayload["labels"] = d.Get("labels")
		baseInstanceConfigUpdate = true
	}

	if baseInstanceConfigUpdate {
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
		instance = result.Instance
	}

	// Resize Instance Configuration
	var resizeChanges bool = false

	resizeInstancePayload := make(map[string]interface{})

	if d.HasChange("plan_id") {
		resizeChanges = true
		resizeInstancePayload["instance"] = map[string]interface{}{
			"plan": map[string]interface{}{
				"id": d.Get("plan_id"),
			},
		}
	}

	if d.HasChange("storage_volume") || d.HasChange("network_interface") {
		resizeChanges = true
		resizeInstancePayload["networkInterfaces"] = parseMVMNetworkInterfaces(d.Get("network_interface").([]interface{}), fetchInstance.Interfaces)
		resizeInstancePayload["volumes"] = parseMVMStorageVolumes(d.Get("storage_volume").([]interface{}), fetchInstance.Volumes)
		resizeInstancePayload["deleteOriginalVolumes"] = true
	}

	// Check if plan, storage, or nics have changed
	if resizeChanges {
		resizeReq := &morpheus.Request{Body: resizeInstancePayload}
		resizeResp, err := client.ResizeInstance(toInt64(id), resizeReq)
		if err != nil {
			log.Printf("API FAILURE: %s - %s", resizeResp, err)
			return diag.FromErr(err)
		}
		result := resizeResp.Result.(*morpheus.UpdateInstanceResult)
		instance = result.Instance
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"resizing", "pending"},
		Target:  []string{"running"},
		Refresh: func() (interface{}, string, error) {
			instanceDetails, err := client.GetInstance(toInt64(id), &morpheus.Request{})
			if err != nil {
				return "", "", err
			}
			result := instanceDetails.Result.(*morpheus.GetInstanceResult)
			instance := result.Instance
			return result, instance.Status, nil
		},
		Timeout:      30 * time.Minute,
		MinTimeout:   1 * time.Minute,
		Delay:        1 * time.Minute,
		PollInterval: 30 * time.Second,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error updating instance: %s", err)
	}

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

	stateConf := &resource.StateChangeConf{
		Pending: []string{"removing", "pendingRemoval", "stopping", "pending", "warning"},
		Target:  []string{"removed"},
		Refresh: func() (interface{}, string, error) {
			instanceDetails, err := client.GetInstance(toInt64(id), &morpheus.Request{})
			if instanceDetails.StatusCode == 404 {
				return "", "removed", nil
			}
			if err != nil {
				return "", "", err
			}
			result := instanceDetails.Result.(*morpheus.GetInstanceResult)
			instance := result.Instance
			return result, instance.Status, nil
		},
		Timeout:      30 * time.Minute,
		MinTimeout:   1 * time.Minute,
		Delay:        1 * time.Minute,
		PollInterval: 30 * time.Second,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error deleting instance: %s", err)
	}

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

func generateMVMInstanceCreatePayload(interfaces []interface{}) []map[string]interface{} {
	var networkInterfaces []map[string]interface{}
	for i := 0; i < len(interfaces); i++ {
		row := make(map[string]interface{})
		item := (interfaces)[i].(map[string]interface{})
		if item["network_id"] != nil {
			row["network"] = map[string]interface{}{
				"id": fmt.Sprintf("network-%d", item["network_id"].(int)),
			}
		}
		if item["ip_address"] != nil {
			row["ipAddress"] = item["ip_address"] //.(string)
		}
		if item["ip_mode"] != nil {
			row["ipMode"] = item["ip_mode"] // .(string)
		}
		if item["network_interface_type_id"] != nil {
			row["networkInterfaceTypeId"] = item["network_interface_type_id"] //.(int)
		}
		networkInterfaces = append(networkInterfaces, row)
	}
	return networkInterfaces
}

func parseMVMNetworkInterfaces(interfaces []interface{}, existingInterfaces []morpheus.NetworkInterface) []map[string]interface{} {
	var networkInterfaces []map[string]interface{}
	//var existingNetworkInterfaceIDs []string

	for i := 0; i < len(interfaces); i++ {
		row := make(map[string]interface{})
		item := (interfaces)[i].(map[string]interface{})

		for _, nic := range existingInterfaces {
			if nic.ID == item["id"] {
				row["id"] = nic.ID
			}
		}
		if item["network_id"] != nil {
			row["network"] = map[string]interface{}{
				"id": fmt.Sprintf("network-%d", item["network_id"].(int)),
			}
		}
		if item["ip_address"] != nil {
			row["ipAddress"] = item["ip_address"] //.(string)
		}
		if item["ip_mode"] != nil {
			row["ipMode"] = item["ip_mode"] // .(string)
		}
		if item["network_interface_type_id"] != nil {
			row["networkInterfaceTypeId"] = item["network_interface_type_id"] //.(int)
		}
		networkInterfaces = append(networkInterfaces, row)
	}
	return networkInterfaces
}

func parseMVMStorageVolumes(volumes []interface{}, existingVolumes Volumes) []map[string]interface{} {
	var storageVolumes []map[string]interface{}
	var existingVolumeUUIDs []string

	// Iterate through the existing storage volumes and fetch the UUID to identify
	// which defined volumes already exist on the instance
	for _, vol := range existingVolumes {
		existingVolumeUUIDs = append(existingVolumeUUIDs, vol.Uuid)
	}

	for i := 0; i < len(volumes); i++ {
		row := make(map[string]interface{})
		item := (volumes)[i].(map[string]interface{})

		for _, vol := range existingVolumes {
			if vol.Uuid == item["uuid"] {
				row["id"] = vol.ID
			}
		}

		// If the defined disk does not already existing
		if !containsString(existingVolumeUUIDs, item["uuid"].(string)) {
			row["id"] = -1
		}

		if item["root"] != nil {
			row["rootVolume"] = item["root"]
		}
		if item["name"] != nil {
			row["name"] = item["name"] // .(string)
		}
		if item["size"] != nil {
			row["size"] = item["size"] // .(int)
		}
		if item["storage_type"] != nil {
			row["storageType"] = item["storage_type"] // .(int)
		}
		if item["datastore_id"] != nil {
			row["datastoreId"] = item["datastore_id"] // .(int)
		}
		storageVolumes = append(storageVolumes, row)
	}
	return storageVolumes // .([]map[string]interface{})
}

type Volumes []struct {
	ID                   interface{} `json:"id"`
	Name                 string      `json:"name"`
	ShortName            string      `json:"shortName"`
	Description          string      `json:"description"`
	ControllerId         int64       `json:"controllerId"`
	ControllerMountPoint string      `json:"controllerMountPoint"`
	Resizeable           interface{} `json:"resizeable"`
	PlanResizable        interface{} `json:"planResizable"`
	Size                 interface{} `json:"size"`
	StorageType          interface{} `json:"storageType"`
	RootVolume           interface{} `json:"rootVolume"`
	UnitNumber           string      `json:"unitNumber"`
	DeviceName           string      `json:"deviceName"`
	DeviceDisplayName    string      `json:"deviceDisplayName"`
	Type                 struct {
		ID   int64  `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"type"`
	TypeId           int64  `json:"typeId"`
	Category         string `json:"category"`
	Status           string `json:"status"`
	StatusMessage    string `json:"statusMessage"`
	ConfigurableIOPS bool   `json:"configurableIOPS"`
	MaxStorage       int64  `json:"maxStorage"`
	DisplayOrder     int64  `json:"displayOrder"`
	MaxIOPS          string `json:"maxIOPS"`
	Uuid             string `json:"uuid"`
	Active           bool   `json:"active"`
	Zone             struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"zone"`
	ZoneId    int64 `json:"zoneId"`
	Datastore struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"datastore"`
	DatastoreId   interface{} `json:"datastoreId"`
	StorageGroup  string      `json:"storageGroup"`
	Namespace     string      `json:"namespace"`
	StorageServer string      `json:"storageServer"`
	Source        string      `json:"source"`
	Owner         struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"owner"`
}
