resource "morpheus_monitoring_setting" "tf_example_guidance_setting" {
  morpheus_auto_create_checks         = true
  morpheus_availability_time_frame    = ""
  morpheus_availability_precision     = ""
  morpheus_default_check_interval     = ""
  servicenow_monitoring_enabled       = false
  servicenow_integration_id           = 1
  servicenow_new_incident_action      = ""
  servicenow_close_incident_action    = ""
  servicenow_severity_info_impact     = ""
  servicenow_severity_warning_impact  = ""
  servicenow_severity_critical_impact = ""
  new_relic_monitoring_enabled        = false
  new_relic_license_key               = ""
}