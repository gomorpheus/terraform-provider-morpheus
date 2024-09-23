data "morpheus_networks" "tf_example_networks" {
  cloud_id       = 3
  sort_ascending = true
  filter {
    name   = "name"
    values = ["Test*"]
  }
}