terraform {
  required_providers {
    morpheus = {
      source  = "gomorpheus/morpheus"
      version = "0.12.0"
    }
  }
}

provider "morpheus" {
  url      = "https://morpheus.test.local"
  username = "administrator"
  password = "password"
}