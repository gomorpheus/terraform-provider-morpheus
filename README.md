# Terraform Provider for Morpheus

[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/gomorpheus/terraform-provider-morpheus?label=release)](https://github.com/gomorpheus/terraform-provider-morpheus/releases) [![license](https://img.shields.io/github/license/gomorpheus/terraform-provider-morpheus.svg)]()

<img src="https://morpheusdata.com/wp-content/uploads/2020/04/morpheus-logo-v2.svg" width="300px">

- Website: https://www.morpheusdata.com/
- Docs: [Morpheus Documentation](https://docs.morpheusdata.com)
- Support: [Morpheus Support](https://support.morpheusdata.com)


This is the Terraform provider for the Morpheus Data Cloud Management Platform (CMP). It interfaces with the [Morpheus API](https://apidocs.morpheusdata.com/) using the morpheus-go-sdk client. Like all [Terraform Providers](https://github.com/terraform-providers/), it is written in Go.

This is being developed in conjunction with [morpheus-go-sdk](https://github.com/gomorpheus/morpheus-go-sdk).  

## Requirements
------------

* [Terraform](https://www.terraform.io/) | 0.13+
* [Go](https://golang.org/dl/) 1.16 (to build the provider plugin)


## Getting Started
---------------------

The best way to get started using the Morpheus Terraform provider is by following the [getting started guide](docs/guides/getting_started.md).

## Supported Resources
----------------------

The following list of resources are supported by the Morpheus Terraform provider:

| Resource Name | Description |
|------|---------------|
| morpheus_ansible_playbook_task | Morpheus ansible playbook automation task resource |
| morpheus_checkbox_option_type | Morpheus checkbox option type resource |
| morpheus_environment | Morpheus environment resource |
| morpheus_execute_schedule | Morpheus execute schedule resource |
| morpheus_group | Morpheus group resource |
| morpheus_hidden_option_type | Morpheus hidden option type resource |
| morpheus_manual_option_list | Morpheus manual option list resource |
| morpheus_max_cores_policy | Morpheus max cores policy resource |
| morpheus_max_hosts_policy | Morpheus max hosts policy resource |
| morpheus_max_vms_policy | Morpheus max vms policy resource |
| morpheus_network_domain | Morpheus network domain resource |
| morpheus_number_option_type | Morpheus number option type resource |
| morpheus_operational_workflow | Morpheus operational automation workflow resource |
| morpheus_password_option_type | Morpheus password option type resource |
| morpheus_provisioning_workflow | Morpheus provisioning automation workflow resource |
| morpheus_python_script_task | Morpheus python script automation task resource |
| morpheus_rest_option_list | Morpheus REST API option list resource |
| morpheus_select_list_option_type | Morpheus select list option type resource |
| morpheus_tenant | Morpheus tenant resource |
| morpheus_terraform_spec_template | Morpheus Terraform spec template resource |
| morpheus_text_option_type | Morpheus text option type resource |
| morpheus_typeahead_option_type | Morpheus typeahead option type resource |
| morpheus_vsphere_cloud | Morpheus VMware vSphere cloud resource |
| morpheus_vsphere_instance | Morpheus VMware vSphere instance resource |
| morpheus_workflow_policy | Morpheus workflow policy resource for assigning a workflow to a group, cloud, role, user or globally |

## Supported Data Sources
----------------------

The following list of data sources are supported by the Morpheus Terraform provider:

| Data Source Name | Description |
|------------------|-------------|
| morpheus_cloud | Morpheus cloud data source |
| morpheus_environment | Morpheus environment data source|
| morpheus_execute_schedule | Morpheus execute schedule data source |
| morpheus_group | Morpheus group data source |
| morpheus_instance_layout | Morpheus isntance layout data source |
| morpheus_instance_type | Morpheus instance type data source |
| morpheus_network | Morpheus network data source |
| morpheus_option_type | Morpheus option type data source |
| morpheus_plan | Morpheus plan data source |
| morpheus_resource_pool | Morpheus resources pool data source |
| morpheus_task | Morpheus automation task data source |
| morpheus_workflow | Morpheus workflow data soure |

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

## Generating Docs
----------------------
From the root of the repo run:

```
go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
```

## Developing the provider
-------------------------

See the [`contributing`](contributing/) directory for more developer documentation.
