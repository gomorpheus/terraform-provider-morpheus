resource "morpheus_cloud_formation_app_blueprint" "tf_example_cloud_formation_app_blueprint_git" {
  name                   = "example_cloud_formation_app_blueprint_git"
  description            = "Example cloud formation app blueprint"
  category               = "cloudformation"
  install_agent          = true
  cloud_init_enabled     = true
  capability_iam         = true
  capability_named_iam   = true
  capability_auto_expand = true
  source_type            = "repository"
  working_path           = "./test"
  integration_id         = 3
  repository_id          = 1
  version_ref            = "main"
}