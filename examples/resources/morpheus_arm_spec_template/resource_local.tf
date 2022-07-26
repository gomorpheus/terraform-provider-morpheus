resource "morpheus_arm_spec_template" "tfexample_arm_spec_template_local" {
  name                = "tf-terraform-spec-example-local"
  category            = ""
  source_type         = "local"
  spec_path           = "http://example.com/spec.tf"
}