resource "morpheus_terraform_app_blueprint" "tfapp_blueprint_git" {
  name              = "tfappbluedemogit"
  description       = "testing terraform"
  category          = "terraformdemo"
  source_type       = "repository"
  working_path      = "./test"
  integration_id    = 3
  repository_id     = 1
  version_ref       = "main"
  terraform_version = "1.1.1"
  terraform_options = "-var 'foo=bar'"
  tfvar_secret      = "tfvars/rdsdemo-secrets"
}