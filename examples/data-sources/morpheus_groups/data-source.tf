data "morpheus_groups" "terraform_test" {
  sort_ascending = false
  filter {
    name   = "name"
    values = ["^tf*"]
  }

  filter {
    name   = "location"
    values = ["denver"]
  }
}