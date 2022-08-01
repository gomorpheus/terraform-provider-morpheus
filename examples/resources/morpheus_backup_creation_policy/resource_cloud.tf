resource "morpheus_backup_creation_policy" "tf_example_backup_creation_policy_cloud" {
  name             = "tf_example_backup_creation_policy_cloud"
  description      = "tfvsphere"
  enabled          = true
  enforcement_type = "fixed"
  create_backup    = true
  scope            = "cloud"
  cloud_id         = 1
}