resource "morpheus_user" "tf_output_writer_user" {
  username              = "tfwriter"
  first_name            = "tfwriter"
  last_name             = "tfwriter"
  email                 = "test@test.local"
  password              = "PmWFEAE#92331"
  role_ids              = [morpheus_user_role.terrform_output_user_role.id]
  receive_notifications = false
}

resource "morpheus_user" "tf_output_reader_user" {
  username              = "tfreader"
  first_name            = "tfreader"
  last_name             = "tfreader"
  email                 = "test@test.local"
  password              = "PmWFEAE#92331"
  role_ids              = [morpheus_user_role.terrform_output_user_role.id]
  receive_notifications = false
}

resource "morpheus_cypher_access_policy" "tf_writer_policy" {
  name          = "tf_writer_policy"
  description   = "terraform example user cypher access policy"
  enabled       = true
  key_path      = ".*"
  read_access   = true
  write_access  = false
  update_access = true
  list_access   = false
  delete_access = true
  scope         = "user"
  user_id       = morpheus_user.tf_output_writer_user.id
}

resource "morpheus_cypher_access_policy" "tf_reader_policy" {
  name          = "tf_reader_policy"
  description   = "terraform example user cypher access policy"
  enabled       = true
  key_path      = ".*"
  read_access   = true
  write_access  = true
  update_access = true
  list_access   = true
  delete_access = true
  scope         = "user"
  user_id       = morpheus_user.tf_output_reader_user.id
}

data "morpheus_permission_set" "tf_output_permissions" {
  default_group_permission             = "none"
  default_instance_type_permission     = "none"
  default_blueprint_permission         = "none"
  default_report_type_permission       = "none"
  default_catalog_item_type_permission = "none"
  default_vdi_pool_permission          = "none"
  default_workflow_permission          = "none"
  default_task_permission              = "none"
  default_persona                      = "standard"
  feature_permission {
    code   = "services-cypher"
    access = "read"
  }
}

resource "morpheus_user_role" "terrform_output_user_role" {
  name               = "terraform-shared-output"
  description        = "Terraform provider example user role"
  multitenant_role   = false
  multitenant_locked = false
  permission_set     = data.morpheus_permission_set.tf_output_permissions.json
}
