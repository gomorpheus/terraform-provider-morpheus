resource "morpheus_ruby_script_task" "tfexample_ruby_git" {
  name                = "tfexample_ruby_git"
  code                = "tfexample_ruby_git"
  labels              = ["demo", "terraform"]
  source_type         = "repository"
  result_type         = "json"
  script_path         = "example.rb"
  version_ref         = "master"
  repository_id       = 1
  retryable           = true
  retry_count         = 1
  retry_delay_seconds = 10
  allow_custom_config = true
}