resource "morpheus_file_template" "tfexample_file_template" {
  name             = "tf-terraform-file-template"
  labels           = ["demo", "template", "terraform"]
  file_name        = "tfcustom.cnf"
  file_path        = "/etc/my.cnf.d"
  phase            = "preProvision"
  file_content     = file("${path.module}/custom.cnf")
  file_owner       = "root"
  setting_name     = "myCnf"
  setting_category = "master"
}