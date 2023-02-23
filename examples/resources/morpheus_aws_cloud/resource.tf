resource "morpheus_aws_cloud" "tf_example_aws_cloud" {
  name                       = "tf-aws-demo"
  code                       = "tf-aws-demo"
  location                   = "colorado"
  visibility                 = "public"
  tenant_id                  = 1
  enabled                    = true
  automatically_power_on_vms = true
  region                     = "us-east-1"
  access_key                 = "ADMEI422IMWIF824"
  secret_key                 = "34MPW23DQQFEWNGN112WEG"
  inventory                  = "full"
  vpc                        = "all"
  appliance_url              = "https://morpheus.local"
  time_zone                  = "America/Denver"
  ebs_encryption             = true
  datacenter_id              = "tfawsdemo"
  guidance                   = "manual"
  costing                    = "full"
  agent_install_mode         = "cloudInit"
}