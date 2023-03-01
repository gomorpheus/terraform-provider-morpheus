resource "morpheus_select_list_option_type" "tf_example_select_list_option_type" {
  name                     = "tf_example_select_list_option_type"
  description              = "Terraform select list option type example"
  labels                   = ["demo", "terraform"]
  field_name               = "tfSelectExample"
  export_meta              = true
  dependent_field          = "dependent_example"
  visibility_field         = "visibility_example"
  require_field            = "require_example"
  show_on_edit             = true
  editable                 = true
  display_value_on_details = true
  field_label              = "numbers"
  default_value            = "testing"
  help_block               = "fiwefw"
  required                 = true
  option_list_id           = 3
}