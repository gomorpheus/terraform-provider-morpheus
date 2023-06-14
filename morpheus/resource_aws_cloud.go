package morpheus

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
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

func resourceAWSCloud() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus AWS cloud resource.",
		CreateContext: resourceAWSCloudCreate,
		ReadContext:   resourceAWSCloudRead,
		UpdateContext: resourceAWSCloudUpdate,
		DeleteContext: resourceAWSCloudDelete,
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
			"account_number": {
				Description: "The AWS account number associated with the cloud integration",
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
				ValidateFunc: validation.StringInSlice([]string{"private", "public", ""}, false),
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
			"region": {
				Type:        schema.TypeString,
				Description: "The AWS region associated with the cloud integration",
				Required:    true,
			},
			"credential_id": {
				Description: "The ID of the credential store entry used for authentication",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"access_key": {
				Type:        schema.TypeString,
				Description: "The AWS access key used for authentication",
				Optional:    true,
				Computed:    true,
			},
			"secret_key": {
				Type:        schema.TypeString,
				Description: "The AWS secret key used for authentication",
				Optional:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256_hash := hex.EncodeToString(h.Sum(nil))
					return strings.EqualFold(old, sha256_hash)
				},
				RequiredWith: []string{"access_key"},
			},
			"role_arn": {
				Type:        schema.TypeString,
				Description: "The AWS IAM role ARN to assume for authentication",
				Optional:    true,
				Computed:    true,
			},
			"use_host_iam_credentials": {
				Description:   "Whether to use the IAM profile associated with the Morpheus server or not",
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"access_key", "secret_key"},
			},
			"inventory": {
				Type:         schema.TypeString,
				Description:  "Whether to import existing virtual machines (off, basic, full)",
				ValidateFunc: validation.StringInSlice([]string{"off", "basic", "full", ""}, false),
				Optional:     true,
				Computed:     true,
			},
			"vpc": {
				Type:        schema.TypeString,
				Description: "The VPC ID for a specific VPC (all or the AWS VPC id (vpc-25e6dae))",
				Optional:    true,
				Default:     "all",
			},
			"ebs_encryption": {
				Description: "Determines whether to configure default EBS volume encryption or not",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
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
			"guidance": {
				Type:         schema.TypeString,
				Description:  "Whether to enable guidance recommendations on the cloud (manual, off)",
				ValidateFunc: validation.StringInSlice([]string{"manual", "off"}, false),
				Optional:     true,
				Computed:     true,
			},
			"costing": {
				Type:         schema.TypeString,
				Description:  "Whether to enable costing on the cloud (off, costing, full)",
				ValidateFunc: validation.StringInSlice([]string{"off", "costing", "full"}, false),
				Optional:     true,
				Computed:     true,
			},
			"agent_install_mode": {
				Type:         schema.TypeString,
				Description:  "The method used to install the Morpheus agent on virtual machines provisioned in the cloud (ssh, cloudInit)",
				ValidateFunc: validation.StringInSlice([]string{"ssh", "cloudInit", ""}, false),
				Optional:     true,
				Computed:     true,
			},
		},
	}
}

func resourceAWSCloudCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	region := fmt.Sprintf("ec2.%s.amazonaws.com", d.Get("region").(string))
	config["endpoint"] = region

	if d.Get("credential_id").(int) != 0 {
		credential := make(map[string]interface{})
		credential["type"] = "access-key-secret"
		credential["id"] = d.Get("credential_id").(int)
		cloud["credential"] = credential
	} else {
		credential := make(map[string]interface{})
		credential["type"] = "local"
		cloud["credential"] = credential
		config["accessKey"] = d.Get("access_key").(string)
		config["secretKey"] = d.Get("secret_key").(string)
	}

	if d.Get("use_host_iam_credentials").(bool) {
		config["useHostCredentials"] = "on"
	} else {
		config["useHostCredentials"] = "off"
	}
	config["stsAssumeRole"] = d.Get("role_arn").(string)

	cloud["inventoryLevel"] = d.Get("inventory").(string)
	config["isVpc"] = "true"

	if d.Get("vpc").(string) == "all" {
		config["vpc"] = ""
	} else {
		config["vpc"] = d.Get("vpc").(string)
	}

	config["ebsEncryption"] = d.Get("ebs_encryption").(bool)
	if d.Get("ebs_encryption").(bool) {
		config["ebsEncryption"] = "on"
	} else {
		config["ebsEncryption"] = "off"
	}

	config["applianceUrl"] = d.Get("appliance_url")
	cloud["timezone"] = d.Get("time_zone").(string)
	config["datacenterName"] = d.Get("datacenter_id")
	cloud["guidanceMode"] = d.Get("guidance").(string)
	cloud["costingMode"] = d.Get("costing").(string)
	cloud["agentMode"] = d.Get("agent_install_mode").(string)

	//	config["useHostCredentials"] = d.Get("use_host_iam_credentials").(bool)

	cloud["config"] = config

	cloudType := make(map[string]interface{})
	cloudType["code"] = "amazon"
	cloud["zoneType"] = cloudType

	payload := map[string]interface{}{
		"zone": cloud,
	}

	req := &morpheus.Request{Body: payload}

	jsonRequest, _ := json.Marshal(req.Body)
	log.Printf("API JSON REQUEST: %s", string(jsonRequest))

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
	resourceAWSCloudRead(ctx, d, meta)
	return diags
}

func resourceAWSCloudRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		region := strings.Split(cloud.RegionCode, ".")
		d.Set("region", region[1])
		d.Set("credential_id", cloud.Credential.ID)
		d.Set("access_key", cloud.Config.AccessKey)
		d.Set("secret_key", cloud.Config.SecretKeyHash)
		if cloud.Config.UseHostCredentials == "" {
			d.Set("use_host_iam_credentials", false)
		} else {
			d.Set("use_host_iam_credentials", true)
		}
		d.Set("role_arn", cloud.Config.StsAssumeRole)
		d.Set("inventory", cloud.InventoryLevel)
		if cloud.Config.VPC == "" {
			d.Set("vpc", "all")
		} else {
			d.Set("vpc", cloud.Config.VPC)
		}
		if cloud.Config.EbsEncryption == "on" {
			d.Set("ebs_encryption", true)
		} else {
			d.Set("ebs_encryption", false)
		}
		d.Set("appliance_url", cloud.Config.ApplianceUrl)
		d.Set("time_zone", cloud.TimeZone)
		d.Set("datacenter_id", cloud.Config.DatacenterName)
		d.Set("guidance", cloud.GuidanceMode)
		d.Set("costing", cloud.CostingMode)
		d.Set("agent_install_mode", cloud.AgentMode)
		d.Set("account_number", cloud.ExternalID)
		return diags
	}
}

func resourceAWSCloudUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	region := fmt.Sprintf("ec2.%s.amazonaws.com", d.Get("region").(string))
	config["endpoint"] = region

	if d.Get("credential_id").(int) != 0 {
		credential := make(map[string]interface{})
		credential["type"] = "access-key-secret"
		credential["id"] = d.Get("credential_id").(int)
		cloud["credential"] = credential
	} else {
		credential := make(map[string]interface{})
		credential["type"] = "local"
		cloud["credential"] = credential
		config["accessKey"] = d.Get("access_key").(string)
		config["secretKey"] = d.Get("secret_key").(string)
	}

	if d.Get("use_host_iam_credentials").(bool) {
		config["useHostCredentials"] = "on"
	} else {
		config["useHostCredentials"] = "off"
	}
	config["stsAssumeRole"] = d.Get("role_arn").(string)

	cloud["inventoryLevel"] = d.Get("inventory").(string)
	config["isVpc"] = "true"

	if d.Get("vpc").(string) == "all" {
		config["vpc"] = ""
	} else {
		config["vpc"] = d.Get("vpc").(string)
	}

	config["ebsEncryption"] = d.Get("ebs_encryption").(bool)
	if d.Get("ebs_encryption").(bool) {
		config["ebsEncryption"] = "on"
	} else {
		config["ebsEncryption"] = "off"
	}

	config["applianceUrl"] = d.Get("appliance_url")
	cloud["timezone"] = d.Get("time_zone").(string)
	config["datacenterName"] = d.Get("datacenter_id")
	cloud["guidanceMode"] = d.Get("guidance").(string)
	cloud["costingMode"] = d.Get("costing").(string)
	cloud["agentMode"] = d.Get("agent_install_mode").(string)

	//	config["useHostCredentials"] = d.Get("use_host_iam_credentials").(bool)

	cloud["config"] = config

	cloudType := make(map[string]interface{})
	cloudType["code"] = "amazon"
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
	return resourceAWSCloudRead(ctx, d, meta)
}

func resourceAWSCloudDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
