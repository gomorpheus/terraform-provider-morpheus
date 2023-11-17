package morpheus

import (
	"context"
	"log"
	"time"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceStandardCloud() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus cloud resource.",
		CreateContext: resourceStandardCloudCreate,
		ReadContext:   resourceStandardCloudRead,
		UpdateContext: resourceStandardCloudUpdate,
		DeleteContext: resourceStandardCloudDelete,

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
			"import_existing_vms": {
				Type:        schema.TypeBool,
				Description: "Whether to import existing virtual machines",
				Optional:    true,
				Default:     false,
			},
			"enable_network_interface_type_selection": {
				Type:        schema.TypeBool,
				Description: "Whether to enable the user to select the network interface type during provisioning",
				Optional:    true,
				Default:     false,
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
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			/* Awaiting SDK Support
			"logo_image_name": {
				Type:        schema.TypeString,
				Description: "The file name of the cloud logo image",
				Optional:    true,
				Computed:    true,
			},
			"logo_image_path": {
				Type:        schema.TypeString,
				Description: "The file path of the cloud logo image including the file name",
				Optional:    true,
				Computed:    true,
			},
			"dark_logo_image_name": {
				Type:        schema.TypeString,
				Description: "The file name of the cloud dark mode logo image",
				Optional:    true,
				Computed:    true,
			},
			"dark_logo_image_path": {
				Type:        schema.TypeString,
				Description: "The file path of the cloud dark mode logo image including the file name",
				Optional:    true,
				Computed:    true,
			},
			*/
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceStandardCloudCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	account["id"] = d.Get("tenant_id").(int)
	cloud["account"] = account
	cloud["accountId"] = d.Get("tenant_id").(int)
	// Enabled
	cloud["enabled"] = d.Get("enabled").(bool)
	// Automatically Power On VMs
	cloud["autoRecoverPowerState"] = d.Get("automatically_power_on_vms").(bool)

	config := make(map[string]interface{})

	// Inventory Existing Instances
	if d.Get("import_existing_vms").(bool) {
		config["importExisting"] = "on"
	} else {
		config["importExisting"] = ""
	}

	// Enable Network Interface Type Selection
	if d.Get("enable_network_interface_type_selection").(bool) {
		config["enableNetworkTypeSelection"] = "on"
	} else {
		config["enableNetworkTypeSelection"] = ""
	}

	config["certificateProvider"] = "internal"

	// Domain
	// Appliance URL
	config["applianceUrl"] = d.Get("appliance_url")
	// Time Zone
	cloud["timezone"] = d.Get("time_zone").(string)
	// Datacenter ID
	config["datacenterName"] = d.Get("datacenter_id")
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
	cloudType["code"] = "standard"
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
		PollInterval: 30 * time.Second,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error creating cloud: %s", err)
	}

	/* Awaiting SDK support
	var filePayloads []*morpheus.FilePayload

	if d.Get("logo_image_path") != "" && d.Get("logo_image_name") != "" {
		data, err := os.ReadFile(d.Get("logo_image_path").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		filePayload := &morpheus.FilePayload{
			ParameterName: "logo",
			FileName:      d.Get("logo_image_name").(string),
			FileContent:   data,
		}
		filePayloads = append(filePayloads, filePayload)
	}
	if d.Get("dark_logo_image_path") != "" && d.Get("dark_logo_image_name") != "" {
		darkLogoData, err := os.ReadFile(d.Get("dark_logo_image_path").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		darkLogoPayload := &morpheus.FilePayload{
			ParameterName: "darkLogo",
			FileName:      d.Get("dark_logo_image_name").(string),
			FileContent:   darkLogoData,
		}
		filePayloads = append(filePayloads, darkLogoPayload)
	}

	response, err := client.UpdateCloudLogo(cloudOutput.ID, filePayloads, &morpheus.Request{})
	if err != nil {
		log.Printf("API FAILURE: %s - %s", response, err)
	}
	log.Printf("API RESPONSE: %s", response)
	*/
	// Successfully created resource, now set id
	d.SetId(int64ToString(cloudOutput.ID))
	resourceStandardCloudRead(ctx, d, meta)
	return diags
}

func resourceStandardCloudRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		d.Set("visibility", cloud.Visibility)
		d.Set("tenant_id", int(cloud.AccountID))
		d.Set("enabled", cloud.Enabled)
		d.Set("automatically_power_on_vms", cloud.AutoRecoverPowerState)

		if cloud.Config.ImportExisting == "on" {
			d.Set("import_existing_vms", true)
		} else {
			d.Set("import_existing_vms", false)
		}

		if cloud.Config.EnableNetworkTypeSelection == "on" {
			d.Set("enable_network_interface_type_selection", true)
		} else {
			d.Set("enable_network_interface_type_selection", false)
		}
		d.Set("appliance_url", cloud.Config.ApplianceUrl)
		d.Set("time_zone", cloud.TimeZone)
		d.Set("datacenter_id", cloud.Config.DatacenterName)
		d.Set("guidance", cloud.GuidanceMode)
		d.Set("costing", cloud.CostingMode)
		d.Set("agent_install_mode", cloud.AgentMode)
		/* Awaiting SDK Support
		imagePath := strings.Split(cloud.ImagePath, "/")
		opt := strings.Replace(imagePath[len(imagePath)-1], "_original", "", 1)
		d.Set("logo_image_name", opt)
		darkImagePath := strings.Split(catalogItem.DarkImagePath, "/")
		darkOpt := strings.Replace(darkImagePath[len(darkImagePath)-1], "_original", "", 1)
		d.Set("dark_logo_image_name", darkOpt)
		*/
		return diags
	}
}

func resourceStandardCloudUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	// Inventory Existing Instances
	if d.Get("import_existing_vms").(bool) {
		config["importExisting"] = "on"
	} else {
		config["importExisting"] = ""
	}

	// Enable Network Interface Type Selection
	if d.Get("enable_network_interface_type_selection").(bool) {
		config["enableNetworkTypeSelection"] = "on"
	} else {
		config["enableNetworkTypeSelection"] = ""
	}
	// Domain
	// Appliance URL
	config["applianceUrl"] = d.Get("appliance_url")
	// Time Zone
	cloud["timezone"] = d.Get("time_zone").(string)
	// Datacenter ID
	config["datacenterName"] = d.Get("datacenter_id")
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
	cloudType["code"] = "standard"
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

	/* Awaiting SDK Support
	var filePayloads []*morpheus.FilePayload

	if d.HasChange("logo_image_path") || d.HasChange("logo_image_name") {
		data, err := os.ReadFile(d.Get("logo_image_path").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		filePayload := &morpheus.FilePayload{
			ParameterName: "logo",
			FileName:      d.Get("logo_image_name").(string),
			FileContent:   data,
		}
		filePayloads = append(filePayloads, filePayload)
	}
	if d.HasChange("dark_logo_image_path") || d.HasChange("dark_logo_image_name") {
		darkLogoData, err := os.ReadFile(d.Get("dark_logo_image_path").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		darkLogoPayload := &morpheus.FilePayload{
			ParameterName: "darkLogo",
			FileName:      d.Get("dark_logo_image_name").(string),
			FileContent:   darkLogoData,
		}
		filePayloads = append(filePayloads, darkLogoPayload)
	}

	response, err := client.UpdateCatalogItemLogo(catalogItemResult.ID, filePayloads, &morpheus.Request{})
	if err != nil {
		log.Printf("API FAILURE: %s - %s", response, err)
	}
	log.Printf("API RESPONSE: %s", response)
	*/

	// Successfully updated resource, now set id
	d.SetId(int64ToString(cloudOutput.ID))
	return resourceStandardCloudRead(ctx, d, meta)
}

func resourceStandardCloudDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
