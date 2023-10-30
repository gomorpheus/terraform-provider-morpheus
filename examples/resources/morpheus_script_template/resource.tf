resource "morpheus_script_template" "tfexample_script_template" {
  name           = "tf-terraform-script-template"
  labels         = ["demo", "template", "terraform"]
  script_type    = "bash"
  script_phase   = "provision"
  script_content = <<EOF
echo "testing"
EOF
  run_as_user    = "root"
  sudo           = true
}