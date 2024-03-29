---
page_title: "morpheus_cypher_tfvars Resource - terraform-provider-morpheus"
subcategory: ""
description: |-
  Provides a Morpheus cypher tfvars secret resource.
---

# morpheus_cypher_tfvars

Provides a Morpheus cypher tfvars secret resource.

## Example Usage

```terraform
resource "morpheus_cypher_tfvars" "tf_example_cypher_tfvars" {
  key   = "securetfvars"
  value = <<EOT
account=12345
password=supersecure
EOT
  ttl   = 86400
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `key` (String) The path of the cypher tfvars secret, excluding the secret prefix
- `value` (String, Sensitive) The value of the cypher tfvars secret

### Optional

- `ttl` (Number) The time to live of the cypher tfvars secret

### Read-Only

- `id` (String) The ID of the cypher tfvars secret

## Import

Import is supported using the following syntax:

```shell
terraform import morpheus_cypher_tfvars.tf_example_cypher_tfvars 1
```
