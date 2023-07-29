resource "morpheus_servicenow_integration" "tf_example_servicenow_integration" {
  name                = "terraform servicenow integration"
  enabled             = true
  url                 = "https://servicenowprod.service-now.com"
  username            = "my-snow-username"
  password            = "my-snow-password"
  cmdb_custom_mapping = <<EOF
{
"object_id":"<%=instance.name%>",
"SN_field_id2":"<%=morph.varname2%>",
"SN_field_id3":"<%=morph.varname3%>"
}
  EOF
  cmdb_class_mapping = {
    "Amazon Instance"                = "cmdb_ci_ec2_instance"
    "Hyper-V Hypervisor - Unmanaged" = "cmdb_ci_hyper_v_instance"
    "VMware Windows VM"              = "cmdb_ci_vmware_instance"
  }
  default_cmdb_business_class = "demo"
}