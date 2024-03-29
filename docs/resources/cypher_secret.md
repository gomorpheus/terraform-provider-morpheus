---
page_title: "morpheus_cypher_secret Resource - terraform-provider-morpheus"
subcategory: ""
description: |-
  Provides a Morpheus cypher secret resource.
---

# morpheus_cypher_secret

Provides a Morpheus cypher secret resource.

## Example Usage

```terraform
resource "morpheus_cypher_secret" "tf_example_cypher_secret" {
  key   = "apipassword"
  value = "password123"
  ttl   = 86400
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `key` (String) The path of the cypher secret, excluding the secret prefix
- `value` (String, Sensitive) The value of the cypher secret

### Optional

- `ttl` (Number) The time to live of the cypher secret

### Read-Only

- `id` (String) The ID of the cypher secret

## Import

Import is supported using the following syntax:

```shell
terraform import morpheus_cypher_secret.tf_example_cypher_secret 1
```
