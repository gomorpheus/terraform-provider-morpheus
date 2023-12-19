resource "morpheus_form" "mem" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  option_types = []

  field_group {
    name                 = "fg1"
    description          = "testin"
    collapsible          = true
    collapsed_by_deafult = true
    visibility_field     = "testing"
    option_types      = [
        data.morpheus_option_type.versions.json,
        data.morpheus_form_option_type.test.json
    ]
  }
}

data "morpheus_option_type" "versions" {
  name = "App Versions"
}

data "morpheus_form_option_type" "test" {
  type        = ""
  name        = ""
  description = ""
}