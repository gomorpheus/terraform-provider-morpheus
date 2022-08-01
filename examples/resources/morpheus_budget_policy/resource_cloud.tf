resource "morpheus_budget_policy" "tf_example_budget_policy_cloud" {
  name         = "tf_example_budget_policy_cloud"
  description  = "terraform example cloud budget policy"
  enabled      = true
  max_price    = "4000"
  currency     = "USD"
  unit_of_time = "hour"
  scope        = "cloud"
  cloud_id     = 1
}