---
page_title: "morpheus_radio_list_option_type Resource - terraform-provider-morpheus"
subcategory: ""
description: |-
  Provides a Morpheus radio list option type resource
---

# morpheus_radio_list_option_type

Provides a Morpheus radio list option type resource

## Example Usage

```terraform
resource "morpheus_radio_list_option_type" "tf_example_radio_list_option_type" {
  name                     = "tf_example_radio_list_option_type"
  description              = "Terraform radio list option type example"
  labels                   = ["demo", "terraform"]
  field_name               = "radioExample"
  export_meta              = true
  dependent_field          = "dependent_example"
  visibility_field         = "visibility_example"
  require_field            = "require_example"
  show_on_edit             = true
  editable                 = true
  display_value_on_details = true
  field_label              = "Radio Example"
  default_value            = "example"
  help_block               = "Terraform radio list option type example"
  required                 = true
  option_list_id           = 3
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the radio list option type

### Optional

- `default_value` (String) The default value of the option type
- `dependent_field` (String) The field or code used to trigger the reloading of the field
- `description` (String) The description of the radio list option type
- `display_value_on_details` (Boolean) Display the selected value of the radio list option type on the associated resource's details page
- `editable` (Boolean) Whether the value of the option type can be edited after the initial request
- `export_meta` (Boolean) Whether to export the radio list option type as a tag
- `field_label` (String) The label associated with the field in the UI
- `field_name` (String) The field name of the radio list option type
- `help_block` (String) Text that provides additional details about the use of the option type
- `labels` (Set of String) The organization labels associated with the option type(Only supported on Morpheus 5.5.3 or higher)
- `option_list_id` (Number) The ID of the associated option list
- `require_field` (String) The field or code used to trigger the required status of the field
- `required` (Boolean) Whether the option type is required
- `show_on_edit` (Boolean) Whether the option type will display in the edit section of the provisioned resource
- `visibility_field` (String) The field or code used to trigger the visibility of the field

### Read-Only

- `id` (String) The ID of the radio list option type

## Import

Import is supported using the following syntax:

```shell
terraform import morpheus_radio_list_option_type.tf_example_radio_list_option_type 1
```
