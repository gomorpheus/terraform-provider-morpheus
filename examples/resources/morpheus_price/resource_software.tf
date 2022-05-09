resource "morpheus_price" "tf_example_price" {
  name          = "terraform-test"
  code          = "terraform-test"
  tenant_id     = 1
  price_type    = "software"
  software      = "visio"
  price_unit    = "minute"
  incur_charges = "always"
  currency      = "USD"
  cost          = 38.00
}
