resource "morpheus_group" "tf_example_group" {
  name      = "tfgroup"
  code      = "tfgroup"
  location  = "denver"
  cloud_ids = [1]
}