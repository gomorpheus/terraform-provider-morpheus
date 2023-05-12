resource "morpheus_ruby_script_task" "tfexample_ruby_url" {
  name                = "tfexample_ruby_url"
  code                = "tfexample_ruby_url"
  labels              = ["demo", "terraform"]
  source_type         = "url"
  result_type         = "json"
  script_path         = "https://example.com/example.rb"
  retryable           = true
  retry_count         = 1
  retry_delay_seconds = 10
  allow_custom_config = true
}