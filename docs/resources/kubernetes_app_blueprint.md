---
page_title: "morpheus_kubernetes_app_blueprint Resource - terraform-provider-morpheus"
subcategory: ""
description: |-
  Provides a Morpheus kubernetes app blueprint resource
---

# morpheus_kubernetes_app_blueprint

Provides a Morpheus kubernetes app blueprint resource

## Example Usage

Creating the Kubernetes app blueprint with local content in yaml format:

```terraform
resource "morpheus_kubernetes_app_blueprint" "tfexample_kubernetes_app_blueprint_yaml" {
  name              = "tf-kubernetes-app-blueprint-example-yaml"
  description       = "tf example kubernetes app blueprint"
  category          = "k8s"
  source_type       = "yaml"
  blueprint_content = <<TFEOF
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
TFEOF
}
```

Creating the Kubernetes app blueprint with Kubernetes spec templates:

```terraform
resource "morpheus_kubernetes_app_blueprint" "tfexample_kubernetes_app_blueprint_spec" {
  name              = "tf-kubernetes-app-blueprint-example-spec"
  description       = "tf example kubernetes app blueprint"
  category          = "k8s"
  source_type       = "spec"
  spec_template_ids = [2, 3]
}
```

Creating the Kubernetes app blueprint with the blueprint fetched via git:

```terraform
resource "morpheus_kubernetes_app_blueprint" "tfexample_kubernetes_app_blueprint_git" {
  name           = "tf-kubernetes-app-blueprint-example-git"
  description    = "tf example kubernetes app blueprint"
  category       = "k8s"
  source_type    = "repository"
  integration_id = 3
  repository_id  = 1
  version_ref    = "main"
  working_path   = "./test"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the kubernetes app blueprint
- `source_type` (String) The source of the kubernetes app blueprint (yaml, spec or repository)

### Optional

- `blueprint_content` (String) The content of the kubernetes app blueprint. Used when the yaml source type is specified
- `category` (String) The category of the kubernetes app blueprint
- `description` (String) The description of the kubernetes app blueprint
- `integration_id` (Number) The ID of the git integration
- `repository_id` (Number) The ID of the git repository
- `spec_template_ids` (List of Number) A list of kubernetes spec template ids associated with the app blueprint
- `version_ref` (String) The git reference of the repository to pull (main, master, etc.)
- `working_path` (String) The path of the kubernetes app blueprint in the git repository

### Read-Only

- `id` (String) The ID of the kubernetes app blueprint

## Import

Import is supported using the following syntax:

```shell
terraform import morpheus_kubernetes_app_blueprint.tf_example_kubernetes_app_blueprint 1
```
