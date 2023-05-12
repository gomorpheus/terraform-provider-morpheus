resource "morpheus_email_task" "tfexample_email_git" {
  name                        = "tfexample_email_git"
  code                        = "tfexample_email_git"
  labels                      = ["demo", "terraform"]
  email_address               = "<%=instance.createdByEmail%>"
  subject                     = "<%=instance.hostname%> provisioning complete"
  source                      = "repository"
  content_path                = "example.txt"
  repository_id               = 1
  version_ref                 = "main"
  skip_wrapped_email_template = false
  retryable                   = true
  retry_count                 = 1
  retry_delay_seconds         = 10
  allow_custom_config         = true
}