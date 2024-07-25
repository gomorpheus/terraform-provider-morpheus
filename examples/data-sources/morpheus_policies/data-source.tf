data "morpheus_policies" "example_policies" {
  sort_ascending = true

  filter {
    name   = "name"
    values = ["Test*"]
  }

  filter {
    name   = "type"
    values = ["Max VMs", "Workflow"]
  }
}