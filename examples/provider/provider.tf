terraform {
  required_providers {
    morpheus = {
      source  = "gomorpheus/morpheus"
      version = "0.9.2"
    }
  }
}

provider "morpheus" {
  url      = var.morpheus_url
  username = var.morpheus_username
  password = var.morpheus_password
}