resource "morpheus_hidden_option_type" "tf_example_hidden_option_type" {
  name                     = "tf_example_hidden_option_type"
  description              = "Terraform hidden option type example"
  field_name               = "test1"
  export_meta              = true
  dependent_field          = "servicemsh"
  visibility_field         = "demotestin"
  display_value_on_details = true
  default_value            = "testing"
}