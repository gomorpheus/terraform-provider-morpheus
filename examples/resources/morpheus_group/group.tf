resource "morpheus_group" "terraform_group" {
  name      = "tfgroup"
  code      = "tfgroup"
  location  = "denver"
  cloud_ids = [1]
}