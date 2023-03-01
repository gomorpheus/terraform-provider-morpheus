resource "morpheus_text_option_type" "tf_example_text_option_type" {
  name                     = "tf_example_text_option_type"
  description              = "Terraform text option type example"
  labels                   = ["demo", "terraform"]
  field_name               = "test1"
  export_meta              = true
  dependent_field          = "dependent_example"
  visibility_field         = "visibility_example"
  require_field            = "require_example"
  show_on_edit             = true
  editable                 = true
  display_value_on_details = true
  field_label              = "numbers"
  placeholder              = "fewf"
  default_value            = "testing"
  help_block               = "fiwefw"
  required                 = true
  verify_pattern           = "a\\D{4}"
}