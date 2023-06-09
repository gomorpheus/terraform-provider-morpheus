data "morpheus_git_integration" "tf_example_git_integration" {
  name = "MorpheusAutomation"
}

resource "morpheus_shell_script_task" "tfexample_shell_local" {
  name           = "tfexample_shell_local"
  code           = "tfexample_shell_local"
  labels         = ["demo", "terraform"]
  result_type    = "json"
  source_type    = "repository"
  script_path    = "Shell/preinstallcheck.sh"
  version_ref    = "main"
  repository_id  = data.morpheus_git_integration.tf_example_git_integration.repository_ids["morpheus-automation-examples"]
  execute_target = "resource"
}