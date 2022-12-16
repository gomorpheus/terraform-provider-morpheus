resource "morpheus_power_schedule_policy" "tf_example_power_schedule_policy_group" {
  name                         = "tf_example_power_schedule_policy_group"
  description                  = "terraform example group power schedule policy"
  enabled                      = true
  enforcement_type             = "fixed"
  power_schedule_id            = 2
  hide_power_schedule_if_fixed = true
  scope                        = "group"
  group_id                     = 1
}