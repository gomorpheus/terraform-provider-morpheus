resource "morpheus_chef_bootstrap_task" "cheftask" {
  name                = "terraform_example_chef"
  code                = "terraform_example_chef"
  labels              = ["demo", "terraform"]
  server_id           = 9
  environment         = "dev"
  run_list            = "role[web]"
  data_bag_key        = "test123"
  data_bag_key_path   = "/etc/chef/databag_secret"
  node_name           = "demonode"
  node_attributes     = <<EOF
{
  "test":"demo"
}
EOF
  retryable           = true
  retry_count         = 1
  retry_delay_seconds = 10
  allow_custom_config = true
  visibility          = "public"
}