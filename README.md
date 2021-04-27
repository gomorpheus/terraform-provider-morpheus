# Terraform Provider for Morpheus

<img src="https://morpheusdata.com/wp-content/uploads/2020/04/morpheus-logo-v2.svg" width="300px">

- Website: https://www.morpheusdata.com/
- Docs: [Morpheus Documentation](https://docs.morpheusdata.com)
- Support: [Morpheus Support](https://support.morpheusdata.com)


This is the Terraform provider for the Morpheus Data Cloud Management Platform (CMP). It interfaces with the [Morpheus API](https://apidocs.morpheusdata.com/) using the morpheus-go-sdk client. Like all [Terraform Providers](https://github.com/terraform-providers/), it is written in Go.

This is being developed in conjunction with [morpheus-go-sdk](https://github.com/gomorpheus/morpheus-go-sdk).  

**BETA** This library is actively under development and is only available as a prototype. A fully featured version will be available in the near future.

## Requirements
------------

* [Terraform](https://www.terraform.io/) | 0.13+
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

Enter the provider directory.

```sh
cd $GOPATH/src/github.com/gomorpheus/terraform-provider-morpheus
```

Build the provider using `make dev`. This will place the provider onto your system in a [Terraform 0.13-compliant](https://www.terraform.io/upgrade-guides/0-13.html#in-house-providers) manner.

```bash
make dev
```

You'll need to ensure that your Terraform file contains the information necessary to find the plugin when running `terraform init`. `make dev` will use a version number of 0.0.1, so the following block will work:

```hcl
terraform {
  required_providers {
    morpheus = {
      source = "localhost/providers/morpheus"
      version = "0.0.1"
    }
  }
}
```

## Using the Provider

---------------------

When the provider is out of beta the documentation will be available alongside the Terraform provider on the Terraform registry, but during the beta phase the best resource is this [guide](docs/guides/getting_started.md).

## Developing the provider
-------------------------

See the [`contributing`](contributing/) directory for more developer documentation.