resource "morpheus_number_option_type" "tf_example_number_option_type" {
  name                     = "tf_example_number_option_type"
  description              = "Terraform number option type example"
  field_name               = "number_example"
  export_meta              = true
  dependent_field          = "dependent_example"
  visibility_field         = "visibility_example"
  display_value_on_details = true
  field_label              = "Number Example"
  placeholder              = "12"
  default_value            = "1"
  help_block               = "Provide a number"
  required                 = true
}