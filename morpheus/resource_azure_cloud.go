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

func resourceAzureCloud() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus Azure cloud resource.",
		CreateContext: resourceAzureCloudCreate,
		ReadContext:   resourceAzureCloudRead,
		UpdateContext: resourceAzureCloudUpdate,
		DeleteContext: resourceAzureCloudDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(45 * time.Minute),
			Delete: schema.DefaultTimeout(45 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the cloud",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The name of the cloud integration",
				Type:        schema.TypeString,
				Required:    true,
			},
			"code": {
				Description: "Optional code for use with policies",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"location": {
				Description: "Optional location for the cloud",
				Type:        schema.TypeString,
				Optional:    true,
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
				Computed:    true,
			},
			"cloud_type": {
				Type:         schema.TypeString,
				Description:  "The Azure cloud type (global, usgov, german, china)",
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"global", "usgov", "german", "china"}, false),
			},
			"azure_subscription_id": {
				Type:        schema.TypeString,
				Description: "The Azure subscription ID used for authentication",
				Required:    true,
			},
			"azure_tenant_id": {
				Type:        schema.TypeString,
				Description: "The Azure tenant ID used for authentication",
				Required:    true,
			},
			"credential_id": {
				Description: "The ID of the credential store entry used for authentication",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"azure_client_id": {
				Type:        schema.TypeString,
				Description: "The Azure client ID used for authentication",
				Optional:    true,
				Computed:    true,
			},
			"azure_client_secret": {
				Type:        schema.TypeString,
				Description: "The Azure client secret used for authentication",
				Optional:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256_hash := hex.EncodeToString(h.Sum(nil))
					return strings.EqualFold(old, sha256_hash)
				},
				RequiredWith: []string{"azure_client_id"},
			},
			"region": {
				Type:        schema.TypeString,
				Description: "The Azure region associated with the cloud integration",
				Required:    true,
			},
			"resource_group": {
				Type:        schema.TypeString,
				Description: "The Azure resource group associated with the cloud integration",
				Optional:    true,
				Computed:    true,
			},
			"import_existing_instances": {
				Type:        schema.TypeBool,
				Description: "Whether to import existing instances",
				Optional:    true,
				Default:     false,
			},
			//			"account_type": {
			//				Type:        schema.TypeString,
			//				Description: "The Azure cloud account type (Standard, EA, SP)",
			//				Optional:    true,
			//				Default:     "Standard",
			//			},
			"rpc_mode": {
				Type:         schema.TypeString,
				Description:  "The method for interacting with cloud workloads (guestexec (Azure Run Command) or rpc (SSH/WinRM))",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"guestexec", "rpc"}, true),
				Default:      "guestexec",
			},
			"appliance_url": {
				Type:        schema.TypeString,
				Description: "The URL used by workloads provisioned in the cloud for interacting with the Morpheus server",
				Optional:    true,
				Computed:    true,
			},
			"time_zone": {
				Type:        schema.TypeString,
				Description: "The time zone for the cloud",
				Optional:    true,
				Computed:    true,
			},
			"datacenter_id": {
				Type:        schema.TypeString,
				Description: "An arbitrary id used to reference the datacenter for the cloud",
				Optional:    true,
				Computed:    true,
			},
			"config_management_integration_id": {
				Type:        schema.TypeString,
				Description: "The id of the configuration management integration associated with the Azure cloud",
				Optional:    true,
				Computed:    true,
			},
			"guidance": {
				Type:         schema.TypeString,
				Description:  "Whether to enable guidance recommendations on the cloud (manual, off)",
				ValidateFunc: validation.StringInSlice([]string{"manual", "off"}, false),
				Optional:     true,
				Computed:     true,
			},
			"costing": {
				Type:         schema.TypeString,
				Description:  "Whether to enable costing on the cloud (off, costing, reservations, full)",
				ValidateFunc: validation.StringInSlice([]string{"off", "costing", "reservations", "full"}, false),
				Optional:     true,
				Computed:     true,
			},
			"agent_install_mode": {
				Type:         schema.TypeString,
				Description:  "The method used to install the Morpheus agent on instances provisioned in the cloud (ssh, cloudInit)",
				ValidateFunc: validation.StringInSlice([]string{"ssh", "cloudInit", ""}, false),
				Optional:     true,
				Computed:     true,
			},
		},
	}
}

func resourceAzureCloudCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	cloud := make(map[string]interface{})
	cloud["name"] = d.Get("name").(string)
	cloud["code"] = d.Get("code").(string)
	cloud["location"] = d.Get("location").(string)
	cloud["visibility"] = d.Get("visibility").(string)

	account := make(map[string]interface{})
	account["id"] = d.Get("tenant_id").(string)
	cloud["account"] = account
	cloud["accountId"] = d.Get("tenant_id").(string)

	cloud["enabled"] = d.Get("enabled").(bool)
	cloud["autoRecoverPowerState"] = d.Get("automatically_power_on_vms").(bool)

	config := make(map[string]interface{})

	config["cloudType"] = d.Get("cloud_type").(string)
	config["subscriberId"] = d.Get("azure_subscription_id").(string)
	config["tenantId"] = d.Get("azure_tenant_id").(string)

	if d.Get("credential_id").(int) != 0 {
		credential := make(map[string]interface{})
		credential["type"] = "client-id-secret"
		credential["id"] = d.Get("credential_id").(int)
		cloud["credential"] = credential
	} else {
		credential := make(map[string]interface{})
		credential["type"] = "local"
		cloud["credential"] = credential
		config["clientId"] = d.Get("azure_client_id").(string)
		config["clientSecret"] = d.Get("azure_client_secret").(string)
	}

	if d.Get("region").(string) == "all" {
		cloud["regionCode"] = ""
	} else {
		cloud["regionCode"] = d.Get("region").(string)
	}

	// Resource Group
	if d.Get("resource_group").(string) == "all" {
		config["resourceGroup"] = ""
	} else {
		config["resourceGroup"] = d.Get("resource_group").(string)
	}

	// Inventory Existing Instances
	if d.Get("import_existing_instances").(bool) {
		config["importExisting"] = "on"
	} else {
		config["importExisting"] = ""
	}

	// RPC Mode
	config["rpcMode"] = d.Get("rpc_mode")

	config["applianceUrl"] = d.Get("appliance_url")
	cloud["timezone"] = d.Get("time_zone").(string)
	config["datacenterName"] = d.Get("datacenter_id")
	config["configManagementId"] = d.Get("config_management_integration_id").(string)
	cloud["guidanceMode"] = d.Get("guidance").(string)
	cloud["costingMode"] = d.Get("costing").(string)
	cloud["agentMode"] = d.Get("agent_install_mode").(string)

	cloud["config"] = config

	cloudType := make(map[string]interface{})
	cloudType["code"] = "azure"
	cloud["zoneType"] = cloudType

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
		Delay:        1 * time.Minute,
		PollInterval: 1 * time.Minute,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error creating cloud: %s", err)
	}

	// Successfully created resource, now set id
	d.SetId(int64ToString(cloudOutput.ID))
	resourceAzureCloudRead(ctx, d, meta)
	return diags
}

func resourceAzureCloudRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		d.Set("tenant_id", strconv.Itoa(int(cloud.AccountID)))
		d.Set("enabled", cloud.Enabled)
		d.Set("automatically_power_on_vms", cloud.AutoRecoverPowerState)
		d.Set("cloud_type", cloud.Config.CloudType)
		d.Set("azure_subscription_id", cloud.Config.SubscriberID)
		d.Set("azure_tenant_id", cloud.Config.TenantID)
		d.Set("credential_id", cloud.Credential.ID)
		if cloud.Credential.ID == 0 {
			d.Set("azure_client_id", cloud.Config.ClientID)
			d.Set("azure_client_secret", cloud.Config.ClientSecretHash)
		} else {
			d.Set("credential_id", cloud.Credential.ID)
		}
		if cloud.RegionCode == "" {
			d.Set("region", "all")
		} else {
			d.Set("region", cloud.RegionCode)
		}
		if cloud.Config.ResourceGroup == "" {
			d.Set("resource_group", "all")
		} else {
			d.Set("resource_group", cloud.Config.ResourceGroup)
		}
		if cloud.Config.ImportExisting == "on" {
			d.Set("import_existing_instances", true)
		} else {
			d.Set("import_existing_instances", false)
		}
		d.Set("rpc_mode", cloud.Config.RPCMode)
		d.Set("appliance_url", cloud.Config.ApplianceUrl)
		d.Set("time_zone", cloud.TimeZone)
		d.Set("datacenter_id", cloud.Config.DatacenterName)
		d.Set("config_management_integration_id", cloud.Config.ConfigManagementID)
		d.Set("guidance", cloud.GuidanceMode)
		d.Set("costing", cloud.CostingMode)
		d.Set("agent_install_mode", cloud.AgentMode)
		return diags
	}
}

func resourceAzureCloudUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	cloud := make(map[string]interface{})
	cloud["name"] = d.Get("name").(string)
	cloud["code"] = d.Get("code").(string)
	cloud["location"] = d.Get("location").(string)
	cloud["visibility"] = d.Get("visibility").(string)

	account := make(map[string]interface{})
	account["id"] = d.Get("tenant_id").(string)
	cloud["account"] = account
	cloud["accountId"] = d.Get("tenant_id").(string)

	cloud["enabled"] = d.Get("enabled").(bool)
	cloud["autoRecoverPowerState"] = d.Get("automatically_power_on_vms").(bool)

	config := make(map[string]interface{})

	config["cloudType"] = d.Get("cloud_type").(string)
	config["subscriberId"] = d.Get("azure_subscription_id").(string)
	config["tenantId"] = d.Get("azure_tenant_id").(string)

	if d.Get("credential_id").(int) != 0 {
		credential := make(map[string]interface{})
		credential["type"] = "client-id-secret"
		credential["id"] = d.Get("credential_id").(int)
		cloud["credential"] = credential
	} else {
		credential := make(map[string]interface{})
		credential["type"] = "local"
		cloud["credential"] = credential
		config["clientId"] = d.Get("azure_client_id").(string)
		config["clientSecret"] = d.Get("azure_client_secret").(string)
	}

	if d.Get("region").(string) == "all" {
		cloud["regionCode"] = ""
	} else {
		cloud["regionCode"] = d.Get("region").(string)
	}

	// Resource Group
	if d.Get("resource_group").(string) == "all" {
		config["resourceGroup"] = ""
	} else {
		config["resourceGroup"] = d.Get("resource_group").(string)
	}

	// Inventory Existing Instances
	if d.Get("import_existing_instances").(bool) {
		config["importExisting"] = "on"
	} else {
		config["importExisting"] = ""
	}

	// RPC Mode
	config["rpcMode"] = d.Get("rpc_mode")

	config["applianceUrl"] = d.Get("appliance_url")
	cloud["timezone"] = d.Get("time_zone").(string)
	config["datacenterName"] = d.Get("datacenter_id")
	if d.HasChange("config_management_integration_id") {
		config["configManagementId"] = d.Get("config_management_integration_id").(string)
	}
	cloud["guidanceMode"] = d.Get("guidance").(string)
	cloud["costingMode"] = d.Get("costing").(string)
	cloud["agentMode"] = d.Get("agent_install_mode").(string)

	cloud["config"] = config

	cloudType := make(map[string]interface{})
	cloudType["code"] = "azure"
	cloud["zoneType"] = cloudType

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
	return resourceAzureCloudRead(ctx, d, meta)
}

func resourceAzureCloudDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
