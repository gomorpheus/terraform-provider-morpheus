resource "morpheus_kubernetes_app_blueprint" "tfexample_kubernetes_app_blueprint_git" {
  name           = "tf-kubernetes-app-blueprint-example-git"
  description    = "tf example kubernetes app blueprint"
  category       = "k8s"
  source_type    = "repository"
  integration_id = 3
  repository_id  = 1
  version_ref    = "main"
  working_path   = "./test"
}