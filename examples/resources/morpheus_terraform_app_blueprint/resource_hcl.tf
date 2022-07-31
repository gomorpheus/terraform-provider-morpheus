resource "morpheus_terraform_app_blueprint" "tfapp_blueprint" {
  name              = "tfappbluedemo"
  description       = "testing terraform"
  category          = "terraformdemo"
  source_type       = "hcl"
  blueprint_content = <<EOF
variable "master_username" {
  type = string
}

variable "master_password" {
  type      = string
  sensitive = true
}

variable "engine_version" {
  type = string
}

variable "instance_class" {
  type = string
}

resource "local_file" "foo" {
    content  = "foo!"
    filename = "${path.module}/foo.bar"
}
EOF
  terraform_version = "1.1.1"
  terraform_options = "-var 'foo=bar'"
  tfvar_secret      = "tfvars/rdsdemo-secrets"
}