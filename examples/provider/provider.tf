terraform {
  required_providers {
    morpheus = {
      source  = "morpheusdata.com/gomorpheus/morpheus"
      version = "0.3.1"
    }
  }
}

provider "morpheus" {
  url      = var.morpheus_url
  username = var.morpheus_username
  password = var.morpheus_password
}