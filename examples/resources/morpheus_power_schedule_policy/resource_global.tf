resource "morpheus_power_schedule_policy" "tf_example_power_schedule_policy_global" {
  name                         = "tf_example_power_schedule_policy_global"
  description                  = "terraform example global power schedule policy"
  enabled                      = true
  enforcement_type             = "fixed"
  power_schedule_id            = 2
  hide_power_schedule_if_fixed = true
  scope                        = "global"
}