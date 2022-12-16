resource "morpheus_email_task" "tfexample_email" {
  name                        = "tfexample_email"
  code                        = "tfexample_email"
  email_address               = "<%=instance.createdByEmail%>"
  subject                     = "<%=instance.hostname%> provisioning complete"
  source                      = "local"
  content                     = "Your instance <%=instance.hostname%> was provisioned."
  skip_wrapped_email_template = false
  retryable                   = true
  retry_count                 = 1
  retry_delay_seconds         = 10
  allow_custom_config         = true
}