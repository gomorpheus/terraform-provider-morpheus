resource "morpheus_monitoring_setting" "tf_example_guidance_setting" {
  morpheus_auto_create_checks         = true
  morpheus_availability_time_frame    = 30
  morpheus_availability_precision     = 4
  morpheus_default_check_interval     = 120
  servicenow_monitoring_enabled       = true
  servicenow_integration_id           = 1
  servicenow_new_incident_action      = "create"
  servicenow_close_incident_action    = "activity"
  servicenow_severity_info_impact     = "high"
  servicenow_severity_warning_impact  = "high"
  servicenow_severity_critical_impact = "low"
  new_relic_monitoring_enabled        = true
  new_relic_license_key               = "ABC123"
}