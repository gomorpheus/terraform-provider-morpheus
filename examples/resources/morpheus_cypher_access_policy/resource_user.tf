resource "morpheus_cypher_access_policy" "tf_example_cypher_access_policy_user" {
  name          = "tf_example_cypher_access_policy_user"
  description   = "terraform example user cypher access policy"
  enabled       = true
  key_path      = ".*"
  read_access   = true
  write_access  = true
  update_access = true
  list_access   = true
  delete_access = true
  scope         = "user"
  user_id       = 1
}