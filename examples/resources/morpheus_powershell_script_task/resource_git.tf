resource "morpheus_powershell_script_task" "tfexample_powershell_git" {
  name                = "tfexample_powershell_git"
  code                = "tfexample_powershell_git"
  labels              = ["demo", "terraform"]
  source_type         = "repository"
  result_type         = "json"
  script_path         = "example.ps"
  version_ref         = "master"
  repository_id       = 1
  elevated_shell      = true
  retryable           = true
  retry_count         = 1
  retry_delay_seconds = 10
  allow_custom_config = true
}