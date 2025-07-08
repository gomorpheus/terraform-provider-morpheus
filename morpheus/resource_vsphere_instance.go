package morpheus

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceVsphereInstance() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus instance resource.",
		CreateContext: resourceVsphereInstanceCreate,
		ReadContext:   resourceVsphereInstanceRead,
		UpdateContext: resourceVsphereInstanceUpdate,
		DeleteContext: resourceVsphereInstanceDelete,
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
				Description: "The id of type of instance to provision, specify this or 'instance_type_code'",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"instance_type_code": {
				Description:  "The code of type of instance to provision, specify this or 'instance_type_id'",
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"instance_type_id"},
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
				ForceNew:    true,
				Required:    true,
			},
			"resource_pool_id": {
				Description: "The ID of the resource pool to provision the instance to",
				Type:        schema.TypeInt,
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
			},
			"user_group_id": {
				Description: "The id of the user group associated with the instance",
				Type:        schema.TypeInt,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
			},
			"folder_id": {
				Description: "The VMware folder to use when provisioning the instance",
				Type:        schema.TypeInt,
				ForceNew:    true,
				Optional:    true,
			},
			"asset_tag": {
				Description: "The asset tag associated with the instance",
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
							Description: "The ID of the datastore, specify this or datastore_auto_selection",
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
						},
						"datastore_auto_selection": {
							Description:  "Whether to automatically select the datastore, values can be 'auto' or 'autoCluster', specify this or datastore_id",
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"auto", "autoCluster"}, false),
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
			"connection_info": {
				Description: "Connection information for the instance, a list - this is returned by the API",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Description: "The IP address to connect to",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"port": {
							Description: "The port to connect to",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"name": {
							Description: "The name of the connection protocol",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
		CustomizeDiff: customdiff.All(
			volumesCustomizeDiff,
			customdiff.ForceNewIfChange("instance_type_code", func(ctx context.Context, old, new, meta interface{}) bool {
				// We will force a new instance if instance_type_code has a non-zero value, which means that it has been
				// set by the user
				return new.(string) != ""
			}),
			customdiff.ForceNewIfChange("instance_type_id", func(ctx context.Context, old, new, meta interface{}) bool {
				// We will force a new instance if instance_type_id has a non-zero value, which means that it has been
				// set by the user
				return new.(int) != 0
			}),
		),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// volumesCustomizeDiff is a custom diff function to ensure that only one of datastore_id or datastore_auto_selection is set
func volumesCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	if d.HasChange("volumes") {
		volumes := d.Get("volumes").([]interface{})
		for _, volume := range volumes {
			volumeMap := volume.(map[string]interface{})
			// Check if both datastore_id and datastore_auto_selection are set
			// We check for non-zero values of each, zero values (i.e. > 0 or != "") will be ignored
			dataStoreID := volumeMap["datastore_id"].(int)
			dataStoreAutoSelection := volumeMap["datastore_auto_selection"].(string)
			if dataStoreID != 0 && dataStoreAutoSelection != "" {
				return fmt.Errorf("only one of 'datastore_id' or 'datastore_auto_selection' can be set")
			}
		}
	}

	return nil
}

func resourceVsphereInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	group := d.Get("group_id").(int)
	cloud := d.Get("cloud_id").(int)
	name := d.Get("name").(string)

	// Service Plan
	planResp, err := client.GetPlan(int64(d.Get("plan_id").(int)), &morpheus.Request{})
	if err != nil {
		return diag.FromErr(err)
	}
	planResult, ok := planResp.Result.(*morpheus.GetPlanResult)
	if !ok {
		return diag.Errorf("Plan response is not of type *morpheus.GetPlanResult")
	}
	plan := planResult.Plan

	// Instance Type
	// The Schema validation will ensure that only one of "instance_type_code" or "instance_type_id" is set
	// We need the code to create the instance
	instanceTypeCode := d.Get("instance_type_code").(string)
	instanceTypeId := d.Get("instance_type_id").(int)
	if instanceTypeId != 0 {
		instanceTypeResp, err := client.GetInstanceType(int64(d.Get("instance_type_id").(int)), &morpheus.Request{})
		if err != nil {
			return diag.FromErr(err)
		}
		instanceTypeResult, ok := instanceTypeResp.Result.(*morpheus.GetInstanceTypeResult)
		if !ok {
			return diag.Errorf("Instance Type response is not of type *morpheus.GetInstanceTypeResult")
		}
		instanceTypeCode = instanceTypeResult.InstanceType.Code
	}

	// Instance Layout - we only need the ID to create the instance
	instanceLayoutId := int64(d.Get("instance_layout_id").(int))

	// Config
	config := make(map[string]interface{})

	// Resource Pool
	resourcePoolResp, err := client.GetResourcePool(int64(cloud), int64(d.Get("resource_pool_id").(int)), &morpheus.Request{})
	if err != nil {
		return diag.FromErr(err)
	}
	resourcePoolResult, ok := resourcePoolResp.Result.(*morpheus.GetResourcePoolResult)
	if !ok {
		return diag.Errorf("Resource Pool response is not of type *morpheus.GetResourcePoolResult")
	}
	resourcePool := resourcePoolResult.ResourcePool
	config["resourcePoolId"] = resourcePool.ID

	// Custom Options, check for non-zero value
	if d.Get("custom_options") != nil && d.Get("custom_options").(map[string]interface{}) != nil {
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

	// Skip Agent Install
	config["noAgent"] = d.Get("skip_agent_install").(bool)

	// Nested Virtualization
	config["nestedVirtualization"] = d.Get("nested_virtualization").(bool)

	// Folder ID, check for non-zero value
	if d.Get("folder_id") != nil && d.Get("folder_id").(int) != 0 {
		config["vmwareFolderId"] = d.Get("folder_id").(int)
	}

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
			"id": instanceLayoutId,
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
		Delay:        3 * time.Minute,
		PollInterval: 1 * time.Minute,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error creating instance: %s", err)
	}

	// Successfully created resource, now set id
	d.SetId(int64ToString(instance.ID))
	resourceVsphereInstanceRead(ctx, d, meta)
	return diags
}

func resourceVsphereInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	// Special logic for "instance_type_id" and "instance_type_code".
	// Only one of these will have been specified by the user.  We need to figure out which it is
	// and store that value in the State.
	if d.HasChange("instance_type_code") {
		_, newVal := d.GetChange("instance_type_code")
		d.Set("instance_type_code", newVal.(string))
	}
	if d.HasChange("instance_type_id") {
		_, newVal := d.GetChange("instance_type_id")
		d.Set("instance_type_id", newVal.(int))
	}

	d.SetId(int64ToString(instance.ID))
	d.Set("name", instance.Name)
	d.Set("description", instance.Description)
	d.Set("cloud_id", instance.Cloud.ID)
	d.Set("group_id", instance.Group.ID)
	d.Set("instance_layout_id", instance.Layout.ID)
	d.Set("plan_id", instance.Plan.ID)
	d.Set("resource_pool_id", instance.Config["resourcePoolId"])
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
	d.Set("folder_id", instance.Config["vmwareFolderId"])
	if instance.Config["nestedVirtualization"] == "off" {
		d.Set("nested_virtualization", false)
	} else {
		d.Set("nested_virtualization", true)
	}
	d.Set("custom_options", instance.Config["customOptions"])
	d.Set("domain_id", instance.NetworkDomain.Id)

	var networkInterfaces []map[string]interface{}
	// iterate over the array of interfaces
	for i := 0; i < len(instance.Interfaces); i++ {
		row := make(map[string]interface{})
		networkInterface := instance.Interfaces[i]
		row["network_id"] = int(networkInterface.Network.ID)
		row["network_group"] = networkInterface.Network.Group
		row["ip_address"] = networkInterface.IpAddress
		row["ip_mode"] = networkInterface.IpMode
		row["network_interface_type_id"] = networkInterface.NetworkInterfaceTypeId
		networkInterfaces = append(networkInterfaces, row)
	}
	d.Set("interfaces", networkInterfaces)

	var connectionInfo []map[string]interface{}
	// Iterate over the array of connection info
	for i := 0; i < len(instance.ConnectionInfo); i++ {
		row := make(map[string]interface{})
		connection := instance.ConnectionInfo[i]
		row["ip"] = connection.Ip
		row["port"] = connection.Port
		row["name"] = connection.Name
		connectionInfo = append(connectionInfo, row)
	}
	d.Set("connection_info", connectionInfo)

	return diags
}

func resourceVsphereInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	return resourceVsphereInstanceRead(ctx, d, meta)
}

func resourceVsphereInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	// Since in a delete-create cycle we find that the API returns an error that
	// "name must be unique" we will GET the VM instance until such time as the instance
	// isn't available
	stateConf := retry.StateChangeConf{
		Delay:        1 * time.Second,
		Timeout:      5 * time.Minute,
		PollInterval: 1 * time.Second,
		MinTimeout:   1 * time.Second,
		Pending:      []string{"200"},
		Target:       []string{"404"},
		Refresh: func() (interface{}, string, error) {
			resp, err = client.GetInstance(toInt64(id), &morpheus.Request{})
			if err != nil {
				if resp != nil {
					return resp, strconv.Itoa(resp.StatusCode), nil
				}
				return "", "", err
			}

			return resp, strconv.Itoa(resp.StatusCode), nil
		},
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diags
}

func parseNetworkInterfaces(interfaces []interface{}) []map[string]interface{} {
	var networkInterfaces []map[string]interface{}
	for i := 0; i < len(interfaces); i++ {
		row := make(map[string]interface{})
		item := (interfaces)[i].(map[string]interface{})
		if item["network_id"] != nil {
			if item["network_group"].(bool) {
				row["network"] = map[string]interface{}{
					"id": fmt.Sprintf("networkGroup-%d", item["network_id"].(int)),
				}
			} else {
				row["network"] = map[string]interface{}{
					"id": fmt.Sprintf("network-%d", item["network_id"].(int)),
				}
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

func parseStorageVolumes(volumes []interface{}) []map[string]interface{} {
	var storageVolumes []map[string]interface{}
	for i := 0; i < len(volumes); i++ {
		row := make(map[string]interface{})
		item := (volumes)[i].(map[string]interface{})
		if item["id"] != nil {
			row["id"] = item["id"]
		}
		if item["root"] != nil {
			row["rootVolume"] = item["root"]
		}
		if item["name"] != nil {
			row["name"] = item["name"] // .(string)
		}
		// Check for non-zero value of size
		if item["size"] != nil && item["size"].(int) != 0 {
			row["size"] = item["size"] // .(int)
		}
		// Check for non-zero value of size_id
		if item["size_id"] != nil && item["size_id"].(int) != 0 {
			row["sizeId"] = item["size_id"] // .(int)
		}
		// Check for non-zero value of storage_type
		if item["storage_type"] != nil && item["storage_type"].(int) != 0 {
			row["storageType"] = item["storage_type"] // .(int)
		}
		// Check for non-zero value of datastore_id
		if item["datastore_id"] != nil && item["datastore_id"].(int) != 0 {
			row["datastoreId"] = item["datastore_id"] // .(int)
		}
		// If "auto" or "autoCluster" have been specified set the datastoreId to the value
		// Our CustomizeDiff function will ensure that only one of these is set
		if item["datastore_auto_selection"] != nil && item["datastore_auto_selection"].(string) != "" {
			row["datastoreId"] = item["datastore_auto_selection"] // .(string)
		}
		storageVolumes = append(storageVolumes, row)
	}
	return storageVolumes // .([]map[string]interface{})
}
