resource "morpheus_javascript_task" "tfexample_javascript" {
  name                = "tfexample_javascript"
  code                = "tfexample_javascript"
  script_content      = <<EOF
console.log("testing")
EOF
  retryable           = true
  retry_count         = 1
  retry_delay_seconds = 10
  allow_custom_config = true
}