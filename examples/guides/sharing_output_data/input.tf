terraform {
  required_providers {
    morpheus = {
      source  = "gomorpheus/morpheus"
      version = "0.10.0"
    }
  }
}

provider "morpheus" {
  url      = var.morpheus_url
  username = var.morpheus_username
  password = var.morpheus_password
}

data "morpheus_cypher_secret" "shared_test" {
  key = "exampleworkspace/foo_hash"
}

data "morpheus_cypher_secret" "shared_hash_test" {
  key = "exampleworkspace/foo_hashes"
}

locals {
  hashes = jsondecode(data.morpheus_cypher_secret.shared_hash_test.value)
}

output "test" {
  value = data.morpheus_cypher_secret.shared_test.value
}

output "hash_test" {
  value = local.hashes.dataoutput
}