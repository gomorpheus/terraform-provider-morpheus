terraform {
  required_providers {
    morpheus = {
      source  = "morpheusdata.com/gomorpheus/morpheus"
      version = "0.3.1"
    }
  }
}

provider "morpheus" {
  url      = "https://morpheus.test.local"
  username = "administrator"
  password = "password"
}