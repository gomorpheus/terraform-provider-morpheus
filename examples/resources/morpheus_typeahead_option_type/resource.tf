resource "morpheus_typeahead_option_type" "tf_example_typeahead_option_type" {
  name                      = "tf_example_typeahead_option_type"
  description               = "terraform example typeahead option type"
  labels                    = ["demo", "terraform"]
  field_name                = "example"
  export_meta               = true
  dependent_field           = "dependent_example"
  visibility_field          = "visibility_example"
  require_field             = "require_example"
  show_on_edit              = true
  editable                  = true
  display_value_on_details  = true
  field_label               = "numbers"
  placeholder               = "enter text here"
  default_value             = "testing"
  help_block                = "terraform example typeahead"
  allow_multiple_selections = true
  required                  = true
  option_list_id            = 3
}