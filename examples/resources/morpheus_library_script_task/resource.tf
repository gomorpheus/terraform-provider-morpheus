data "morpheus_script_template" "example_script_template" {
  name = "My library script task"
}

resource "morpheus_library_script_task" "tf_example_library_script_task" {
  name                = "Example Terraform Libary Script Task"
  code                = "tf-example-library-script-task"
  labels              = ["demo", "library", "terraform"]
  execute_target      = "resource"
  script_template     = data.morpheus_script_template.example_script_template.name
  script_template_id  = data.morpheus_script_template.example_script_template.id
  retryable           = true
  retry_count         = 1
  retry_delay_seconds = 10
  allow_custom_config = true
}