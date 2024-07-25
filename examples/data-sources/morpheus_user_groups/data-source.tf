data "morpheus_user_groups" "example_user_groups" {
  sort_ascending = true

  filter {
    name   = "name"
    values = ["Test*"]
  }
}