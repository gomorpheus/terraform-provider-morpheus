resource "morpheus_cypher_tfvars" "tf_example_cypher_tfvars" {
  key   = "securetfvars"
  value = <<EOT
account=12345
password=supersecure
EOT
  ttl   = 86400
}