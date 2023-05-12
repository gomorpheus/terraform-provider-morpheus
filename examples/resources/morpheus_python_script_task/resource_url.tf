resource "morpheus_python_script_task" "tfexample_python_url" {
  name                = "tfexample_python_url"
  code                = "tfexample_python_url"
  labels              = ["demo", "terraform"]
  source_type         = "url"
  result_type         = "json"
  script_path         = "https://example.com/example.py"
  command_arguments   = "example"
  additional_packages = "pyyaml"
  python_binary       = "/usr/bin/python3"
  retryable           = true
  retry_count         = 1
  retry_delay_seconds = 10
  allow_custom_config = true
}