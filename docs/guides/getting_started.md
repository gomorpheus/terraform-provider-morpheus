---
subcategory: ""
page_title: "Getting started with the Morpheus Terraform provider"
description: |-
    A guide to getting started with the Morpheus Terraform provider.
---

# Getting Started with the Morpheus Provider

This guide walks you through getting started with using the Morpheus Terraform provider to configure and manage your Morpheus platform.

## Before you begin

[Install Terraform](https://www.terraform.io/intro/getting-started/install.html)
and read the Terraform getting started guide that follows. This guide will
assume basic proficiency with Terraform.

## Configuring the provider

Configure the Terraform provider by specifying the provider information according the [Terraform 0.13-compliant](https://www.terraform.io/upgrade-guides/0-13.html#in-house-providers) provider installation standard.
Create a `provider.tf` file with the following content to define the provider configuration.

```terraform
terraform {
  required_providers {
    morpheus = {
      source  = "gomorpheus/morpheus"
      version = "0.11.0"
    }
  }
}

provider "morpheus" {
  url      = "https://morpheus.test.local"
  username = "administrator"
  password = "password"
}
```

The provider also supports the use of an [access token](auth.md#access-token) instead of specifying a username and password to authentication to the Morpheus platform. 

## Creating your first Morpheus resource
Once the provider is configured, you can apply the Morpheus resources defined in your Terraform file. The following is an example Terraform file containing a Morpheus environment resource. Create a `main.tf` file with the following content to define the environment resource.

```terraform
resource "morpheus_environment" "tfdemo" {
  active      = true
  code        = "tfdemo"
  description = "Terraform provider demo environment"
  name        = "TFDemo"
}
```

Use `terraform init` to initialize the specified version of the Morpheus provider:

```
$ terraform init

Initializing the backend...

Initializing provider plugins...
- Finding morpheusdata.com/gomorpheus/morpheus versions matching "0.11.0"...
- Installing morpheusdata.com/gomorpheus/morpheus v0.11.0...
- Installed morpheusdata.com/gomorpheus/morpheus v0.11.0 (unauthenticated)

Terraform has created a lock file .terraform.lock.hcl to record the provider
selections it made above. Include this file in your version control repository
so that Terraform can guarantee to make the same selections by default when
you run "terraform init" in the future.

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
```

Next, use `terraform plan` to display a list of resources to be created, and highlight any possible unknown attributes at apply time.

```
terraform plan

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # morpheus_environment.tfdemo will be created
  + resource "morpheus_environment" "tfdemo" {
      + active      = true
      + code        = "tfdemo"
      + description = "Terraform provider demo environment"
      + id          = (known after apply)
      + name        = "TFDemo"
      + visibility  = "private"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

Use `terraform apply` to create the resource shown above.


```
terraform apply --auto-approve

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # morpheus_environment.tfdemo will be created
  + resource "morpheus_environment" "tfdemo" {
      + active      = true
      + code        = "tfdemo"
      + description = "Terraform provider demo environment"
      + id          = (known after apply)
      + name        = "TFDemo"
      + visibility  = "private"
    }

Plan: 1 to add, 0 to change, 0 to destroy.
morpheus_environment.tfdemo: Creating...
morpheus_environment.tfdemo: Creation complete after 1s [id=5]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

Congratulations! You've successfully created your first Morpheus resource using the Terraform provider.