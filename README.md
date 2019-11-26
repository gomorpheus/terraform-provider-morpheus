# Terraform Provider for Morpheus

- Website: https://www.morpheusdata.com/
- Docs: [Morpheus Documentation](https://docs.morpheusdata.com)
- Support: [Morpheus Support](https://support.morpheusdata.com)

<img src="https://www.morpheusdata.com/wp-content/uploads/2018/06/cropped-morpheus_highres.png" width="600px">

This is the Terraform provider for the Morpheus data appliance. It interfaces with the [Morpheus API](https://bertramdev.github.io/morpheus-apidoc/) using the morpheusapi client. Like all [Terraform Providers](https://github.com/terraform-providers/), it is written in Go.

This is being developed in conjunction with [morpheusapi](https://github.com/gomorpheus/morpheus-go/morpheusapi).  

**BETA** This library is actively under development and is only available as a prototype, **version 0.1**. A fully featured version will be available in the near future.

## Requirements

* [Terraform](https://www.terraform.io/) | 0.12+
* [Go](https://golang.org/dl/) | 1.13

## Using the Provider

This is an example of a terraform configuration that will create a cloud and group using the `morpheus` provider.

### Example Config

```
provider "morpheus" {
  url          = "https://yourmorpheus.com"
  access_token = "a3a4c6fa-fb54-14af-a09b-13bdd19e5ae5"
}

resource "morpheus_cloud" "tftest_cloud" {
  name = "tftest"
  type = "vmware"
  code = "vmware"
  location = "US East"
  description = "A VMware vCenter cloud created with Terraform."
  config = {
    apiUrl = "https://10.0.0.150/sdk"
    username = "administrator@yourcompany.com"
    password = "b24n32jh4g98"
    datacenter = "labs-denver"
    cluster = "Test"
  }
}

resource "morpheus_group" "tftest_group" {
  name         = "tftest"
  location     = "Test Bunker"
  clouds = [morpheus_cloud.tftest_cloud.name]
}

```

#### Provider Settings

The `morpheus` provider has the following configuration options.

Name | Description
--------- | -----------
url | The URL of your Morpheus appliance. eg. https://yourmorpheus.com
access_token | A valid Morpheus API access token. eg. "d3a4c6fa-fb54-44af"
username | Morpheus username.
password | Morpheus password.

There are 2 different ways to authenticate.

* `access_token`
* `username` and `password`

Be sure to utilize [variables](#https://learn.hashicorp.com/terraform/getting-started/variables.html) to set secret values like `access_token` and `password` in your configuration.

For more information on configuring morpheus resources, visit the [Provider Wiki](/gomorpheus/terraform-provider-morpheus/wiki/CLI-Manual).

### Testing the provider

If you are working on the privder, always remember to [build](#Building the provider) first.

Use terraform to create resources.

```bash
terraform init && terraform plan && terraform apply
```

Use `[morpheus](https://github.com/gomorpheus/morpheus-cli)` to see that the resources were created.

```bash
morpheus groups list -s tftest
morpheus network-domains list -s tftest
```

Use terraform to destroy resources.

```bash
terraform destroy
```

<!-- 
### Installing the plugin
To use a released provider in your Terraform environment, run [`terraform init`](https://www.terraform.io/docs/commands/init.html) and Terraform will automatically install the provider. To specify a particular provider version when installing released providers, see the [Terraform documentation on provider versioning](https://www.terraform.io/docs/configuration/providers.html#version-provider-versions).

To instead use a custom-built provider in your Terraform environment (e.g. the provider binary from the build instructions below), follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory,  run `terraform init` to initialize it. -->


### Building the provider

Our `Makefile` is under construction. You can use the following steps to build the provider.

First install dependencies.

```bash
go get -v github.com/gomorpheus/morpheus-go/morpheusapi/...
```

<!-- Alternatively, you could just use: `cd $GOPATH/src/github.com/gomorpheus && git clone https://github.com/gomorpheus/morpheus-go/morpheusapi.git`. -->

Build the executable using `go build`.

```bash
go build -o terraform-provider-morpheus
```

## Developing the provider

Please help contribute to the Terraform provider for Morpheus.

First, you'll need to install Go and Terraform, see [Requirements](#requirements).

<!--
*Note:* This project uses [Go Modules](https://blog.golang.org/using-go-modules) making it safe to work with it outside of your existing [GOPATH](http://golang.org/doc/code.html#GOPATH). The instructions that follow assume a directory in your home directory outside of the standard GOPATH (i.e `$HOME/development/terraform-providers/`).

Clone repository to: `$HOME/development/terraform-providers/`

```sh
$ mkdir -p $HOME/development/terraform-providers/; cd $HOME/development/terraform-providers/
$ git clone git@github.com:gomorpheus/terraform-provider-morpheus
...
```

Enter the provider directory and run `make tools`. This will install the needed tools for the provider.

```sh
$ make tools
```

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-morpheus
...
```

**WARNING** Makefile is not yet ready. See [Building the Provider](#building-the-provider) to build the provider manually.
-->

### External Resources

- [Morpheus API](https://bertramdev.github.io/morpheus-apidoc/)
- [Writing Custom Providers](https://www.terraform.io/docs/extend/writing-custom-providers.html)
- [Terraform Plugin SDK](https://github.com/hashicorp/terraform-plugin-sdk)

