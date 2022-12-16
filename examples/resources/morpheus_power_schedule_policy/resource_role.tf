resource "morpheus_power_schedule_policy" "tf_example_power_schedule_policy_role" {
  name                         = "tf_example_power_schedule_policy_role"
  description                  = "terraform example role power schedule policy"
  enabled                      = true
  enforcement_type             = "fixed"
  power_schedule_id            = 2
  hide_power_schedule_if_fixed = true
  scope                        = "role"
  role_id                      = 1
  apply_to_each_user           = true
}