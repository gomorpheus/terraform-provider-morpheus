resource "morpheus_budget_policy" "tf_example_budget_policy_group" {
  name         = "tf_example_budget_policy_group"
  description  = "terraform example group budget policy"
  enabled      = true
  max_price    = "4000"
  currency     = "USD"
  unit_of_time = "hour"
  scope        = "group"
  group_id     = 1
}