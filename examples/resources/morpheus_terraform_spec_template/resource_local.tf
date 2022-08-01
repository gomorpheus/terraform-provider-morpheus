resource "morpheus_terraform_spec_template" "tfexample_terraform_spec_template_local" {
  name         = "tf-terraform-spec-example-local"
  source_type  = "local"
  spec_content = <<EOF
resource "aws_instance" "instance_1" {
  ami           = "ami-0b91a410940e82c54"
  instance_type = "t2.micro"
}
EOF
}