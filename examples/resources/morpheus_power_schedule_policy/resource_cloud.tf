resource "morpheus_power_schedule_policy" "tf_example_power_schedule_policy_cloud" {
  name                         = "tf_example_power_schedule_policy_cloud"
  description                  = "terraform example cloud power schedule policy"
  enabled                      = true
  enforcement_type             = "fixed"
  power_schedule_id            = 2
  hide_power_schedule_if_fixed = true
  scope                        = "cloud"
  cloud_id                     = 1
}