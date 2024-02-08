resource "morpheus_form" "tf_example_form" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  option_type {
    id = 2182
  }

  option_type {
    name = "test select"
    code = "select-input"
    description = "Testing stuff"
    type = "select"
    field_label = "Select Test"
    field_name = "selectTest"
    default_value = "Demo123"
    placeholder = "Testing 123"
    help_block = "Is this working now"
    option_list_id = 1
    required = true
    export_meta = true
    display_value_on_details = true
    locked = true
    hidden = true
    exclude_from_search = true
  }

  option_type {
    name = "test select"
    code = "select-input"
    description = "Testing stuff"
    type = "select"
    field_label = "Select Test"
    field_name = "selectTest"
    default_value = "Demo123"
    placeholder = "Testing 123"
    help_block = "Is this working now"
    option_list_id = 1
    required = true
    export_meta = true
    display_value_on_details = true
    locked = true
    hidden = true
    exclude_from_search = true
  }

  option_type {
    name = "test radio"
    code = "radio-input"
    description = "Testing stuff"
    type = "radio"
    field_label = "Radio Test"
    field_name = "radioTest"
    default_value = "Demo123"
    placeholder = "Testing 123"
    help_block = "Is this working now"
    option_list_id = 1
    required = true
    export_meta = true
    display_value_on_details = true
    locked = true
    hidden = true
    exclude_from_search = true
  }

  option_type {
    name = "test text"
    code = "test-input"
    description = "Testing stuff"
    type = "text"
    field_label = "Testin"
    field_name = "test"
    default_value = "Demo123"
    placeholder = "Testing 123"
    help_block = "Is this working now"
    required = true
    export_meta = true
    display_value_on_details = true
    locked = true
    hidden = true
    exclude_from_search = true
  }

  option_type {
    name = "checkbox input"
    code = "checkbox-input"
    description = "Testing stuff"
    type = "checkbox"
    field_label = "checkbox input"
    field_name = "checkboxInput"
    default_value = "test"
    placeholder = "Testing 123"
    help_block = "Is this working now"
    required = true
    export_meta = true
    display_value_on_details = true
    locked = true
    hidden = true
    exclude_from_search = true
  }

  option_type {
    name = "hidden input"
    code = "hidden-input"
    description = "Testing stuff"
    type = "hidden"
    field_label = "hidden input"
    field_name = "hiddenInput"
    default_value = "test"
    placeholder = "Testing 123"
    help_block = "Is this working now"
    required = true
    export_meta = true
    display_value_on_details = true
    locked = true
    hidden = true
    exclude_from_search = true
  }

  option_type {
    name = "number input"
    code = "number-input"
    description = "Testing stuff"
    type = "number"
    field_label = "number input"
    field_name = "numberInput"
    default_value = "4"
    placeholder = "Testing 123"
    help_block = "Is this working now"
    required = true
    export_meta = true
    display_value_on_details = true
    locked = true
    hidden = true
    exclude_from_search = true
    min_value = 3
    max_value = 44
    step = 2
  }

  field_group {
    name                 = "fg2"
    description          = "testin"
    collapsible          = true
    collapsed_by_deafult = true
    //   visibility_field     = "testing"
    //    option_type_ids      = []
    option_type {
      id = 2182
    }
  }

  field_group {
    name                 = "fg1"
    description          = "testin"
    collapsible          = true
    collapsed_by_deafult = true
    //    visibility_field     = "testing"
  }

}