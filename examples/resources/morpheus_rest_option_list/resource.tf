resource "morpheus_rest_option_list" "restoptionlistdemo" {
  name               = "tfrestexample"
  description        = "tetin"
  visibility         = "private"
  source_url         = "https://api.github.com/repos/hashicorp/consul/releases"
  real_time          = true
  ignore_ssl_errors  = true
  source_method      = "GET"
  initial_dataset    = <<POLICY
  [{"name": "Level 1","value":"level1"},
  {"name": "Level 2","value":"level2"},
  {"name": "Level 3","value":"level3"}
  ]
  POLICY
  translation_script = <<POLICY
      for(var x=0;x < 5; x++) {
          results.push({name: data[x].name,value:data[x].name});
        }
  POLICY
}