// Pin the version
terraform {
  required_providers {
    morpheus = {
      source  = "morpheus/morpheus"
      version = "~> 0.1"
    }
  }
}

// Configure the provider
provider "morpheus" {
  url      = "${var.morpheus_url}"
  username = "${var.morpheus_username}"
  password = "${var.morpheus_password}"       
}