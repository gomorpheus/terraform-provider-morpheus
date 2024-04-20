---
subcategory: ""
page_title: "Sharing Data Between Morpheus Terraform Apps/Instances"
description: |-
    A guide to obtaining client credentials and adding them to provider configuration.
---

# Sharing Data Between Morpheus Terraform Apps/Instances

This guide walks you through utilizing Morpheus Cypher for sharing output data between Morpheus Terraform apps and instances.

## Overview

One of the challenges with deploying and managing infrastructure with Terraform
is the need to share information between code that has different lifecycles. An
example of this would be a cloud virtual network (i.e. - VPC, VNet, etc.) and the resources 
that are deployed into the cloud network. In this scenario the code for the resources
that will be deployed into the network will need information about the network such as
a name or an ID for reference purposes.

### What about Terraform remote state?

Terraform has a native remote_state data source that can be used by a downstream
deployment to reference the state and outputs of an upstream deployment. One of the major downsides 
of using the remote_state data source is the need to expose the entire state file to the downstream resources. 
The upstream statefile may contain sensitive information that doesn't need to be accessed by the downstream resource 
but the remote_state data source is an all or nothing situation. HashiCorp actually recommends using other alternatives 
(https://developer.hashicorp.com/terraform/language/state/remote-state-data) where possible to avoid granting full state access. 

## Morpheus Solution Overview

The commonly recommened solution is to utilize provider specific data sources to pull the information directly
from the source (i.e. - cloud) or leverage an external data store. The native Morpheus solution utilizes the platform
built-in secure storage funcitonality, Cypher. This will allow Terraform code to write and read data that should be shared
to Cypher. 

### Requirements

The following Morpheus objects are required to implement the solution:

**Upstream Resource Requirements**

* **Writer Service Account:** A Morpheus service account with permissions to write data to Morpheus cypher
* **Writer User Role:** A Morpheus user role with permission to write to Cypher
* **Writer Cypher Policy (optional):** A Morpheus cypher policy associated with the "writer" service account that restricts what cypher paths the service account has access to write to. 
* **Writer TFVar Secret:** A Morpheus TFVars cypher secret for storing the API details and credentials used by the Morpheus Terraform provider to write the outputs to cypher that should be accessed by a downstream resource or resources.

**Downstream Resource Requirements**

* **Reader Service Account:** A Morpheus service account with permissions to read data from Morpheus cypher
* **Reader User Role:** A Morpheus user role with permission to read from Cypher.
* **Reader Cypher Policy (optional):** A Morpheus cypher policy associated with the "reader" service account that restricts what cypher paths the service account has access to read from. 
* **Reader TFVars Secret:** A Morpheus TFVars cypher secret for storing the API details and credentials used by the Morpheus Terraform provider to allow a downstream resource to read the upstream outputs from cypher.

#### Cypher Policies

Cypher policies in the Morpheus platform are used to granularly restrict access to the 
Cypher secure key/value store. The policies in this case are used to align with the security
principle of least privilege.

**Writer Policy**

The service account used to write the output data to cypher can be further 
constrained from a permissions standpoint by associating a Morpheus cypher policy
with the user account. This could be used in situations where you only want the account
to be able to write to a specific path in the cypher store. This is valuable
when you want to use multiple service accounts for greater segregation.

**Reader Policy**

The service account used to read the output data from cypher should be further
constrained from a permissions standpoint by associating a Morpheus cypher policy
with the user account. This ensures that the Terraform deployment only has access to
read the outputs that you explicitly intend for it to read.

### Terraform Code Examples

The following code is an example of the design in practice:

**Upstream Resource**

The upstream resource in this case creates a few local files and "outputs" the file hash attribute that's generated after the file has been created.
The file's hash is what needs to be made available to the downstream resource. In a real situation, the information being output is a resource ID, name, or other pertinent detail about the resource.

```terraform
terraform {
  required_providers {
    morpheus = {
      source  = "gomorpheus/morpheus"
      version = "0.9.9"
    }
  }
}

provider "morpheus" {
  url      = var.morpheus_url
  username = var.morpheus_username
  password = var.morpheus_password
}

resource "local_file" "foo" {
  content  = "foo!1"
  filename = "${path.module}/foo.bar"
}

resource "local_file" "foo2" {
  content  = "foo!2"
  filename = "${path.module}/foo2.bar"
}

output "file_md5_hash" {
  value = local_file.foo.content_md5
}

resource "morpheus_cypher_secret" "foo_hash_cypher_secret" {
  key   = "exampleworkspace/foo_hash"
  value = local_file.foo.content_md5
}

locals {
  file_hashes = [local_file.foo.content_md5, local_file.foo2.content_md5]
}

resource "morpheus_cypher_secret" "foo_hashes_cypher_secret" {
  key   = "exampleworkspace/foo_hashes"
  value = jsonencode({"dataoutput" : local.file_hashes})
}
```

**Downstream Resource**

The downstream resource in this example reads the hashes of the file created in the upstream resource.
In a real scenario, the information being read from cypher would be a resource ID, name, or other pertinent information.

```terraform
terraform {
  required_providers {
    morpheus = {
      source  = "gomorpheus/morpheus"
      version = "0.9.9"
    }
  }
}

provider "morpheus" {
  url      = var.morpheus_url
  username = var.morpheus_username
  password = var.morpheus_password
}

data "morpheus_cypher_secret" "shared_test" {
  key = "exampleworkspace/foo_hash"
}

data "morpheus_cypher_secret" "shared_hash_test" {
  key = "exampleworkspace/foo_hashes"
}

locals {
  hashes = jsondecode(data.morpheus_cypher_secret.shared_hash_test.value)
}

output "test" {
  value = data.morpheus_cypher_secret.shared_test.value
}

output "hash_test" {
  value = local.hashes.dataoutput
}
```
