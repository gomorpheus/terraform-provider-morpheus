resource "morpheus_powershell_script_task" "tfexample_powershell_url" {
  name                = "tfexample_powershell_url"
  code                = "tfexample_powershell_url"
  source_type         = "url"
  result_type         = "json"
  script_path         = "https://example.com/example.ps"
  elevated_shell                = true
  retryable           = true
  retry_count         = 1
  retry_delay_seconds = 10
  allow_custom_config = true
}