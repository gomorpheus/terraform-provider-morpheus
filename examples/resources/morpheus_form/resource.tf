resource "morpheus_form" "tf_example_form" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  option_type {
    id = 2182
  }

  option_type {
    name                     = "test select"
    code                     = "select-input"
    description              = "Testing stuff"
    type                     = "select"
    field_label              = "Select Test"
    field_name               = "selectTest"
    default_value            = "Demo123"
    placeholder              = "Testing 123"
    help_block               = "Is this working now"
    option_list_id           = 1
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = true
    exclude_from_search      = true
  }

  option_type {
    name                     = "test select"
    code                     = "select-input"
    description              = "Testing stuff"
    type                     = "select"
    field_label              = "Select Test"
    field_name               = "selectTest"
    default_value            = "Demo123"
    placeholder              = "Testing 123"
    help_block               = "Is this working now"
    option_list_id           = 1
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = true
    exclude_from_search      = true
  }

  option_type {
    name                     = "test radio"
    code                     = "radio-input"
    description              = "Testing stuff"
    type                     = "radio"
    field_label              = "Radio Test"
    field_name               = "radioTest"
    default_value            = "Demo123"
    placeholder              = "Testing 123"
    help_block               = "Is this working now"
    option_list_id           = 1
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = true
    exclude_from_search      = true
  }

  option_type {
    name                     = "test text"
    code                     = "test-input"
    description              = "Testing stuff"
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
    name                     = "checkbox input"
    code                     = "checkbox-input"
    description              = "Testing stuff"
    type                     = "checkbox"
    field_label              = "checkbox input"
    field_name               = "checkboxInput"
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
    name                     = "hidden input"
    code                     = "hidden-input"
    description              = "Testing stuff"
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
    name                     = "number input"
    code                     = "number-input"
    description              = "Testing stuff"
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

  option_type {
    name        = "group_name"
    code        = "group_name"
    field_name  = "f_group"
    field_label = "Name of the group"
    type        = "group"
    required    = true
  }

  option_type {
    name        = "cloud_provider"
    code        = "cloud_provider"
    field_name  = "f_cloud"
    field_label = "Provider of the Cloud"
    type        = "cloud"
    required    = true
    group_code = "group_name"
  }

  field_group {
    name        = "Configuration"

    option_type {
      name        = "debian_layout"
      code        = "debian_layout"
      field_name  = "f_layout"
      field_label = "Template"
      type        = "layout"
      instance_type_code = "debian"
      required    = true
      group_code = "group_name"
      cloud_code = "cloud_provider"
    }

    option_type {
      name        = "plan_choice"
      code        = "plan_choice"
      field_name  = "f_plan"
      field_label = "Gabarit"
      help_block = "Gabarit"
      type        = "plan"
      required    = true
      group_code = "group_name"
      cloud_code = "cloud_provider"
      layout_code = "debian_layout" 
      pool_code = "pool_choice"
    }

    option_type {
      name        = "labels"
      code        = "labels"
      field_name  = "f_tags"
      field_label = "Labels"
      help_block = "Labels"
      type        = "tag"
      required    = false
      default_value = "[\n {\n  \"name\": \"lab\",\n  \"value\": \"sandbox\"\n }\n]" 
    }

    option_type {
      name        = "pool_choice"
      code        = "pool_choice"
      field_name  = "f_pool"
      field_label = "Resource Pool"
      help_block = "Resource Pool (RAM, CPU...) to use"
      type        = "resourcePool"
      required    = true
      hidden   = true
      group_code = "group_name"
      cloud_code = "cloud_provider"
      layout_code = "debian_layout" 
      plan_code = "plan_choice"
    }
  }

  field_group {
    name        = "Exposition"

    option_type {
      name        = "network_interface"
      code        = "network_interface"
      field_name  = "f_network"
      field_label = "Network"
      type        = "networkManager"
      required    = false
      group_code = "group_name"
      cloud_code = "cloud_provider"
      layout_code = "debian_layout"
      pool_code = "pool_choice"
    }
  }

}