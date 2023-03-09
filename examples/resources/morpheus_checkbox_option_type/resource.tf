resource "morpheus_checkbox_option_type" "tf_example_checkbox_option_type" {
  name                     = "tfcheckboxexample"
  description              = "Terraform checkbox option type example"
  labels                   = ["demo", "terraform"]
  field_name               = "checkbox_example"
  export_meta              = true
  dependent_field          = "dependent_example"
  visibility_field         = "visibility_example"
  require_field            = "require_example"
  show_on_edit             = true
  editable                 = true
  display_value_on_details = true
  field_label              = "Checkbox Example"
  default_checked          = true
}