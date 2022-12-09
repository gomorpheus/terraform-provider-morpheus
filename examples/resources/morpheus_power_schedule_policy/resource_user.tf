resource "morpheus_power_schedule_policy" "tf_example_power_schedule_policy_user" {
  name                         = "tf_example_power_schedule_policy_user"
  description                  = "terraform example user power schedule policy"
  enabled                      = true
  enforcement_type             = "fixed"
  power_schedule_id            = 2
  hide_power_schedule_if_fixed = true
  scope                        = "user"
  user_id                      = 1
}