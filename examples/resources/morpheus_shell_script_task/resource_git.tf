resource "morpheus_shell_script_task" "tfexample_shell_git" {
  name                = "tfexample_shell_git"
  code                = "tfexample_shell_git"
  source_type         = "repository"
  result_type         = "json"
  script_path         = "example.sh"
  version_ref         = "master"
  repository_id       = 1
  sudo                = true
  retryable           = true
  retry_count         = 1
  retry_delay_seconds = 10
  allow_custom_config = true
}