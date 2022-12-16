resource "morpheus_cypher_access_policy" "tf_example_cypher_access_policy_role" {
  name               = "tf_example_cypher_access_policy_role"
  description        = "terraform example role cypher access policy"
  enabled            = true
  key_path           = ".*"
  read_access        = true
  write_access       = true
  update_access      = true
  list_access        = true
  delete_access      = true
  scope              = "role"
  role_id            = 1
  apply_to_each_user = true
}