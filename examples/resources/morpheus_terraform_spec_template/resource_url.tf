resource "morpheus_terraform_spec_template" "tfexample_terraform_spec_template_url" {
  name                = "tf-terraform-spec-example-url"
  source_type         = "url"
  spec_path           = "http://example.com/spec.tf"
}