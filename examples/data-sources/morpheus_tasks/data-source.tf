data "morpheus_tasks" "example_tasks" {
  sort_ascending = true

  filter {
    name   = "name"
    values = ["Test*"]
  }

  filter {
    name   = "type"
    values = ["Shell Script","Python Script"]
  }
}