resource "morpheus_form" "tf_example_form" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

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
    group_code  = "group_name"
  }

  field_group {
    name = "Configuration"

    option_type {
      name               = "debian_layout"
      code               = "debian_layout"
      field_name         = "f_layout"
      field_label        = "Template"
      type               = "layout"
      instance_type_code = "debian"
      required           = true
      group_code         = "group_name"
      cloud_code         = "cloud_provider"
    }

    option_type {
      name        = "plan_choice"
      code        = "plan_choice"
      field_name  = "f_plan"
      field_label = "Gabarit"
      help_block  = "Gabarit"
      type        = "plan"
      required    = true
      group_code  = "group_name"
      cloud_code  = "cloud_provider"
      layout_code = "debian_layout"
      pool_code   = "pool_choice"
    }

    option_type {
      name          = "labels"
      code          = "labels"
      field_name    = "f_tags"
      field_label   = "Labels"
      help_block    = "Labels"
      type          = "tag"
      required      = false
      default_value = "[\n {\n  \"name\": \"lab\",\n  \"value\": \"sandbox\"\n }\n]"
    }

    option_type {
      name        = "pool_choice"
      code        = "pool_choice"
      field_name  = "f_pool"
      field_label = "Resource Pool"
      help_block  = "Resource Pool (RAM, CPU...) to use"
      type        = "resourcePool"
      required    = true
      hidden      = true
      group_code  = "group_name"
      cloud_code  = "cloud_provider"
      layout_code = "debian_layout"
      plan_code   = "plan_choice"
    }
  }

  field_group {
    name = "Exposition"

    option_type {
      name        = "network_interface"
      code        = "network_interface"
      field_name  = "f_network"
      field_label = "Network"
      type        = "networkManager"
      required    = false
      group_code  = "group_name"
      cloud_code  = "cloud_provider"
      layout_code = "debian_layout"
      pool_code   = "pool_choice"
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
