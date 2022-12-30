resource "morpheus_textarea_option_type" "tf_example_textarea_option_type" {
  name                     = "tf_example_textarea_option_type"
  description              = "Terraform text area option type example"
  field_name               = "textareaExample"
  export_meta              = true
  dependent_field          = "upstreamExample"
  visibility_field         = "upstreamExample"
  require_field            = "upstreamExample"
  show_on_edit             = true
  editable                 = true
  display_value_on_details = true
  field_label              = "Text Area Example"
  rows                     = "5"
  placeholder              = "example text"
  default_value            = "example"
  help_block               = "Terraform text area option type example"
  required                 = true
}