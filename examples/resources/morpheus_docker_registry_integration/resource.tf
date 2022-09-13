resource "morpheus_docker_registry_integration" "tf_example_docker_registry_integration" {
  name     = "tfexampledockerregistry"
  enabled  = true
  url      = "https://index.docker.io/v1/"
  username = "admin"
  password = "password123"
}