resource "morpheus_budget_policy" "tf_example_budget_policy_user" {
  name         = "tf_example_budget_policy_user"
  description  = "terraform example user budget policy"
  enabled      = true
  max_price    = "4000"
  currency     = "USD"
  unit_of_time = "hour"
  scope        = "user"
  user_id      = 1
}