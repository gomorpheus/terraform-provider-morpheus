package morpheus

import (
	"context"
	"encoding/json"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceMonitoringSetting() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus monitoring setting resource.",
		CreateContext: resourceMonitoringSettingCreate,
		ReadContext:   resourceMonitoringSettingRead,
		UpdateContext: resourceMonitoringSettingUpdate,
		DeleteContext: resourceMonitoringSettingDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the monitoring setting",
				Computed:    true,
			},
			"morpheus_auto_create_checks": {
				Type:        schema.TypeBool,
				Description: "When enabled a Monitoring Check will automatically be create for Instances and Apps",
				Optional:    true,
				Computed:    true,
			},
			"morpheus_availability_time_frame": {
				Type:        schema.TypeInt,
				Description: "The number of days availability should be calculated for",
				Optional:    true,
				Computed:    true,
			},
			"morpheus_availability_precision": {
				Type:         schema.TypeInt,
				Description:  "The number of decimal places availability should be displayed in, can be anywhere between 0 and 5",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1, 2, 3, 4, 5}),
			},
			"morpheus_default_check_interval": {
				Type:         schema.TypeInt,
				Description:  "The default interval in minutes to use when creating new checks (1, 2, 3, 4, 5, 10, 15, 20, 25, 30, 45, 60, 120, 180)",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 3, 4, 5, 10, 15, 20, 25, 30, 45, 60, 120, 180}),
			},
			"servicenow_monitoring_enabled": {
				Type:         schema.TypeBool,
				Description:  "Whether the ServiceNow monitoring integration is enabled",
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"servicenow_integration_id"},
			},
			"servicenow_integration_id": {
				Type:         schema.TypeInt,
				Description:  "The id of the ServiceNow monitoring integration",
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"servicenow_monitoring_enabled"},
			},
			"servicenow_new_incident_action": {
				Type:         schema.TypeString,
				Description:  "The Service Now action to take when a Morpheus incident is created (create, none)",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"create", "none"}, false),
				RequiredWith: []string{"servicenow_monitoring_enabled"},
			},
			"servicenow_close_incident_action": {
				Type:         schema.TypeString,
				Description:  "The Service Now action to take when a Morpheus incident is closed (activity, close, none)",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"activity", "close", "none"}, false),
				RequiredWith: []string{"servicenow_monitoring_enabled"},
			},
			"servicenow_severity_info_impact": {
				Type:         schema.TypeString,
				Description:  "The ServiceNow impact level to map to the Morpheus info severity (high, medium, low)",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"high", "medium", "low"}, false),
				RequiredWith: []string{"servicenow_monitoring_enabled"},
			},
			"servicenow_severity_warning_impact": {
				Type:         schema.TypeString,
				Description:  "The ServiceNow impact level to map to the Morpheus warning severity (high, medium, low)",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"high", "medium", "low"}, false),
				RequiredWith: []string{"servicenow_monitoring_enabled"},
			},
			"servicenow_severity_critical_impact": {
				Type:         schema.TypeString,
				Description:  "The ServiceNow impact level to map to the Morpheus critical severity (high, medium, low)",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"high", "medium", "low"}, false),
				RequiredWith: []string{"servicenow_monitoring_enabled"},
			},
			"new_relic_monitoring_enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the New Relic monitoring integration is enabled",
				Optional:    true,
				Computed:    true,
			},
			"new_relic_license_key": {
				Type:        schema.TypeString,
				Description: "The New Relic license key",
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceMonitoringSettingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	monitoringSettings := make(map[string]interface{})

	monitoringSettings["autoManageChecks"] = d.Get("morpheus_auto_create_checks")

	availabilityTimeFrame, availabilityTimeFrameok := d.GetOk("morpheus_availability_time_frame")
	if availabilityTimeFrameok {
		monitoringSettings["availabilityTimeFrame"] = availabilityTimeFrame
	}

	availabilityPrecision, availabilityPrecisionok := d.GetOk("morpheus_availability_precision")
	if availabilityPrecisionok {
		monitoringSettings["availabilityPrecision"] = availabilityPrecision
	}

	defaultCheckInterval, defaultCheckIntervalok := d.GetOk("morpheus_default_check_interval")
	if defaultCheckIntervalok {
		monitoringSettings["defaultCheckInterval"] = defaultCheckInterval
	}

	serviceNowSettings := make(map[string]interface{})

	serviceNowNewIntegrationId, serviceNowNewIntegrationIdok := d.GetOk("servicenow_integration_id")
	if serviceNowNewIntegrationIdok {
		serviceNowIntegration := make(map[string]interface{})
		serviceNowIntegration["id"] = serviceNowNewIntegrationId
		serviceNowSettings["integration"] = serviceNowIntegration
		serviceNowSettings["enabled"] = d.Get("servicenow_monitoring_enabled")
		monitoringSettings["serviceNow"] = serviceNowSettings
	}

	serviceNowNewIncidentAction, serviceNowNewIncidentActionok := d.GetOk("servicenow_new_incident_action")
	if serviceNowNewIncidentActionok {
		serviceNowSettings["newIncidentAction"] = serviceNowNewIncidentAction
		monitoringSettings["serviceNow"] = serviceNowSettings
	}

	serviceNowCloseIncidentAction, serviceNowCloseIncidentActionok := d.GetOk("servicenow_close_incident_action")
	if serviceNowCloseIncidentActionok {
		serviceNowSettings["closeIncidentAction"] = serviceNowCloseIncidentAction
		monitoringSettings["serviceNow"] = serviceNowSettings
	}

	serviceNowInfoImpact, serviceNowInfoImpactok := d.GetOk("servicenow_severity_info_impact")
	if serviceNowInfoImpactok {
		serviceNowSettings["infoMapping"] = serviceNowInfoImpact
		monitoringSettings["serviceNow"] = serviceNowSettings
	}

	serviceNowWarningImpact, serviceNowWarningImpactok := d.GetOk("servicenow_severity_warning_impact")
	if serviceNowWarningImpactok {
		serviceNowSettings["warningMapping"] = serviceNowWarningImpact
		monitoringSettings["serviceNow"] = serviceNowSettings
	}

	serviceNowCriticalImpact, serviceNowCriticalImpactok := d.GetOk("servicenow_severity_critical_impact")
	if serviceNowCriticalImpactok {
		serviceNowSettings["criticalMapping"] = serviceNowCriticalImpact
		monitoringSettings["serviceNow"] = serviceNowSettings
	}

	newRelicSettings := make(map[string]interface{})

	newRelicSettings["enabled"] = d.Get("new_relic_monitoring_enabled")
	monitoringSettings["newRelic"] = newRelicSettings

	newRelicLicenseKey, newRelicLicenseKeyok := d.GetOk("new_relic_license_key")
	if newRelicLicenseKeyok {
		newRelicSettings["licenseKey"] = newRelicLicenseKey
		monitoringSettings["newRelic"] = newRelicSettings
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"monitoringSettings": monitoringSettings,
		},
	}

	jsonRequest, _ := json.Marshal(req.Body)
	log.Printf("API JSON REQUEST: %s", string(jsonRequest))

	resp, err := client.UpdateMonitoringSettings(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.UpdateMonitoringSettingsResult)
	_ = result.MonitoringSettings
	// Successfully created resource, now set id
	d.SetId(int64ToString(1))

	resourceMonitoringSettingRead(ctx, d, meta)
	return diags
}

func resourceMonitoringSettingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error

	resp, err = client.GetMonitoringSettings(&morpheus.Request{})

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

	// store resource data
	result := resp.Result.(*morpheus.GetMonitoringSettingsResult)
	monitoringSetting := result.MonitoringSettings
	d.SetId(int64ToString(1))
	d.Set("morpheus_auto_create_checks", monitoringSetting.AutoManageChecks)
	d.Set("morpheus_availability_time_frame", monitoringSetting.AvailabilityTimeFrame)
	d.Set("morpheus_availability_precision", monitoringSetting.AvailabilityPrecision)
	d.Set("morpheus_default_check_interval", monitoringSetting.DefaultCheckInterval)
	d.Set("servicenow_monitoring_enabled", monitoringSetting.ServiceNow.Enabled)
	d.Set("servicenow_integration_id", monitoringSetting.ServiceNow.Integration.ID)
	d.Set("servicenow_new_incident_action", monitoringSetting.ServiceNow.NewIncidentAction)
	d.Set("servicenow_close_incident_action", monitoringSetting.ServiceNow.CloseIncidentAction)
	d.Set("servicenow_severity_info_impact", monitoringSetting.ServiceNow.InfoMapping)
	d.Set("servicenow_severity_warning_impact", monitoringSetting.ServiceNow.WarningMapping)
	d.Set("servicenow_severity_critical_impact", monitoringSetting.ServiceNow.CriticalMapping)
	d.Set("new_relic_monitoring_enabled", monitoringSetting.NewRelic.Enabled)
	d.Set("new_relic_license_key", monitoringSetting.NewRelic.LicenseKey)

	return diags
}

func resourceMonitoringSettingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	monitoringSettings := make(map[string]interface{})

	if d.HasChange("morpheus_auto_create_checks") {
		monitoringSettings["autoManageChecks"] = d.Get("morpheus_auto_create_checks")
	}

	if d.HasChange("morpheus_availability_time_frame") {
		monitoringSettings["availabilityTimeFrame"] = d.Get("morpheus_availability_time_frame")
	}

	if d.HasChange("morpheus_availability_precision") {
		monitoringSettings["availabilityPrecision"] = d.Get("morpheus_availability_precision")
	}

	if d.HasChange("morpheus_default_check_interval") {
		monitoringSettings["defaultCheckInterval"] = d.Get("morpheus_default_check_interval")
	}

	serviceNowSettings := make(map[string]interface{})

	if d.HasChange("servicenow_monitoring_enabled") {
		serviceNowSettings["enabled"] = d.Get("servicenow_monitoring_enabled")
		monitoringSettings["serviceNow"] = serviceNowSettings
	}

	if d.HasChange("servicenow_integration_id") {
		serviceNowIntegration := make(map[string]interface{})
		serviceNowIntegration["id"] = d.Get("servicenow_integration_id")
		serviceNowSettings["integration"] = serviceNowIntegration
		monitoringSettings["serviceNow"] = serviceNowSettings
	}

	if d.HasChange("servicenow_new_incident_action") {
		serviceNowSettings["newIncidentAction"] = d.Get("servicenow_new_incident_action")
		monitoringSettings["serviceNow"] = serviceNowSettings
	}

	if d.HasChange("servicenow_close_incident_action") {
		serviceNowSettings["closeIncidentAction"] = d.Get("servicenow_close_incident_action")
		monitoringSettings["serviceNow"] = serviceNowSettings
	}

	if d.HasChange("servicenow_severity_info_impact") {
		serviceNowSettings["infoMapping"] = d.Get("servicenow_severity_info_impact")
		monitoringSettings["serviceNow"] = serviceNowSettings
	}

	if d.HasChange("servicenow_severity_warning_impact") {
		serviceNowSettings["warningMapping"] = d.Get("servicenow_severity_warning_impact")
		monitoringSettings["serviceNow"] = serviceNowSettings
	}

	if d.HasChange("servicenow_severity_critical_impact") {
		serviceNowSettings["criticalMapping"] = d.Get("servicenow_severity_critical_impact")
		monitoringSettings["serviceNow"] = serviceNowSettings
	}

	newRelicSettings := make(map[string]interface{})

	if d.HasChange("new_relic_monitoring_enabled") {
		newRelicSettings["enabled"] = d.Get("new_relic_monitoring_enabled")
		monitoringSettings["newRelic"] = newRelicSettings
	}

	if d.HasChange("new_relic_license_key") {
		newRelicSettings["licenseKey"] = d.Get("new_relic_license_key")
		monitoringSettings["newRelic"] = newRelicSettings
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"monitoringSettings": monitoringSettings,
		},
	}

	jsonRequest, _ := json.Marshal(req.Body)
	log.Printf("API JSON REQUEST: %s", string(jsonRequest))

	resp, err := client.UpdateMonitoringSettings(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateMonitoringSettingsResult)
	_ = result.MonitoringSettings
	// Successfully created resource, now set id
	d.SetId(int64ToString(1))

	return resourceMonitoringSettingRead(ctx, d, meta)
}

func resourceMonitoringSettingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	d.SetId("")
	return diags
}
