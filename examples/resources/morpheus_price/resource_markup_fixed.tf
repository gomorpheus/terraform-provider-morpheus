resource "morpheus_price" "tf_example_price" {
  name          = "terraform-test"
  code          = "terraform-test"
  tenant_id     = 1
  price_type    = "fixed"
  price_unit    = "minute"
  incur_charges = "always"
  currency      = "USD"
  cost          = 38.00
  markup_type   = "fixed"
  markup_cost   = 18.00
}
