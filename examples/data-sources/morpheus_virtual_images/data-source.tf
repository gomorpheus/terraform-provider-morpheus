data "morpheus_virtual_images" "example_virtual_images" {
  sort_ascending = true
  source         = "Synced"
  filter {
    name   = "name"
    values = ["Test*"]
  }

  filter {
    name   = "type"
    values = ["vmdk", "iso"]
  }
}