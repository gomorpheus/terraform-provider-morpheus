resource "morpheus_password_option_type" "tf_example_password_option_type" {
  name                     = "tf_example_password_option_type"
  description              = "Terraform password option type example"
  field_name               = "test1"
  export_meta              = true
  dependent_field          = "servicemsh"
  visibility_field         = "demotestin"
  display_value_on_details = true
  field_label              = "numbers"
  placeholder              = "fewf"
  default_value            = "testing"
  help_block               = "fiwefw"
  required                 = true
}