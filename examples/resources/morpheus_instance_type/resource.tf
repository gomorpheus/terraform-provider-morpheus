resource "morpheus_instance_type" "tf_example_instance_type" {
  name               = "tf_example_instance"
  code               = "tf_example_instance"
  description        = "Terraform Example Instance Type"
  labels             = ["demo", "instance", "terraform"]
  category           = "web"
  visibility         = "private"
  image_path         = "tfexample.png"
  image_name         = "tfexample.png"
  featured           = false
  enable_deployments = true
  enable_scaling     = true
  enable_settings    = true
  environment_prefix = "TFEXAMPLE_DEMO"
  option_type_ids    = [1910, 1912]

  evar {
    name   = "first"
    value  = "first"
    export = true
  }

  evar {
    name         = "second"
    masked_value = "second"
    export       = false
  }
}