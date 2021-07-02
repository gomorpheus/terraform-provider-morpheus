resource "morpheus_terraform_spec_template" "tfexample_terraform_spec_template_git" {
  name                = "tf-terraform-spec-example-git"
  source_type         = "repository"
  repository_id       = 2
  version_ref         = "main"
  spec_path           = "Instance Types/Terraform/CloudResource/aws/vpc.tf"
}