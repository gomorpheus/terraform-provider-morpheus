resource "morpheus_radio_list_option_type" "tf_example_radio_list_option_type" {
  name                     = "tf_example_radio_list_option_type"
  description              = "Terraform radio list option type example"
  field_name               = "radioExample"
  export_meta              = true
  dependent_field          = "upstreamExample"
  visibility_field         = "upstreamExample"
  require_field            = "upstreamExample"
  show_on_edit             = true
  editable                 = true
  display_value_on_details = true
  field_label              = "Radio Example"
  default_value            = "example"
  help_block               = "Terraform radio list option type example"
  required                 = true
  option_list_id           = 3
}