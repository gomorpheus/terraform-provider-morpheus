resource "morpheus_groovy_script_task" "tfexample_groovy_local" {
  name                = "tfexample_groovy_local"
  code                = "tfexample_groovy_local"
  source_type         = "local"
  script_content      = <<EOF
println "hello"
EOF
  retryable           = true
  retry_count         = 1
  retry_delay_seconds = 10
  allow_custom_config = true
}