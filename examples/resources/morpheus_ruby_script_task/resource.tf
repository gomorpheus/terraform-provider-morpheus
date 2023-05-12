resource "morpheus_ruby_script_task" "tfexample_ruby_local" {
  name                = "tfexample_ruby_local"
  code                = "tfexample_ruby_local"
  labels              = ["demo", "terraform"]
  source_type         = "local"
  script_content      = <<EOF
  puts "testing"
EOF
  retryable           = true
  retry_count         = 1
  retry_delay_seconds = 10
  allow_custom_config = true
}