package morpheus

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"strconv"
	"strings"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
			"description": {
				Description: "The user friendly description of the cloud",
				Type:        schema.TypeString,
				Optional:    true,
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
			"username": {
				Type:        schema.TypeString,
				Description: "The username of the VMware vSphere account",
				Required:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "The password of the VMware vSphere account",
				Required:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256_hash := hex.EncodeToString(h.Sum(nil))
					return strings.ToLower(old) == strings.ToLower(sha256_hash)
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
				Description:  "",
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
				Description: "An arbitrary id used to reference the datacenter for the cloud",
				Optional:    true,
			},
			"guidance": {
				Type:        schema.TypeString,
				Description: "Whether to enable guidance recommendations on the cloud (manual, off)",
				Optional:    true,
				Default:     "off",
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
				ValidateFunc: validation.StringInSlice([]string{"private", "public", ""}, false),
				Default:      "private",
			},
			"tenant_id": {
				Description: "The id of the morpheus tenant the cloud is assigned to",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceVsphereCloudCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	code := d.Get("code").(string)
	location := d.Get("location").(string)
	visibility := d.Get("visibility").(string)
	enabled := d.Get("enabled").(bool)
	automatically_power_on_vms := d.Get("automatically_power_on_vms").(bool)

	config := make(map[string]interface{})
	config["certificateProvider"] = "internal"
	config["apiUrl"] = d.Get("api_url")
	config["username"] = d.Get("username")
	config["password"] = d.Get("password")
	config["datacenter"] = d.Get("datacenter")
	config["apiVersion"] = d.Get("api_version")

	// Select all clusters by passing an
	// empty string to the API
	if d.Get("cluster") == "all" {
		config["cluster"] = ""
	} else {
		config["cluster"] = d.Get("cluster")
	}

	if d.Get("resource_pool") == "all" {
		config["resourcePool"] = ""
	} else {
		config["resourcePool"] = d.Get("resource_pool")
	}

	config["rpcMode"] = d.Get("rpc_mode")
	config["diskStorageType"] = d.Get("storage_type")
	config["datacenterName"] = d.Get("datacenter_id")
	config["applianceUrl"] = d.Get("appliance_url")

	if d.Get("enable_disk_type_selection").(bool) {
		config["enableDiskTypeSelection"] = "on"
	} else {
		config["enableDiskTypeSelection"] = ""
	}

	if d.Get("enable_storage_type_selection").(bool) {
		config["enableStorageTypeSelection"] = "on"
	} else {
		config["enableStorageTypeSelection"] = ""
	}

	if d.Get("enable_network_interface_type_selection").(bool) {
		config["enableNetworkTypeSelection"] = "on"
	} else {
		config["enableNetworkTypeSelection"] = ""
	}

	if d.Get("hide_host_selection").(bool) {
		config["hideHostSelection"] = "on"
	} else {
		config["hideHostSelection"] = ""
	}

	if d.Get("import_existing_vms").(bool) {
		config["importExisting"] = "on"
	} else {
		config["importExisting"] = ""
	}

	if d.Get("enable_hypervisor_console").(bool) {
		config["enableVnc"] = "on"
	} else {
		config["enableVnc"] = ""
	}

	time_zone := d.Get("time_zone").(string)
	agent_install_mode := d.Get("agent_install_mode").(string)
	costing := d.Get("costing").(string)
	keyboard_layout := d.Get("keyboard_layout")
	guidance := d.Get("guidance")

	payload := map[string]interface{}{
		"zone": map[string]interface{}{
			"name":                  name,
			"code":                  code,
			"location":              location,
			"enabled":               enabled,
			"agentMode":             agent_install_mode,
			"autoRecoverPowerState": automatically_power_on_vms,
			"costingMode":           costing,
			"consoleKeymap":         keyboard_layout,
			"description":           d.Get("description").(string),
			"accountId":             d.Get("tenant_id").(string),
			"account": map[string]interface{}{
				"id": d.Get("tenant_id").(string),
			},
			"guidanceMode": guidance,
			"timezone":     time_zone,
			"zoneType": map[string]interface{}{
				"code": "vmware",
			},
			"config":     config,
			"visibility": visibility,
		},
	}

	req := &morpheus.Request{Body: payload}

	resp, err := client.CreateCloud(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.CreateCloudResult)
	cloud := result.Cloud
	// Successfully created resource, now set id
	d.SetId(int64ToString(cloud.ID))
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
			//return diag.FromErr(err)
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
		//return diag.Errorf("Cloud not found in response data.") // should not happen
	} else {
		d.SetId(int64ToString(cloud.ID))
		d.Set("name", cloud.Name)
		d.Set("code", cloud.Code)
		d.Set("location", cloud.Location)
		d.Set("visibility", cloud.Visibility)
		d.Set("enabled", cloud.Enabled)
		d.Set("tenant_id", strconv.Itoa(int(cloud.AccountID)))
		d.Set("api_url", cloud.Config.APIUrl)
		d.Set("username", cloud.Config.Username)
		d.Set("password", cloud.Config.PasswordHash)
		d.Set("api_version", cloud.Config.APIVersion)
		d.Set("datacenter", cloud.Config.Datacenter)
		d.Set("cluster", cloud.Config.Cluster)
		d.Set("rpc_mode", cloud.Config.RPCMode)

		if cloud.Config.Cluster == "" {
			d.Set("cluster", "all")
		} else {
			d.Set("cluster", cloud.Config.Cluster)
		}

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
		d.Set("guidance", cloud.GuidanceMode)
		d.Set("costing", cloud.CostingMode)
		d.Set("agent_install_mode", cloud.AgentMode)
		return diags
	}
}

func resourceVsphereCloudUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	name := d.Get("name").(string)
	code := d.Get("code").(string)
	location := d.Get("location").(string)
	enabled := d.Get("enabled").(bool)
	automatically_power_on_vms := d.Get("automatically_power_on_vms").(bool)
	visibility := d.Get("visibility").(string)

	config := make(map[string]interface{})
	config["certificateProvider"] = "internal"
	config["apiUrl"] = d.Get("api_url")
	config["username"] = d.Get("username")

	if d.HasChange("password") {
		config["password"] = d.Get("password")
	}

	config["datacenter"] = d.Get("datacenter")
	config["apiVersion"] = d.Get("api_version")
	// Select all clusters by passing an
	// empty string to the API
	if d.Get("cluster") == "all" {
		config["cluster"] = ""
	} else {
		config["cluster"] = d.Get("cluster")
	}
	config["rpcMode"] = d.Get("rpc_mode")
	config["diskStorageType"] = d.Get("storage_type")
	config["datacenterName"] = d.Get("datacenter_id")
	config["applianceUrl"] = d.Get("appliance_url")
	if d.Get("hide_host_selection") == nil {
		config["hideHostSelection"] = ""
	} else {
		config["hideHostSelection"] = "on"
	}

	if d.Get("import_existing") == nil {
		config["importExisting"] = ""
	} else {
		config["importExisting"] = "on"
	}

	if d.Get("enable_hypervisor_console") == nil {
		config["enableVnc"] = ""
	} else {
		config["enableVnc"] = "on"
	}

	if d.Get("enable_disk_type_selection").(bool) {
		config["enableDiskTypeSelection"] = "on"
	} else {
		config["enableDiskTypeSelection"] = ""
	}

	if d.Get("enable_storage_type_selection").(bool) {
		config["enableStorageTypeSelection"] = "on"
	} else {
		config["enableStorageTypeSelection"] = ""
	}

	if d.Get("enable_network_interface_type_selection").(bool) {
		config["enableNetworkTypeSelection"] = "on"
	} else {
		config["enableNetworkTypeSelection"] = ""
	}

	agent_install_mode := d.Get("agent_install_mode").(string)
	costing := d.Get("costing").(string)
	keyboard_layout := d.Get("keyboard_layout")
	guidance := d.Get("guidance")
	time_zone := d.Get("time_zone")

	payload := map[string]interface{}{
		"zone": map[string]interface{}{
			"name":                  name,
			"code":                  code,
			"location":              location,
			"enabled":               enabled,
			"agentMode":             agent_install_mode,
			"autoRecoverPowerState": automatically_power_on_vms,
			"description":           d.Get("description").(string),
			"accountId":             d.Get("tenant_id").(string),
			"account": map[string]interface{}{
				"id": d.Get("tenant_id").(string),
			},
			"costingMode":   costing,
			"consoleKeymap": keyboard_layout,
			"guidanceMode":  guidance,
			"timezone":      time_zone,
			"zoneType": map[string]interface{}{
				"code": "vmware",
			},
			"config":     config,
			"visibility": visibility,
		},
	}

	req := &morpheus.Request{Body: payload}
	resp, err := client.UpdateCloud(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateCloudResult)
	cloud := result.Cloud
	// Successfully updated resource, now set id
	d.SetId(int64ToString(cloud.ID))
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
