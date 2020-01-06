# Terraform Provider for Morpheus

- Website: https://www.morpheusdata.com/
- Docs: [Morpheus Documentation](https://docs.morpheusdata.com)
- Support: [Morpheus Support](https://support.morpheusdata.com)

<img src="https://www.morpheusdata.com/wp-content/uploads/2018/06/cropped-morpheus_highres.png" width="600px">

This is the Terraform provider for the Morpheus data appliance. It interfaces with the [Morpheus API](https://bertramdev.github.io/morpheus-apidoc/) using the morpheus-go-sdk client. Like all [Terraform Providers](https://github.com/terraform-providers/), it is written in Go.

This is being developed in conjunction with [morpheus-go-sdk](https://github.com/gomorpheus/morpheus-go-sdk).  

**BETA** This library is actively under development and is only available as a prototype. A fully featured version will be available in the near future.

## Requirements

* [Terraform](https://www.terraform.io/) | 0.12+
* [Go](https://golang.org/dl/) | 1.13

## Using the Provider

This is an example of a terraform configuration that will create a cloud and group using the `morpheus` provider.

Create a file named `main.tf` under the `examples` directory.

```sh
cd examples
touch main.tf
```

#### main.tf

```
provider "morpheus" {
  url          = "https://api.gomorpheus.com"
  access_token = "a3a4c6ea-fb54-42af-109b-63bdd19e5ae1"
}

resource "morpheus_cloud" "example" {
  name = "tftest"
  type = "vmware"
  code = "vmware"
  location = "US East"
  description = "A VMware vCenter cloud created with Terraform."
  config = {
    apiUrl = "https://10.0.0.5/sdk"
    username = "administrator@yourlabs.com"
    password = "vcenterpassword"
    datacenter = "labs-dc"
    cluster = "QA-vSAN"
  }
}

resource "morpheus_group" "example" {
  name         = "tftest"
  location     = "Test Bunker 2"
  clouds = [morpheus_cloud.example.name]
}

resource "morpheus_network_domain" "example" {
  name         = "terraform.gomorpheus.com"
  description  = "A test domain record"
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

For more information on configuring morpheus resources, visit the [Provider Wiki](/gomorpheus/terraform-provider-morpheus/wiki).

#### Resource Examples

Example | Description
--------- | -----------
[example_instance.tf.json](/gomorpheus/terraform-provider-morpheus/blob/master/examples/example_instance.tf.json) | Basic instance resource example.


### Testing the provider

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


### Building the provider

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

## Developing the provider

Please help contribute to the Terraform provider for Morpheus.

First, you'll need to install Go and Terraform, see [Requirements](#requirements).

Currently, we are in the process of building out all the available morpheus resources.

### Developing the SDK

While working on the provider, you may also be working on the [morpheus-go-sdk](https://github.com/gomorpheus/morpheus-go-sdk), which can be found at `$GOPATH/src/github.com/gomorpheus/morpheus-go-sdk`.


### External Resources

- [Morpheus API](https://bertramdev.github.io/morpheus-apidoc/)
- [Writing Custom Providers](https://www.terraform.io/docs/extend/writing-custom-providers.html)
- [Terraform Plugin SDK](https://github.com/hashicorp/terraform-plugin-sdk)

