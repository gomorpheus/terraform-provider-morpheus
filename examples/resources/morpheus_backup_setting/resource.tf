resource "morpheus_backup_setting" "tf_example_backup_setting" {
  scheduled_backups                = true
  create_backups                   = true
  backup_appliance                 = false
  default_backup_storage_bucket_id = 17
  default_backup_schedule_id       = 3
  retention_days                   = 21
}