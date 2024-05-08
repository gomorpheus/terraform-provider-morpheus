resource "morpheus_app_blueprint_catalog_item" "tf_example_app_blueprint_catalog_item" {
  name                 = "tfexample_app_blueprint_catalog"
  description          = "terraform example app blueprint catalog item"
  logo_image_path      = "tfexample.png"
  logo_image_name      = "tfexample.png"
  dark_logo_image_path = "tfexampledark.png"
  dark_logo_image_name = "tfexampledark.png"
  enabled              = true
  featured             = true
  labels               = ["aws", "demo", "testing"]
  content              = file("${path.module}/catalog-data.md")
  visibility           = "public"
  blueprint_id         = 5
  option_type_ids      = [2056, 2006, 2058]
  app_spec             = file("${path.module}/appSpec.yaml")
}
