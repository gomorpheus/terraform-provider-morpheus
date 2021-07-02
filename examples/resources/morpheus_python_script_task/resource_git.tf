resource "morpheus_python_script_task" "tfexample_python_git" {
  name                = "tfexample_python_git"
  code                = "tfexample_python_git"
  source_type         = "repository"
  result_type         = "json"
  script_path         = "example.py"
  version_ref         = "master"
  repository_id       = 1
  command_arguments   = "example"
  additional_packages = "pyyaml"
  python_binary       = "/usr/bin/python3"
  retryable           = true
  retry_count         = 1
  retry_delay_seconds = 10
  allow_custom_config = true
}