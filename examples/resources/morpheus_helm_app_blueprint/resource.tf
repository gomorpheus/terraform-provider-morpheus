resource "morpheus_helm_app_blueprint" "tf_helm_app_blueprint" {
  name           = "helmappblueprint"
  description    = "testing terraform"
  category       = "helmapps"
  integration_id = 3
  repository_id  = 1
  version_ref    = "main"
  working_path   = "./test"
}
