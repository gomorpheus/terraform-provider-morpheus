resource "morpheus_helm_app_blueprint" "tf_example_helm_app_blueprint" {
  name           = "helmappblueprint"
  description    = "tf example helm app blueprint"
  category       = "helm"
  integration_id = 3
  repository_id  = 1
  version_ref    = "main"
  working_path   = "./test"
}
