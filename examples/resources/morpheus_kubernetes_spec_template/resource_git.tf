resource "morpheus_kubernetes_spec_template" "tfexample_kubernetes_spec_template_git" {
  name          = "tf-kubernetes-spec-example-git"
  source_type   = "repository"
  repository_id = 2
  version_ref   = "main"
  spec_path     = "./spec.yaml"
}