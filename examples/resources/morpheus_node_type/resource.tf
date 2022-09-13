data "morpheus_virtual_image" "example_virtual_image" {
  name = "Ubuntu 20.04 Template"
}

resource "morpheus_node_type" "tf_example_node" {
  name             = "tf_example_node_type"
  short_name       = "tfexamplenodetype"
  technology       = "vmware"
  version          = "2.0"
  category         = "tfexample"
  virtual_image_id = data.morpheus_virtual_image.example_virtual_image.id

  file_template_ids = [
    data.morpheus_file_template.tfexample.id,
    113
  ]

  script_template_ids = [
    data.morpheus_script_template.tfscript1.id,
    data.morpheus_script_template.tfscript2.id
  ]

  service_port {
    name     = "web"
    port     = "8080"
    protocol = "HTTP"
  }

  service_port {
    name     = "secureweb"
    port     = "8443"
    protocol = "HTTPS"
  }
}