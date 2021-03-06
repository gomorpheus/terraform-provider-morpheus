---
page_title: "morpheus_checkbox_option_type Resource - terraform-provider-morpheus"
subcategory: ""
description: |-
  Provides a Morpheus checkbox option type resource
---

# morpheus_checkbox_option_type

Provides a Morpheus checkbox option type resource

## Example Usage

```terraform
resource "morpheus_checkbox_option_type" "tf_example_checkbox_option_type" {
  name                     = "tfcheckboxexample"
  description              = "Terraform checkbox option type example"
  field_name               = "checkbox_example"
  export_meta              = true
  dependent_field          = "dependent_example"
  visibility_field         = "visibility_example"
  display_value_on_details = true
  field_label              = "Checkbox Example"
  default_checked          = "on"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **field_label** (String) The label associated with the field in the UI
- **field_name** (String) The field name of the checkbox option type
- **name** (String) The name of the checkbox option type

### Optional

- **default_checked** (String) Whether the checkbox option type is checked by default (on or off)
- **dependent_field** (String) The field or code used to trigger the reloading of the field
- **description** (String) The description of the checkbox option type
- **display_value_on_details** (Boolean) Display the selected value of the checkbox option type on the associated resource's details page
- **export_meta** (Boolean) Whether to export the checkbox option type as a tag
- **visibility_field** (String) The field or code used to trigger the visibility of the field

### Read-Only

- **id** (String) The ID of the checkbox option type

## Import

Import is supported using the following syntax:

```shell
terraform import morpheus_checkbox_option_type.tf_example_checkbox_option_type 1
```
