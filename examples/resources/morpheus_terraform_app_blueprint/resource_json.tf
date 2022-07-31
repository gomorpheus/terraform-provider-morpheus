resource "morpheus_terraform_app_blueprint" "tfapp_blueprint_json" {
  name              = "tfappbluedemojson"
  description       = "testing terraform"
  category          = "terraformdemo"
  source_type       = "json"
  blueprint_content = <<EOF
{"test":"demo123"}
EOF
  terraform_version = "1.1.1"
  terraform_options = "-var 'foo=bar'"
  tfvar_secret      = "tfvars/rdsdemo-secrets"
  visibility        = "public"
}