resource "morpheus_groovy_script_task" "tfexample_groovy_url" {
  name                = "tfexample_groovy_url"
  code                = "tfexample_groovy_url"
  labels              = ["demo", "terraform"]
  source_type         = "url"
  result_type         = "json"
  script_path         = "https://example.com/example.groovy"
  retryable           = true
  retry_count         = 1
  retry_delay_seconds = 10
  allow_custom_config = true
}