resource "morpheus_hidden_option_type" "tf_example_hidden_option_type" {
  name                     = "tf_example_hidden_option_type"
  description              = "Terraform hidden option type example"
  field_name               = "hidden_example"
  export_meta              = true
  dependent_field          = "dependent_example"
  visibility_field         = "visibility_example"
  display_value_on_details = true
  default_value            = "example"
}