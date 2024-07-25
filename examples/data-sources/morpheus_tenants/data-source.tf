data "morpheus_tenants" "example_tenants" {
  sort_ascending = true
  filter {
    name   = "name"
    values = ["Test*"]
  }
}