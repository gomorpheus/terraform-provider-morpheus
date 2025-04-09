---
page_title: "morpheus_form Resource - terraform-provider-morpheus"
subcategory: ""
description: |-
  Provides a Morpheus form resource
---

# morpheus_form

Provides a Morpheus form resource

!> **Note:** Existing inputs or option types are supported, using __only__ the id field.

## Example Usage

```terraform
resource "morpheus_form" "tf_example_form" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  option_type {
    id = 12345
  }

  option_type {
    name                     = "tf example select"
    code                     = "select-input"
    description              = "Terraform select example"
    type                     = "select"
    field_label              = "Select Test"
    field_name               = "selectTest"
    default_value            = "test123"
    placeholder              = "Testing 123"
    help_block               = "Select an option"
    option_list_id           = 1
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = true
    exclude_from_search      = true
  }

  option_type {
    name                     = "tf radio example"
    code                     = "radio-input"
    description              = "Terraform radio example"
    type                     = "radio"
    field_label              = "Radio Test"
    field_name               = "radioTest"
    default_value            = "Demo123"
    placeholder              = "Testing 123"
    help_block               = "Select an option"
    option_list_id           = 1
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = true
    exclude_from_search      = true
  }

  option_type {
    name                     = "tf text example"
    code                     = "test-input"
    description              = "Terraform text example"
    type                     = "text"
    field_label              = "Testin"
    field_name               = "test"
    default_value            = "Demo123"
    placeholder              = "Testing 123"
    help_block               = "Is this working now"
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = true
    exclude_from_search      = true
  }

  option_type {
    name                       = "tf checkbox example"
    code                       = "checkbox-input"
    description                = "Terraform checkbox example"
    type                       = "checkbox"
    field_label                = "checkbox input"
    field_name                 = "checkboxInput"
    default_chedefault_checked = true
    placeholder                = "Testing 123"
    help_block                 = "Is this working now"
    required                   = true
    export_meta                = true
    display_value_on_details   = true
    locked                     = true
    hidden                     = true
    exclude_from_search        = true
  }

  option_type {
    name                     = "tf hidden input example"
    code                     = "hidden-input"
    description              = "Terraform hidden input example"
    type                     = "hidden"
    field_label              = "hidden input"
    field_name               = "hiddenInput"
    default_value            = "test"
    placeholder              = "Testing 123"
    help_block               = "Is this working now"
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = true
    exclude_from_search      = true
  }

  option_type {
    name                     = "tf number input example"
    code                     = "number-input"
    description              = "Terraform number example"
    type                     = "number"
    field_label              = "number input"
    field_name               = "numberInput"
    default_value            = "4"
    placeholder              = "Testing 123"
    help_block               = "Is this working now"
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = true
    exclude_from_search      = true
    min_value                = 3
    max_value                = 44
    step                     = 2
  }

  field_group {
    name                 = "fg1"
    description          = "testin"
    collapsible          = true
    collapsed_by_deafult = true
    option_type {
      name                     = "tf field group 1 text input example"
      code                     = "test-input"
      description              = "Terraform text input example"
      type                     = "text"
      field_label              = "Testin"
      field_name               = "test"
      default_value            = "Demo123"
      placeholder              = "Testing 123"
      help_block               = "Is this working now"
      required                 = true
      export_meta              = true
      display_value_on_details = true
      locked                   = true
      hidden                   = false
      exclude_from_search      = true
    }
  }

  field_group {
    name                 = "fg2"
    description          = "testin"
    collapsible          = true
    collapsed_by_deafult = true
    option_type {
      name                     = "tf field group 2 text input example"
      code                     = "test-input"
      description              = "Terraform text input example"
      type                     = "text"
      field_label              = "Testin"
      field_name               = "test"
      default_value            = "Demo123"
      placeholder              = "Testing 123"
      help_block               = "Is this working now"
      required                 = true
      export_meta              = true
      display_value_on_details = true
      locked                   = true
      hidden                   = false
      exclude_from_search      = true
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `code` (String) The form code used for API/CLI automation
- `name` (String) The name of the form

### Optional

- `description` (String) A description of the form
- `field_group` (Block List) Field group to add to the form (see [below for nested schema](#nestedblock--field_group))
- `labels` (Set of String) The organization labels associated with the form
- `option_type` (Block List) Form option type (see [below for nested schema](#nestedblock--option_type))

### Read-Only

- `id` (String) The id of the form

<a id="nestedblock--field_group"></a>
### Nested Schema for `field_group`

Required:

- `name` (String) The name of the field group

Optional:

- `collapsed_by_deafult` (Boolean) Whether the field group is collapsed by default
- `collapsible` (Boolean) Whether the field group can be collapsed
- `description` (String) A description of the field group
- `option_type` (Block List) Field group option type (see [below for nested schema](#nestedblock--field_group--option_type))
- `visibility_field` (String) The field or code used to trigger the visibility of the field group

<a id="nestedblock--field_group--option_type"></a>
### Nested Schema for `field_group.option_type`

Optional:

- `allow_duplicates` (Boolean) Whether duplicate selections are allowed
- `allow_multiple_selections` (Boolean) Whether to allow multiple items to be selected when using a select list or type ahead option type
- `allow_password_peek` (Boolean) Whether the value of the password option type can be revealed by the user to ensure they correctly entered the password
- `code` (String) The code of the option type to add to the field group
- `code_language` (String) The coding language used for highlighting code syntax
- `custom_data` (String) Custom JSON data payload to pass (Must be a JSON string)
- `default_checked` (Boolean) Whether the checkbox option type is checked by default
- `default_value` (String) The default value of the option type
- `delimiter` (String) The delimiter used to separate text array input values
- `dependent_field` (String) The field or code used to trigger the reloading of the field
- `description` (String) A description of the option type to add to the field group
- `display` (String) The memory or storage value to use (GB or MB)
- `display_value_on_details` (Boolean) Display the selected value of the option type on the associated resource's details page
- `exclude_from_search` (Boolean) Whether the option type should be execluded from search or not
- `export_meta` (Boolean) Whether to export the option type as a tag
- `field_label` (String) The label of the option type
- `field_name` (String) The field name of the option type to add to the field group
- `help_block` (String) The help block text for the option type
- `hidden` (Boolean) Whether the option type is hidden or not
- `lock_display` (Boolean) Whether to lock the display or not
- `locked` (Boolean) Whether the option type is locked or not
- `max_value` (Number) The maximum value that can be provided for a number option type
- `min_value` (Number) The minimum number that can be selected for a number option type
- `name` (String) The name of the option type to add to the field group
- `option_list_id` (Number) The id of the option list for option types such as a typeahead or select list
- `placeholder` (String) The placeholder text for the option type
- `remove_select_option` (Boolean) For Select List-type Inputs. When marked, the Input will default to the first item in the list rather than to an empty selection
- `require_field` (String) The field or code used to determine whether the field is required or not
- `required` (Boolean) Whether the option type is required or not
- `show_line_numbers` (Boolean) Whether to show the line numbers for the code editor option type
- `sortable` (Boolean) Whether the selected options can be sorted or not
- `step` (Number) The incrementation number used for the number option type (i.e. - 5s, 10s, 100s, etc.)
- `text_rows` (Number) The number of rows to display for a text area
- `type` (String) The type of option type to add to the field group (checkbox, hidden, number, password, radio, select, text, textarea, byteSize, code-editor, fileContent, logoSelector, textArray, typeahead, environment)
- `verify_pattern` (String) The regex pattern used to validate the entered text
- `visibility_field` (String) The field or code used to trigger the visibility of the field



<a id="nestedblock--option_type"></a>
### Nested Schema for `option_type`

Optional:

- `allow_duplicates` (Boolean) Whether duplicate selections are allowed
- `allow_multiple_selections` (Boolean) Whether to allow multiple items to be selected when using a select list or type ahead option type
- `allow_password_peek` (Boolean) Whether the value of the password option type can be revealed by the user to ensure they correctly entered the password
- `code` (String) The code of the option type to add to the form
- `code_language` (String) The coding language used for highlighting code syntax
- `custom_data` (String) Custom JSON data payload to pass (Must be a JSON string)
- `default_checked` (Boolean) Whether the checkbox option type is checked by default
- `default_value` (String) The default value of the option type
- `delimiter` (String) The delimiter used to separate text array input values
- `dependent_field` (String) The field or code used to trigger the reloading of the field
- `description` (String) A description of the option type to add to the form
- `display` (String) The memory or storage value to use (GB or MB)
- `display_value_on_details` (Boolean) Display the selected value of the option type on the associated resource's details page
- `exclude_from_search` (Boolean) Whether the option type should be execluded from search or not
- `export_meta` (Boolean) Whether to export the option type as a tag
- `field_label` (String) The label used for the option type
- `field_name` (String) The name of the option type field to add to the form
- `help_block` (String) The help message displayed below the option type
- `hidden` (Boolean) Whether to display the option type to the user
- `lock_display` (Boolean) Whether to lock the display or not
- `locked` (Boolean) Whether the option type is locked or not
- `max_value` (Number) The maximum value that can be provided for a number option type
- `min_value` (Number) The minimum number that can be selected for a number option type
- `name` (String) The name of the option type to add to the form
- `option_list_id` (Number) The id of the option list for option types such as a typeahead or select list
- `placeholder` (String) The placeholder text used for the option type
- `remove_select_option` (Boolean) For Select List-type Inputs. When marked, the Input will default to the first item in the list rather than to an empty selection
- `require_field` (String) The field or code used to determine whether the field is required or not
- `required` (Boolean) Whether the option type is required or not
- `show_line_numbers` (Boolean) Whether to show the line numbers for the code editor option type
- `sortable` (Boolean) Whether the selected options can be sorted or not
- `step` (Number) The incrementation number used for the number option type (i.e. - 5s, 10s, 100s, etc.)
- `text_rows` (Number) The number of lines to show for a code editor or text area option type
- `type` (String) The type of option type to add to the form (checkbox, hidden, number, password, radio, select, text, textarea, byteSize, code-editor, fileContent, logoSelector, textArray, typeahead, environment)
- `verify_pattern` (String) The regex pattern used to validate the entered text
- `visibility_field` (String) The field or code used to trigger the visibility of the field

## Import

Import is supported using the following syntax:

```shell
terraform import morpheus_form.tf_example_form 1
```
