resource "morpheus_cloud_formation_spec_template" "tfexample_cloud_formation_spec_template_git" {
  name          = "tf-cloud-formation-spec-example-git"
  source_type   = "repository"
  repository_id = 2
  version_ref   = "main"
  spec_path     = "./spec.yaml"
}