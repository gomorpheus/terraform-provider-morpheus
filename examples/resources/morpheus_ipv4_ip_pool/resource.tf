resource "morpheus_ipv4_ip_pool" "tf_example_ipv4_pool" {
  name = "Terraform Example IPv4 IP pool"
  ip_range {
    starting_address = "192.168.1.1"
    ending_address   = "192.168.1.10"
  }
  ip_range {
    starting_address = "10.0.0.1"
    ending_address   = "10.0.0.10"
  }
}