---
page_title: "morpheus_typeahead_option_type Resource - terraform-provider-morpheus"
subcategory: ""
description: |-
  Provides a Morpheus typeahead option type resource
---

# morpheus_typeahead_option_type

Provides a Morpheus typeahead option type resource

## Example Usage

```terraform
resource "morpheus_typeahead_option_type" "tf_example_cloud" {
  name        = "morpheus_vsphere"
  description = ""
  type        = "vmware"
  code        = ""
  visibility  = "public"
  enabled     = true
  config     = ""
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **name** (String) The name of the typeahead option type

### Optional

- **default_value** (String) The default value of the option type
- **dependent_field** (String) The field or code used to trigger the reloading of the field
- **description** (String) The description of the typeahead option type
- **display_value_on_details** (Boolean) Display the selected value of the text option type on the associated resource's details page
- **export_meta** (Boolean) Whether to export the text option type as a tag
- **field_label** (String) The label associated with the field in the UI
- **field_name** (String) The field name of the typeahead option type
- **help_block** (String) Text that provides additional details about the use of the option type
- **id** (String) The ID of this resource.
- **option_list_id** (Number) The ID of the associated option list
- **placeholder** (String) Text in the field used as a placeholder for example purposes
- **required** (Boolean) Whether the option type is required
- **visibility_field** (String) The field or code used to trigger the visibility of the field

## Import

Import is supported using the following syntax:

```shell
terraform import morpheus_typeahead_option_type.tf_example_typeahead_option_type 1
```