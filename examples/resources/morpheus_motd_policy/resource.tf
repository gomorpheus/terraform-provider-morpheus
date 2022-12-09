resource "morpheus_motd_policy" "tf_example_motd_policy" {
  name         = "tf_example_motd_policy"
  description  = "terraform example global user creation policy"
  enabled      = true
  title        = "TF Example MOTD"
  message      = "This is a test message of the day message"
  message_type = "info"
  full_page    = true
}