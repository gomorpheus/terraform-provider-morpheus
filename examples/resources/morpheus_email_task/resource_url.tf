resource "morpheus_email_task" "tfexample_email_url" {
  name                        = "tfexample_email_url"
  code                        = "tfexample_email_url"
  labels                      = ["demo", "terraform"]
  email_address               = "<%=instance.createdByEmail%>"
  subject                     = "<%=instance.hostname%> provisioning complete"
  source                      = "url"
  content_url                 = "https://example.com/example.txt"
  skip_wrapped_email_template = false
  retryable                   = true
  retry_count                 = 1
  retry_delay_seconds         = 10
  allow_custom_config         = true
}