resource "morpheus_powershell_script_task" "tfexample_powershell_local" {
  name                = "tfexample_powershell_local"
  code                = "tfexample_powershell_local"
  source_type         = "local"
  script_content      = <<EOF
  Write-Output "testing"
EOF
  elevated_shell      = true
  retryable           = true
  retry_count         = 1
  retry_delay_seconds = 10
  allow_custom_config = true
}