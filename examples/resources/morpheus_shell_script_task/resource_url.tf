resource "morpheus_shell_script_task" "tfexample_shell_url" {
  name                = "tfexample_shell_url"
  code                = "tfexample_shell_url"
  labels              = ["demo", "terraform"]
  source_type         = "url"
  result_type         = "json"
  script_path         = "https://example.com/example.sh"
  sudo                = true
  retryable           = true
  retry_count         = 1
  retry_delay_seconds = 10
  allow_custom_config = true
}