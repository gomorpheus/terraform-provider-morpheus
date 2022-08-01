resource "morpheus_backup_creation_policy" "tf_example_backup_creation_policy_role" {
  name             = "tf_example_backup_creation_policy_role"
  description      = "tfvsphere"
  enabled          = true
  enforcement_type = "fixed"
  create_backup    = true
  scope            = "role"
  role_id          = 1
  apply_each_user  = true
}