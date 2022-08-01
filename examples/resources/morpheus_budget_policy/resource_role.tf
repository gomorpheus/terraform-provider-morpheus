resource "morpheus_budget_policy" "tf_example_budget_policy_role" {
  name            = "tf_example_budget_policy_role"
  description     = "terraform example role budget policy"
  enabled         = true
  max_price       = "4000"
  currency        = "USD"
  unit_of_time    = "hour"
  scope           = "role"
  role_id         = 1
  apply_each_user = true
}