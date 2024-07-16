package morpheus

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceApplianceSetting() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus appliance setting resource.",
		CreateContext: resourceApplianceSettingCreate,
		ReadContext:   resourceApplianceSettingRead,
		UpdateContext: resourceApplianceSettingUpdate,
		DeleteContext: resourceApplianceSettingDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the appliance settings",
				Computed:    true,
			},
			"appliance_url": {
				Type:        schema.TypeString,
				Description: "The default URL used for Agent install and Agent functionality",
				Optional:    true,
				Computed:    true,
			},
			"internal_appliance_url": {
				Type:        schema.TypeString,
				Description: "An override for the default appliance URL when PXE boot is utilized",
				Optional:    true,
				Computed:    true,
			},
			"api_allowed_origins": {
				Type:        schema.TypeString,
				Description: "A CORS-related field which specifies the origins that are allowed to access the Morpheus API",
				Optional:    true,
				Computed:    true,
			},
			/* AWAITING API SUPPORT
			"cloud_sync_interval": {
				Type:        schema.TypeInt,
				Description: "The interval at which cloud integrations are synced",
				Optional:    true,
				Computed:    true,
			},
			"cluster_sync_interval": {
				Type:        schema.TypeInt,
				Description: "The interval at which clusters are synced",
				Optional:    true,
				Computed:    true,
			},
			"usage_retainment_period": {
				Type:        schema.TypeInt,
				Description: "The number of days that usage data is stored",
				Optional:    true,
				Computed:    true,
			},
			"incident_retainment_period": {
				Type:        schema.TypeInt,
				Description: "The number of days that incident data is stored",
				Optional:    true,
				Computed:    true,
			},
			"denied_hosts": {
				Type:        schema.TypeString,
				Description: "A comma delimited list of ips/hostnames to be blocked when using HTTP Task Types or REST Datasource Option Lists",
				Optional:    true,
				Computed:    true,
			},
			"approved_hosts": {
				Type:        schema.TypeString,
				Description: "a comma delimited list of ips/hostnames to be allowed when using HTTP Task Types or REST Datasource Option Lists",
				Optional:    true,
				Computed:    true,
			},
			*/
			"stats_retainment_period": {
				Type:        schema.TypeInt,
				Description: "The number of days that incident data is stored",
				Optional:    true,
				Computed:    true,
			},
			"registration_enabled": {
				Type:         schema.TypeBool,
				Description:  "Whether tenant registration is enabled",
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				RequiredWith: []string{"default_role_id", "default_user_role_id"},
			},
			"default_role_id": {
				Type:         schema.TypeString,
				Description:  "The ID of the default role",
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"registration_enabled", "default_user_role_id"},
			},
			"default_user_role_id": {
				Type:         schema.TypeString,
				Description:  "The ID of the default user role",
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"registration_enabled", "default_role_id"},
			},
			"docker_privileged_mode": {
				Type:        schema.TypeBool,
				Description: "Whether Docker privileged mode is enabled",
				Optional:    true,
				Computed:    true,
			},
			/* Awaiting API Support
			"minimum_password_length": {
				Type:        schema.TypeBool,
				Description: "The ID of the default user role",
				Optional:    true,
				Computed:    true,
			},
			"minimum_password_uppercase": {
				Type:        schema.TypeBool,
				Description: "The ID of the default user role",
				Optional:    true,
				Computed:    true,
			},
			"minimum_password_numbers": {
				Type:        schema.TypeBool,
				Description: "The ID of the default user role",
				Optional:    true,
				Computed:    true,
			},
			"minimum_password_symbols": {
				Type:        schema.TypeBool,
				Description: "The ID of the default user role",
				Optional:    true,
				Computed:    true,
			},
			*/
			"smtp_from_address": {
				Type:        schema.TypeString,
				Description: "The email address to send system emails from",
				Optional:    true,
				Computed:    true,
			},
			"smtp_server": {
				Type:        schema.TypeString,
				Description: "The hostname or IP address of the SMTP server",
				Optional:    true,
				Computed:    true,
			},
			"smtp_port": {
				Type:        schema.TypeString,
				Description: "The SMTP server port",
				Optional:    true,
				Computed:    true,
			},
			"smtp_use_ssl": {
				Type:        schema.TypeBool,
				Description: "Whether to use SSL or not when connecting to the SMTP server",
				Optional:    true,
				Computed:    true,
			},
			"smtp_use_tls": {
				Type:        schema.TypeBool,
				Description: "Whether to use TLS or not when connecting to the SMTP server",
				Optional:    true,
				Computed:    true,
			},
			"smtp_username": {
				Type:        schema.TypeString,
				Description: "The username for the user account used to authenticate to the SMTP server",
				Optional:    true,
				Computed:    true,
			},
			"smtp_password": {
				Type:        schema.TypeString,
				Description: "The password for the user account used to authenticate to the SMTP server",
				Optional:    true,
				Sensitive:   true,
				Computed:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256_hash := hex.EncodeToString(h.Sum(nil))
					return strings.EqualFold(old, sha256_hash)
				},
			},
			"proxy_host": {
				Type:        schema.TypeString,
				Description: "The hostname or IP address of the proxy host",
				Optional:    true,
				Computed:    true,
			},
			"proxy_port": {
				Type:        schema.TypeString,
				Description: "The proxy host port",
				Optional:    true,
				Computed:    true,
			},
			"proxy_user": {
				Type:        schema.TypeString,
				Description: "The username for authenticating to the proxy",
				Optional:    true,
				Computed:    true,
			},
			"proxy_password": {
				Type:        schema.TypeString,
				Description: "The password for authenticating to the proxy",
				Optional:    true,
				Sensitive:   true,
				Computed:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256_hash := hex.EncodeToString(h.Sum(nil))
					return strings.EqualFold(old, sha256_hash)
				},
			},
			"proxy_domain": {
				Type:        schema.TypeString,
				Description: "The proxy domain name",
				Optional:    true,
				Computed:    true,
			},
			"proxy_workstation": {
				Type:        schema.TypeString,
				Description: "The proxy workstation",
				Optional:    true,
				Computed:    true,
			},
			"currency_provider": {
				Type:         schema.TypeString,
				Description:  "The currency provider (openexchange, fixer)",
				ValidateFunc: validation.StringInSlice([]string{"openexchange", "fixer"}, true),
				Optional:     true,
				Computed:     true,
			},
			"currency_provider_api_key": {
				Type:        schema.TypeString,
				Description: "The API key for the currency provider",
				Optional:    true,
				Computed:    true,
			},
			/*
				"enabled_cloud_ids": {
					Type:        schema.TypeList,
					Description: "The ids of the cloud types to enable",
					Optional:    true,
					Computed:    true,
					Elem:        &schema.Schema{Type: schema.TypeInt},
				},
			*/
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceApplianceSettingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	applianceSettings := make(map[string]interface{})

	applianceUrl, applianceUrlOk := d.GetOk("appliance_url")
	if applianceUrlOk {
		applianceSettings["applianceUrl"] = applianceUrl
	}

	internalApplianceUrl, internalApplianceUrlOk := d.GetOk("internal_appliance_url")
	if internalApplianceUrlOk {
		applianceSettings["internalApplianceUrl"] = internalApplianceUrl
	}

	apiAllowedOrigins, apiAllowedOriginsOk := d.GetOk("api_allowed_origins")
	if apiAllowedOriginsOk {
		applianceSettings["corsAllowed"] = apiAllowedOrigins
	}

	applianceSettings["registrationEnabled"] = d.Get("registration_enabled").(bool)

	if d.Get("registration_enabled").(bool) {
		defaultRoleID, defaultRoleIDOk := d.GetOk("default_role_id")
		if defaultRoleIDOk {
			applianceSettings["defaultRoleId"] = defaultRoleID
		}

		defaultUserRoleID, defaultUserRoleIDOk := d.GetOk("default_user_role_id")
		if defaultUserRoleIDOk {
			applianceSettings["defaultUserRoleId"] = defaultUserRoleID
		}
	}

	applianceSettings["dockerPrivilegedMode"] = d.Get("docker_privileged_mode").(bool)

	smtpFromAddress, smtpFromAddressOk := d.GetOk("smtp_from_address")
	if smtpFromAddressOk {
		applianceSettings["smtpMailFrom"] = smtpFromAddress
	}

	smtpServer, smtpServerOk := d.GetOk("smtp_server")
	if smtpServerOk {
		applianceSettings["smtpServer"] = smtpServer
	}

	smtpPort, smtpPortOk := d.GetOk("smtp_port")
	if smtpPortOk {
		applianceSettings["smtpPort"] = smtpPort
	}

	smtpSSL, smtpSSLOk := d.GetOk("smtp_use_ssl")
	if smtpSSLOk {
		applianceSettings["smtpSSL"] = smtpSSL
	}

	smtpTLS, smtpTLSOk := d.GetOk("smtp_use_tls")
	if smtpTLSOk {
		applianceSettings["smtpTLS"] = smtpTLS
	}

	smtpUsername, smtpUsernameOk := d.GetOk("smtp_username")
	if smtpUsernameOk {
		applianceSettings["smtpUser"] = smtpUsername
	}

	smtpPassword, smtpPasswordOk := d.GetOk("smtp_password")
	if smtpPasswordOk {
		applianceSettings["smtpPassword"] = smtpPassword
	}

	proxyHost, proxyHostOk := d.GetOk("proxy_host")
	if proxyHostOk {
		applianceSettings["proxyHost"] = proxyHost
	}

	proxyPort, proxyPortOk := d.GetOk("proxy_port")
	if proxyPortOk {
		applianceSettings["proxyPort"] = proxyPort
	}

	proxyUser, proxyUserOk := d.GetOk("proxy_user")
	if proxyUserOk {
		applianceSettings["proxyUser"] = proxyUser
	}

	proxyPassword, proxyPasswordOk := d.GetOk("proxy_password")
	if proxyPasswordOk {
		applianceSettings["proxyPassword"] = proxyPassword
	}

	proxyDomain, proxyDomainOk := d.GetOk("proxy_domain")
	if proxyDomainOk {
		applianceSettings["proxyDomain"] = proxyDomain
	}

	proxyWorkstation, proxyWorkstationOk := d.GetOk("proxy_workstation")
	if proxyWorkstationOk {
		applianceSettings["proxyWorkstation"] = proxyWorkstation
	}

	currencyProvider, currencyProviderOk := d.GetOk("currency_provider")
	if currencyProviderOk {
		applianceSettings["currencyProvider"] = currencyProvider
	}

	currencyKey, currencyKeyOk := d.GetOk("currency_provider_api_key")
	if currencyKeyOk {
		applianceSettings["currencyKey"] = currencyKey
	}

	//applianceSettings["enableZoneTypes"] = d.Get("enabled_cloud_ids")
	//applianceSettings["disableAllZoneTypes"] = true

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"applianceSettings": applianceSettings,
		},
	}

	resp, err := client.UpdateApplianceSettings(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.UpdateApplianceSettingsResult)
	_ = result.ApplianceSettings
	// Successfully created resource, now set id
	d.SetId(int64ToString(1))

	resourceApplianceSettingRead(ctx, d, meta)
	return diags
}

func resourceApplianceSettingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error

	resp, err = client.GetApplianceSettings(&morpheus.Request{})

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
	result := resp.Result.(*morpheus.GetApplianceSettingsResult)
	applianceSetting := result.ApplianceSettings
	d.SetId(int64ToString(1))
	d.Set("appliance_url", applianceSetting.ApplianceURL)
	d.Set("internal_appliance_url", applianceSetting.InternalApplianceURL)
	d.Set("api_allowed_origins", applianceSetting.CorsAllowed)
	d.Set("registration_enabled", applianceSetting.RegistrationEnabled)
	d.Set("default_role_id", applianceSetting.DefaultRoleID)
	d.Set("default_user_role_id", applianceSetting.DefaultUserRoleID)
	d.Set("docker_privileged_mode", applianceSetting.DockerPrivilegedMode)
	d.Set("smtp_from_address", applianceSetting.SMTPMailFrom)
	d.Set("smtp_server", applianceSetting.SMTPServer)
	d.Set("smtp_port", applianceSetting.SMTPPort)
	d.Set("smtp_use_ssl", applianceSetting.SMTPSSL)
	d.Set("smtp_use_tls", applianceSetting.SMTPTLS)
	d.Set("smtp_username", applianceSetting.SMTPUser)
	d.Set("smtp_password", applianceSetting.SMTPPasswordHash)
	d.Set("proxy_host", applianceSetting.ProxyHost)
	d.Set("proxy_port", applianceSetting.ProxyPort)
	d.Set("proxy_user", applianceSetting.ProxyUser)
	d.Set("proxy_password", applianceSetting.ProxyPasswordHash)
	d.Set("proxy_domain", applianceSetting.ProxyDomain)
	d.Set("proxy_workstation", applianceSetting.ProxyWorkstation)
	d.Set("currency_provider", applianceSetting.CurrencyProvider)
	d.Set("currency_provider_api_key", applianceSetting.CurrencyKey)
	//var enabledClouds []int
	//for _, cloud := range applianceSetting.EnabledZoneTypes {
	//	enabledClouds = append(enabledClouds, int(cloud.ID))
	//}

	//cloudIdPayload := matchCloudIdsWithSchema(enabledClouds, d.Get("enabled_cloud_ids").([]interface{}))
	//d.Set("enabled_cloud_ids", cloudIdPayload)
	return diags
}

func resourceApplianceSettingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	applianceSettings := make(map[string]interface{})

	if d.HasChange("appliance_url") {
		applianceSettings["applianceUrl"] = d.Get("appliance_url")
	}

	if d.HasChange("internal_appliance_url") {
		applianceSettings["internalApplianceUrl"] = d.Get("internal_appliance_url")
	}

	if d.HasChange("docker_privileged_mode") {
		applianceSettings["dockerPrivilegedMode"] = d.Get("docker_privileged_mode")
	}

	if d.HasChange("smtp_from_address") {
		applianceSettings["smtpMailFrom"] = d.Get("smtp_from_address")
	}

	if d.HasChange("smtp_server") {
		applianceSettings["smtpServer"] = d.Get("smtp_server")
	}

	if d.HasChange("smtp_port") {
		applianceSettings["smtpPort"] = d.Get("smtp_port")
	}

	if d.HasChange("smtp_use_ssl") {
		applianceSettings["smtpSSL"] = d.Get("smtp_use_ssl")
	}

	if d.HasChange("smtp_use_tls") {
		applianceSettings["smtpTLS"] = d.Get("smtp_use_tls")
	}

	if d.HasChange("smtp_username") {
		applianceSettings["smtpUser"] = d.Get("smtp_username")
	}

	if d.HasChange("smtp_password") {
		applianceSettings["smtpPassword"] = d.Get("smtp_password")
	}

	if d.HasChange("proxy_host") {
		applianceSettings["proxyHost"] = d.Get("proxy_host")
	}

	if d.HasChange("proxy_port") {
		applianceSettings["proxyPort"] = d.Get("proxy_port")
	}

	if d.HasChange("proxy_user") {
		applianceSettings["proxyUser"] = d.Get("proxy_user")
	}

	if d.HasChange("proxy_password") {
		applianceSettings["proxyPassword"] = d.Get("proxy_password")
	}

	if d.HasChange("proxy_domain") {
		applianceSettings["proxyDomain"] = d.Get("proxy_domain")
	}

	if d.HasChange("proxy_workstation") {
		applianceSettings["proxyWorkstation"] = d.Get("proxy_workstation")
	}

	if d.HasChange("currency_provider") {
		applianceSettings["currencyProvider"] = d.Get("currency_provider")
	}

	if d.HasChange("currency_provider_api_key") {
		applianceSettings["currencyKey"] = d.Get("currency_provider_api_key")
	}

	//if d.HasChange("enableZoneTypes") {
	//	applianceSettings["enableZoneTypes"] = d.Get("enabled_cloud_ids")
	//}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"applianceSettings": applianceSettings,
		},
	}

	resp, err := client.UpdateApplianceSettings(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateApplianceSettingsResult)
	_ = result.ApplianceSettings
	// Successfully created resource, now set id
	d.SetId(int64ToString(1))

	return resourceApplianceSettingRead(ctx, d, meta)
}

func resourceApplianceSettingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	d.SetId("")
	return diags
}