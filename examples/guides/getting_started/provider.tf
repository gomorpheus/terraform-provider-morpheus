terraform {
  required_providers {
    morpheus = {
      source  = "gomorpheus/morpheus"
      version = "0.11.0"
    }
  }
}

provider "morpheus" {
  url      = "https://morpheus.test.local"
  username = "administrator"
  password = "password"
}