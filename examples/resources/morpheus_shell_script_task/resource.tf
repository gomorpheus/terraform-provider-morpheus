resource "morpheus_shell_script_task" "tfexample_shell_local" {
  name                = "tfexample_shell_local"
  code                = "tfexample_shell_local"
  labels              = ["demo", "terraform"]
  source_type         = "local"
  script_content      = <<EOF
  echo "testing"
EOF
  sudo                = true
  retryable           = true
  retry_count         = 1
  retry_delay_seconds = 10
  allow_custom_config = true
}