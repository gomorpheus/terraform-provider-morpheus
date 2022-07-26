resource "morpheus_kubernetes_app_blueprint" "tfexample_kubernetes_app_blueprint_git" {
  name           = "tf-kubernetes-spec-example-git"
  description    = ""
  category       = ""
  source_type    = "repository"
  integration_id = 3
  repository_id  = 1
  version_ref    = "main"
  working_path   = "./test"
}