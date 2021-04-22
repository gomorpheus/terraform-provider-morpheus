# Terraform Provider for Morpheus

- Website: https://www.morpheusdata.com/
- Docs: [Morpheus Documentation](https://docs.morpheusdata.com)
- Support: [Morpheus Support](https://support.morpheusdata.com)

<img src="https://morpheusdata.com/wp-content/uploads/2020/04/morpheus-logo-v2.svg" width="200px">

This is the Terraform provider for the Morpheus data appliance. It interfaces with the [Morpheus API](https://apidocs.morpheusdata.com/) using the morpheus-go-sdk client. Like all [Terraform Providers](https://github.com/terraform-providers/), it is written in Go.

This is being developed in conjunction with [morpheus-go-sdk](https://github.com/gomorpheus/morpheus-go-sdk).  

**BETA** This library is actively under development and is only available as a prototype. A fully featured version will be available in the near future.

## Requirements
------------

* [Terraform](https://www.terraform.io/) | 0.12+
* [Go](https://golang.org/dl/) 1.14 (to build the provider plugin)

## Building the provider
-------------------------

Clone repository to: `$GOPATH/src/github.com/gomorpheus/terraform-provider-morpheus`

```sh
mkdir -p $GOPATH/src/github.com/gomorpheus; cd $GOPATH/src/github.com/gomorpheus
git clone git@github.com:gomorpheus/terraform-provider-morpheus
```

As an alternative to cloning manually, you can use `go get`:

```sh
go get -v github.com/gomorpheus/terraform-provider-morpheus/...
```

Enter the provider directory and build the provider.

```sh
cd $GOPATH/src/github.com/gomorpheus/terraform-provider-morpheus
go build -o terraform-provider-morpheus
```

## Using the Provider
---------------------

This is an example of a terraform configuration that will create a cloud and group using the `morpheus` provider.

## Testing the provider
------------------------
If you are actively developing the provider, always remember to [build](#Building-the-provider) first in order to test your changes.

Use terraform to create resources.

```bash
terraform init examples && terraform plan examples && terraform apply examples
```

Use `[morpheus-cli](/gomorpheus/morpheus-cli)` to see that the resources were created.

```bash
morpheus groups list -s tftest
morpheus network-domains list -s tftest
```

Use terraform to destroy resources.

```bash
terraform destroy
```

<!-- 
### Installing the provider
To use a released provider in your Terraform environment, run [`terraform init`](https://www.terraform.io/docs/commands/init.html) and Terraform will automatically install the provider. To specify a particular provider version when installing released providers, see the [Terraform documentation on provider versioning](https://www.terraform.io/docs/configuration/providers.html#version-provider-versions).

To instead use a custom-built provider in your Terraform environment (e.g. the provider binary from the build instructions below), follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory,  run `terraform init` to initialize it. -->

## Developing the provider
-------------------------
If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.13+ is required). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to our `$PATH`.

### Developing the SDK

While working on the provider, you may also be working on the [morpheus-go-sdk](https://github.com/gomorpheus/morpheus-go-sdk), which can be found at `$GOPATH/src/github.com/gomorpheus/morpheus-go-sdk`.