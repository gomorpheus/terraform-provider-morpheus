resource "morpheus_boot_script" "tf_example_boot_script" {
  name    = "TF Example Boot Script"
  content = "ls"
}