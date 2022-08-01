resource "morpheus_backup_creation_policy" "tf_example_backup_creation_policy_user" {
  name             = "tf_example_backup_creation_policy_user"
  description      = "tfvsphere"
  enabled          = true
  enforcement_type = "fixed"
  create_backup    = true
  scope            = "user"
  user_id          = 1
}