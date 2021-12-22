resource "morpheus_execute_schedule" "tf_example_execute_schedule" {
  name        = "Run daily at 7 AM"
  description = "This schedule runs daily at 7 AM Mountain Time"
  enabled     = false
  time_zone   = "America/Denver"
  schedule    = "7 0 * * *"
}