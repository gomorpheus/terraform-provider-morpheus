resource "morpheus_terraform_app_blueprint" "tfapp_blueprint_specs" {
  name              = "tfappbluedemospecs"
  description       = "testing terraform"
  category          = "terraformdemo"
  source_type       = "spec"
  spec_template_ids = [81]
  terraform_version = "1.1.1"
  terraform_options = "-var 'foo=bar'"
  tfvar_secret      = "tfvars/rdsdemo-secrets"
}