resource "morpheus_manual_option_list" "tf_example_manual_option_list" {
  name        = "demo"
  description = "tetin"
  dataset     = <<POLICY
[{"name": "Level 1","value":"level1"},
 {"name": "Level 2","value":"level2"},
 {"name": "Level 3","value":"level3"}
]
POLICY
  real_time   = true
}