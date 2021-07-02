resource "morpheus_python_script_task" "tfexample_python_local" {
  name                = "tfexample_python_local"
  code                = "tfexample_python_local"
  source_type         = "local"
  script_content      = <<EOF
print('morpheus')
print('python')
EOF
  command_arguments   = "example"
  additional_packages = "pyyaml"
  python_binary       = "/usr/bin/python3"
  retryable           = true
  retry_count         = 1
  retry_delay_seconds = 10
  allow_custom_config = true
}