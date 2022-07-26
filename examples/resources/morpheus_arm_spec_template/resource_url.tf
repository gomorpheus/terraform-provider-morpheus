resource "morpheus_arm_spec_template" "tfexample_arm_spec_template_url" {
  name                = "tf-terraform-spec-example-url"
  category            = ""
  source_type         = "url"
  spec_path           = "http://example.com/spec.tf"
}