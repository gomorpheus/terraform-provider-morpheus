resource "morpheus_restart_task" "tfexample_restart" {
  name                = "tfexample_restart"
  code                = "tfexample_restart"
  labels              = ["demo", "terraform"]
  retryable           = true
  retry_count         = 1
  retry_delay_seconds = 10
  allow_custom_config = true
}