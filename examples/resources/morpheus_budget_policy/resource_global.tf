resource "morpheus_budget_policy" "tf_example_budget_policy_global" {
  name         = "tf_example_budget_policy_global"
  description  = "terraform example global budget policy"
  enabled      = true
  max_price    = "4000"
  currency     = "USD"
  unit_of_time = "hour"
  scope        = "global"
}