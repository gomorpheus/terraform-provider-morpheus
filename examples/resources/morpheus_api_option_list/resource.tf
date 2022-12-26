resource "morpheus_api_option_list" "tf_example_api_option_list" {
  name               = "tf_example_api_option_list"
  description        = "Terraform Morpheus API option list example"
  visibility         = "private"
  option_list        = "instances"
  translation_script = <<SCRIPT
  var i=0;
  results = [];
  for(i; i<data.length; i++) {
    results.push({name: data[i].name, value: data[i].name});
  }
  SCRIPT
}