resource "morpheus_key_pair" "tf_example_key_pair" {
  name        = "TF Example Key Pair"
  public_key  = "ssh-rsa AAAAB3Nz"
  private_key = file("privatekey.rsa")
  passphrase  = "12312"
}