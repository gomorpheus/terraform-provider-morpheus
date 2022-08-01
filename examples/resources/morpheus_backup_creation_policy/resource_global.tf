resource "morpheus_backup_creation_policy" "tf_example_backup_creation_policy_global" {
  name             = "tf_example_backup_creation_policy_global"
  description      = "tfvsphere"
  enabled          = true
  enforcement_type = "fixed"
  create_backup    = true
  scope            = "global"
}