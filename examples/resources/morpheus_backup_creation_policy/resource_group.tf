resource "morpheus_backup_creation_policy" "tf_example_backup_creation_policy_group" {
  name             = "tf_example_backup_creation_policy_group"
  description      = "tfvsphere"
  enabled          = true
  enforcement_type = "fixed"
  create_backup    = true
  scope            = "group"
  group_id         = 1
}