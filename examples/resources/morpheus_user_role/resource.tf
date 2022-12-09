resource "morpheus_user_role" "tf_example_user_role" {
  name = "tf_example_user_role"
  description = "Terraform example user role"
  multitenant_role = true
  multitenant_locked = true
  permission {
    global_groups_permission = "full"
    global_instance_types_permission = "full"
    global_blueprints_permission = "full"
    global_report_types_permission = "full"
    global_catalog_item_types_permission = "full"
    global_vdi_pools_permission = "full"
    default_persona = "catalog"

    feature_permission {
      name = "Admin: Appliance Settings"
      access = "read"
    }

    feature_permission {
      name = "Admin: Appliance Settings"
      access = "read"
    }
  }
}