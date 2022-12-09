resource "morpheus_cypher_access_policy" "tf_example_cypher_access_policy_global" {
  name          = "tf_example_cypher_access_policy_global"
  description   = "terraform example global cypher access policy"
  enabled       = true
  key_path      = ".*"
  read_access   = true
  write_access  = true
  update_access = true
  list_access   = true
  delete_access = true
  scope         = "global"
}