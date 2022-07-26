resource "morpheus_arm_spec_template" "tfexample_arm_spec_template_git" {
  name                = "tf-arm-spec-example-git"
  source_type         = "repository"
  repository_id       = 2
  version_ref         = "main"
  spec_path           = ""
}