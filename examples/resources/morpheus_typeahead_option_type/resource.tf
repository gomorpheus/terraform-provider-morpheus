resource "morpheus_typeahead_option_type" "tf_example_typeahead_option_type" {
  name                     = "tf_example_typeahead_option_type"
  description              = "testing"
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
  option_list_id           = 3
}