resource "morpheus_security_package" "tf_example_security_package" {
  name        = "tf_example_security_package"
  description = "Terraform security package example"
  labels      = ["demo", "terraform"]
  enabled     = true
  url         = "https://github.com/ComplianceAsCode/content/releases/download/v0.1.59/scap-security-guide-0.1.59.zip"
}