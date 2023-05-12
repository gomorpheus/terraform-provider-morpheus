resource "morpheus_write_attributes_task" "tfexample_write_attributes" {
  name                = "tfexample_write_attributes"
  code                = "tfexample_write_attributes"
  labels              = ["demo", "terraform"]
  attributes          = <<EOF
{"demo":"test"}
EOF
  retryable           = true
  retry_count         = 1
  retry_delay_seconds = 10
  allow_custom_config = true
}