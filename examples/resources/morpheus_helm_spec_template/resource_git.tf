resource "morpheus_helm_spec_template" "tfexample_helm_spec_template_git" {
  name          = "tf-helm-spec-example-git"
  source_type   = "repository"
  repository_id = 2
  version_ref   = "main"
  spec_path     = "./spec.yaml"
}