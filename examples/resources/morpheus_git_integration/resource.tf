data "morpheus_key_pair" "tf_example_key_pair" {
  name = "morpheusgit"
}

resource "morpheus_git_integration" "tf_example_git_integration" {
  name               = "tftest"
  enabled            = true
  url                = "https://github.com/gomorpheus/tfdemo.git"
  default_branch     = "main"
  key_pair_id        = data.morpheus_key_pair.tf_example_key_pair.id
  enable_git_caching = true
}