resource "morpheus_groovy_script_task" "tfexample_groovy_git" {
  name                = "tfexample_groovy_git"
  code                = "tfexample_groovy_git"
  source_type         = "repository"
  result_type         = "json"
  script_path         = "example.groovy"
  version_ref         = "master"
  repository_id       = 1
  retryable           = true
  retry_count         = 1
  retry_delay_seconds = 10
  allow_custom_config = true
}