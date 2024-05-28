package morpheus

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceVsphereCloud() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus cloud resource.",
		CreateContext: resourceVsphereCloudCreate,
		ReadContext:   resourceVsphereCloudRead,
		UpdateContext: resourceVsphereCloudUpdate,
		DeleteContext: resourceVsphereCloudDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the cloud",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "A unique name scoped to your account for the cloud",
				Type:        schema.TypeString,
				Required:    true,
			},
			"code": {
				Description: "Optional code for use with policies",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"location": {
				Description: "Optional location for your cloud",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"enabled": {
				Description: "Determines whether the cloud is active or not",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"automatically_power_on_vms": {
				Description: "Determines whether to automatically power on cloud virtual machines",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"api_url": {
				Type:        schema.TypeString,
				Description: "The SDK URL of the vCenter server (https://vcenter.morpheus.local/sdk)",
				Required:    true,
			},
			"credential_id": {
				Description:   "The ID of the credential store entry used for authentication",
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"username", "password"},
			},
			"username": {
				Type:          schema.TypeString,
				Description:   "The username of the VMware vSphere account",
				Optional:      true,
				ConflictsWith: []string{"credential_id"},
			},
			"password": {
				Type:          schema.TypeString,
				Description:   "The password of the VMware vSphere account",
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"credential_id"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256_hash := hex.EncodeToString(h.Sum(nil))
					return strings.EqualFold(old, sha256_hash)
				},
			},
			"api_version": {
				Type:        schema.TypeString,
				Description: "The SDK URL of the vCenter server (https://vcenter.morpheus.local/sdk)",
				Required:    true,
			},
			"datacenter": {
				Type:        schema.TypeString,
				Description: "The vSphere datacenter to add",
				Required:    true,
			},
			"cluster": {
				Type:        schema.TypeString,
				Description: "The name of the vSphere cluster",
				Optional:    true,
				Default:     "all",
			},
			"resource_pool": {
				Type:        schema.TypeString,
				Description: "The name of the vSphere resource pool",
				Optional:    true,
			},
			"rpc_mode": {
				Type:         schema.TypeString,
				Description:  "The method for interacting with cloud workloads (guestexec (VMware Tools) or rpc (SSH/WinRM))",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"guestexec", "rpc", ""}, true),
				Default:      "guestexec",
			},
			"hide_host_selection": {
				Type:        schema.TypeBool,
				Description: "Whether to hide the ability to select the vSphere host from the user during provisioning",
				Optional:    true,
				Default:     false,
			},
			"import_existing_vms": {
				Type:        schema.TypeBool,
				Description: "Whether to import existing virtual machines",
				Optional:    true,
				Default:     false,
			},
			"enable_hypervisor_console": {
				Type:        schema.TypeBool,
				Description: "Whether to enable VNC access",
				Optional:    true,
				Default:     false,
			},
			"keyboard_layout": {
				Type:        schema.TypeString,
				Description: "The keyboard layout",
				Optional:    true,
				Default:     "us",
			},
			"enable_disk_type_selection": {
				Type:        schema.TypeBool,
				Description: "Whether to enable the user to select the disk type during provisioning",
				Optional:    true,
				Default:     false,
			},
			"enable_storage_type_selection": {
				Type:        schema.TypeBool,
				Description: "Whether to enable the user to select the storage type during provisioning",
				Optional:    true,
				Default:     false,
			},
			"enable_network_interface_type_selection": {
				Type:        schema.TypeBool,
				Description: "Whether to enable the user to select the network interface type during provisioning",
				Optional:    true,
				Default:     false,
			},
			"storage_type": {
				Type:         schema.TypeString,
				Description:  "The default vSphere VMDK type for virtual machines (thin, thick, thickEager)",
				ValidateFunc: validation.StringInSlice([]string{"thin", "thick", "thickEager"}, true),
				Optional:     true,
				Default:      "thin",
			},
			"appliance_url": {
				Type:        schema.TypeString,
				Description: "The URL used by workloads provisioned in the cloud for interacting with the Morpheus appliance",
				Optional:    true,
			},
			"time_zone": {
				Type:        schema.TypeString,
				Description: "The time zone for the cloud",
				Optional:    true,
			},
			"datacenter_id": {
				Type:        schema.TypeString,
				Description: "A custom id used to reference the datacenter for the cloud",
				Optional:    true,
			},
			"config_management_integration_id": {
				Type:        schema.TypeString,
				Description: "The id of the configuration management ingegration associated with the vSphere cloud",
				Optional:    true,
				Computed:    true,
			},
			"guidance": {
				Type:         schema.TypeString,
				Description:  "Whether to enable guidance recommendations on the cloud (manual, off)",
				ValidateFunc: validation.StringInSlice([]string{"manual", "off", ""}, false),
				Optional:     true,
				Default:      "off",
			},
			"costing": {
				Type:         schema.TypeString,
				Description:  "Whether to enable costing on the cloud (off, costing)",
				ValidateFunc: validation.StringInSlice([]string{"off", "costing", ""}, false),
				Optional:     true,
				Default:      "off",
			},
			"agent_install_mode": {
				Type:         schema.TypeString,
				Description:  "The method used to install the Morpheus agent on virtual machines provisioned in the cloud (ssh, cloudInit)",
				ValidateFunc: validation.StringInSlice([]string{"ssh", "cloudInit", ""}, false),
				Optional:     true,
				Default:      "cloudInit",
			},
			"visibility": {
				Description:  "Determines whether the cloud is visible in sub-tenants or not",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"private", "public"}, false),
				Default:      "private",
			},
			"tenant_id": {
				Description: "The id of the morpheus tenant the cloud is assigned to",
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

func resourceVsphereCloudCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	cloud := make(map[string]interface{})
	// Name
	cloud["name"] = d.Get("name").(string)
	// Code
	cloud["code"] = d.Get("code").(string)
	// Location
	cloud["location"] = d.Get("location").(string)
	// Visibility
	cloud["visibility"] = d.Get("visibility").(string)
	// Tenant
	account := make(map[string]interface{})
	account["id"] = d.Get("tenant_id").(string)
	cloud["account"] = account
	cloud["accountId"] = d.Get("tenant_id").(string)
	// Enabled
	cloud["enabled"] = d.Get("enabled").(bool)
	// Automatically Power On VMs
	cloud["autoRecoverPowerState"] = d.Get("automatically_power_on_vms").(bool)

	config := make(map[string]interface{})
	config["certificateProvider"] = "internal"
	// API URL
	config["apiUrl"] = d.Get("api_url")

	if d.Get("credential_id").(int) != 0 {
		credential := make(map[string]interface{})
		credential["type"] = "username-password"
		credential["id"] = d.Get("credential_id").(int)
		cloud["credential"] = credential
	} else {
		credential := make(map[string]interface{})
		credential["type"] = "local"
		cloud["credential"] = credential
		config["username"] = d.Get("username")
		config["password"] = d.Get("password")
	}
	// Version
	config["apiVersion"] = d.Get("api_version")
	// VDC
	config["datacenter"] = d.Get("datacenter")
	// Cluster
	// Select all clusters by passing an empty string to the API
	if d.Get("cluster") == "all" {
		config["cluster"] = ""
	} else {
		config["cluster"] = d.Get("cluster")
	}
	// Resource Pool
	// Select all resource pools by passing an empty string to the API
	if d.Get("resource_pool") == "all" {
		config["resourcePool"] = ""
	} else {
		config["resourcePool"] = d.Get("resource_pool")
	}
	// RPC Mode
	config["rpcMode"] = d.Get("rpc_mode")
	// Hide Host Selection From Users
	if d.Get("hide_host_selection").(bool) {
		config["hideHostSelection"] = "on"
	} else {
		config["hideHostSelection"] = ""
	}
	// Inventory Existing Instances
	if d.Get("import_existing_vms").(bool) {
		config["importExisting"] = "on"
	} else {
		config["importExisting"] = ""
	}
	// Enable Hypervisor Console
	if d.Get("enable_hypervisor_console").(bool) {
		config["enableVnc"] = "on"
	} else {
		config["enableVnc"] = ""
	}
	// Keyboard Layout
	cloud["consoleKeymap"] = d.Get("keyboard_layout").(string)
	// Enable Disk Type Selection
	if d.Get("enable_disk_type_selection").(bool) {
		config["enableDiskTypeSelection"] = "on"
	} else {
		config["enableDiskTypeSelection"] = ""
	}
	// Enable Storage Type Selection
	if d.Get("enable_storage_type_selection").(bool) {
		config["enableStorageTypeSelection"] = "on"
	} else {
		config["enableStorageTypeSelection"] = ""
	}
	// Enable Network Interface Type Selection
	if d.Get("enable_network_interface_type_selection").(bool) {
		config["enableNetworkTypeSelection"] = "on"
	} else {
		config["enableNetworkTypeSelection"] = ""
	}
	// Storage Type
	config["diskStorageType"] = d.Get("storage_type")
	// Domain
	// Appliance URL
	config["applianceUrl"] = d.Get("appliance_url")
	// Time Zone
	cloud["timezone"] = d.Get("time_zone").(string)
	// Datacenter ID
	config["datacenterName"] = d.Get("datacenter_id")
	// Config Management
	config["configManagementId"] = d.Get("config_management_integration_id").(string)
	// Network Mode
	// Local Firewall
	// Security Server
	// Backup Provider
	// Replication Provider
	// Guidance
	cloud["guidanceMode"] = d.Get("guidance").(string)
	// Costing
	cloud["costingMode"] = d.Get("costing").(string)
	// CMDB
	// CMDB Discovery
	// Agent Install Mode
	cloud["agentMode"] = d.Get("agent_install_mode").(string)
	// VDI Gatway
	cloudType := make(map[string]interface{})
	cloudType["code"] = "vmware"
	cloud["zoneType"] = cloudType

	cloud["config"] = config

	payload := map[string]interface{}{
		"zone": cloud,
	}
	req := &morpheus.Request{Body: payload}

	resp, err := client.CreateCloud(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.CreateCloudResult)
	cloudOutput := result.Cloud

	stateConf := &resource.StateChangeConf{
		Pending: []string{"initializing", "syncing"},
		Target:  []string{"ok"},
		Refresh: func() (interface{}, string, error) {
			cloudDetails, err := client.GetCloud(cloudOutput.ID, &morpheus.Request{})
			if err != nil {
				return "", "", err
			}
			result := cloudDetails.Result.(*morpheus.GetCloudResult)
			cloudStatus := result.Cloud
			return result, cloudStatus.Status, nil
		},
		Timeout:      1 * time.Hour,
		MinTimeout:   1 * time.Minute,
		PollInterval: 1 * time.Minute,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error creating cloud: %s", err)
	}

	// Successfully created resource, now set id
	d.SetId(int64ToString(cloudOutput.ID))
	resourceVsphereCloudRead(ctx, d, meta)
	return diags
}

func resourceVsphereCloudRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindCloudByName(name)
	} else if id != "" {
		resp, err = client.GetCloud(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Cloud cannot be read without name or id")
	}
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)
			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)

	// store resource data
	result := resp.Result.(*morpheus.GetCloudResult)
	cloud := result.Cloud
	if cloud == nil {
		d.SetId("")
		return diags
	} else {
		d.SetId(int64ToString(cloud.ID))
		d.Set("name", cloud.Name)
		d.Set("code", cloud.Code)
		d.Set("location", cloud.Location)
		d.Set("enabled", cloud.Enabled)
		d.Set("automatically_power_on_vms", cloud.AutoRecoverPowerState)
		d.Set("api_url", cloud.Config.APIUrl)
		if cloud.Credential.ID == 0 {
			d.Set("username", cloud.Config.Username)
			d.Set("password", cloud.Config.PasswordHash)
		} else {
			d.Set("credential_id", cloud.Credential.ID)
		}
		d.Set("api_version", cloud.Config.APIVersion)
		d.Set("datacenter", cloud.Config.Datacenter)
		if cloud.Config.Cluster == "" {
			d.Set("cluster", "all")
		} else {
			d.Set("cluster", cloud.Config.Cluster)
		}
		if cloud.Config.ResourcePool == "" {
			d.Set("resource_pool", "all")
		} else {
			d.Set("resource_pool", cloud.Config.ResourcePool)
		}
		d.Set("rpc_mode", cloud.Config.RPCMode)

		if cloud.Config.HideHostSelection == "on" {
			d.Set("hide_host_selection", true)
		} else {
			d.Set("hide_host_selection", false)
		}

		if cloud.Config.ImportExisting == "on" {
			d.Set("import_existing_vms", true)
		} else {
			d.Set("import_existing_vms", false)
		}

		if cloud.Config.EnableVNC == "on" {
			d.Set("enable_hypervisor_console", true)
		} else {
			d.Set("enable_hypervisor_console", false)
		}

		d.Set("keyboard_layout", cloud.ConsoleKeymap)

		if cloud.Config.EnableDiskTypeSelection == "on" {
			d.Set("enable_disk_type_selection", true)
		} else {
			d.Set("enable_disk_type_selection", false)
		}

		if cloud.Config.EnableStorageTypeSelection == "on" {
			d.Set("enable_storage_type_selection", true)
		} else {
			d.Set("enable_storage_type_selection", false)
		}
		if cloud.Config.EnableNetworkTypeSelection == "on" {
			d.Set("enable_network_interface_type_selection", true)
		} else {
			d.Set("enable_network_interface_type_selection", false)
		}
		d.Set("storage_type", cloud.Config.DiskStorageType)
		d.Set("appliance_url", cloud.Config.ApplianceUrl)
		d.Set("time_zone", cloud.TimeZone)
		d.Set("datacenter_id", cloud.Config.DatacenterName)
		d.Set("config_management_integration_id", cloud.Config.ConfigManagementID)
		d.Set("guidance", cloud.GuidanceMode)
		d.Set("costing", cloud.CostingMode)
		d.Set("agent_install_mode", cloud.AgentMode)
		d.Set("visibility", cloud.Visibility)
		d.Set("tenant_id", strconv.Itoa(int(cloud.AccountID)))
		return diags
	}
}

func resourceVsphereCloudUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	cloud := make(map[string]interface{})
	// Name
	cloud["name"] = d.Get("name").(string)
	// Code
	cloud["code"] = d.Get("code").(string)
	// Location
	cloud["location"] = d.Get("location").(string)
	// Visibility
	cloud["visibility"] = d.Get("visibility").(string)
	// Tenant
	account := make(map[string]interface{})
	account["id"] = d.Get("tenant_id").(string)
	cloud["account"] = account
	cloud["accountId"] = d.Get("tenant_id").(string)
	// Enabled
	cloud["enabled"] = d.Get("enabled").(bool)
	// Automatically Power On VMs
	cloud["autoRecoverPowerState"] = d.Get("automatically_power_on_vms").(bool)

	config := make(map[string]interface{})
	config["certificateProvider"] = "internal"
	// API URL
	config["apiUrl"] = d.Get("api_url")

	if d.Get("credential_id").(int) != 0 {
		credential := make(map[string]interface{})
		credential["type"] = "username-password"
		credential["id"] = d.Get("credential_id").(int)
		cloud["credential"] = credential
	} else {
		credential := make(map[string]interface{})
		credential["type"] = "local"
		cloud["credential"] = credential
		if d.HasChange("username") {
			config["username"] = d.Get("username")
		}
		if d.HasChange("password") {
			config["password"] = d.Get("password")
		}
	}
	// Version
	config["apiVersion"] = d.Get("api_version")
	// VDC
	config["datacenter"] = d.Get("datacenter")
	// Cluster
	// Select all clusters by passing an empty string to the API
	if d.Get("cluster") == "all" {
		config["cluster"] = ""
	} else {
		config["cluster"] = d.Get("cluster")
	}
	// Resource Pool
	// Select all resource pools by passing an empty string to the API
	if d.Get("resource_pool") == "all" {
		config["resourcePool"] = ""
	} else {
		config["resourcePool"] = d.Get("resource_pool")
	}
	// RPC Mode
	config["rpcMode"] = d.Get("rpc_mode")
	// Hide Host Selection From Users
	if d.Get("hide_host_selection").(bool) {
		config["hideHostSelection"] = "on"
	} else {
		config["hideHostSelection"] = ""
	}
	// Inventory Existing Instances
	if d.Get("import_existing_vms").(bool) {
		config["importExisting"] = "on"
	} else {
		config["importExisting"] = ""
	}
	// Enable Hypervisor Console
	if d.Get("enable_hypervisor_console").(bool) {
		config["enableVnc"] = "on"
	} else {
		config["enableVnc"] = ""
	}
	// Keyboard Layout
	cloud["consoleKeymap"] = d.Get("keyboard_layout").(string)
	// Enable Disk Type Selection
	if d.Get("enable_disk_type_selection").(bool) {
		config["enableDiskTypeSelection"] = "on"
	} else {
		config["enableDiskTypeSelection"] = ""
	}
	// Enable Storage Type Selection
	if d.Get("enable_storage_type_selection").(bool) {
		config["enableStorageTypeSelection"] = "on"
	} else {
		config["enableStorageTypeSelection"] = ""
	}
	// Enable Network Interface Type Selection
	if d.Get("enable_network_interface_type_selection").(bool) {
		config["enableNetworkTypeSelection"] = "on"
	} else {
		config["enableNetworkTypeSelection"] = ""
	}
	// Storage Type
	config["diskStorageType"] = d.Get("storage_type")
	// Domain
	// Appliance URL
	config["applianceUrl"] = d.Get("appliance_url")
	// Time Zone
	cloud["timezone"] = d.Get("time_zone").(string)
	// Datacenter ID
	config["datacenterName"] = d.Get("datacenter_id")
	// Config Management
	if d.HasChange("config_management_integration_id") {
		config["configManagementId"] = d.Get("config_management_integration_id").(string)
	}
	// Network Mode
	// Local Firewall
	// Security Server
	// Backup Provider
	// Replication Provider
	// Guidance
	cloud["guidanceMode"] = d.Get("guidance").(string)
	// Costing
	cloud["costingMode"] = d.Get("costing").(string)
	// CMDB
	// CMDB Discovery
	// Agent Install Mode
	cloud["agentMode"] = d.Get("agent_install_mode").(string)
	// VDI Gatway
	cloudType := make(map[string]interface{})
	cloudType["code"] = "vmware"
	cloud["zoneType"] = cloudType

	cloud["config"] = config

	payload := map[string]interface{}{
		"zone": cloud,
	}

	req := &morpheus.Request{Body: payload}
	resp, err := client.UpdateCloud(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateCloudResult)
	cloudOutput := result.Cloud
	// Successfully updated resource, now set id
	d.SetId(int64ToString(cloudOutput.ID))
	return resourceVsphereCloudRead(ctx, d, meta)
}

func resourceVsphereCloudDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteCloud(toInt64(id), req)
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
