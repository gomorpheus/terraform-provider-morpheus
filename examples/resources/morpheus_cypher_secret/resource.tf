resource "morpheus_cypher_secret" "tf_example_cypher_secret" {
  key   = "apipassword"
  value = "password123"
  ttl   = 86400
}