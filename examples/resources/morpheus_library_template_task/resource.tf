data "morpheus_file_template" "example_file_template" {
  name = "My file template"
}

resource "morpheus_library_template_task" "tf_example_library_template_task" {
  name                = "Example Terraform Library Template Task"
  code                = "tf-example-library-template-task"
  labels              = ["demo", "library", "terraform"]
  execute_target      = "resource"
  file_template_id    = data.morpheus_file_template.example_file_template.id
  retryable           = true
  retry_count         = 1
  retry_delay_seconds = 10
  allow_custom_config = true
}