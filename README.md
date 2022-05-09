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
| [morpheus_ansible_playbook_task](docs/resources/ansible_playbook_task.md) | Morpheus ansible playbook automation task resource |
| [morpheus_checkbox_option_type](docs/resources/checkbox_option_type.md) | Morpheus checkbox option type resource |
| [morpheus_contact](docs/resources/morpheus_contact.md) | Morpheus contact resource |
| [morpheus_environment](docs/resources/environment.md) | Morpheus environment resource |
| [morpheus_execute_schedule](docs/resources/execute_schedule.md) | Morpheus execute schedule resource |
| [morpheus_group](docs/resources/group.md) | Morpheus group resource |
| [morpheus_hidden_option_type](docs/resources/hidden_option_type.md) | Morpheus hidden option type resource |
| [morpheus_manual_option_list](docs/resources/manual_option_list.md) | Morpheus manual option list resource |
| [morpheus_max_cores_policy](docs/resources/max_cores_policy.md) | Morpheus max cores policy resource |
| [morpheus_max_hosts_policy](docs/resources/max_hosts_policy.md) | Morpheus max hosts policy resource |
| [morpheus_max_vms_policy](docs/resources/max_vms_policy.md) | Morpheus max vms policy resource |
| [morpheus_network_domain](docs/resources/network_domain.md) | Morpheus network domain resource |
| [morpheus_number_option_type](docs/resources/number_option_type.md) | Morpheus number option type resource |
| [morpheus_operational_workflow](docs/resources/operational_workflow.md) | Morpheus operational automation workflow resource |
| [morpheus_password_option_type](docs/resources/password_option_type.md) | Morpheus password option type resource |
| [morpheus_price](docs/resources/price.md) | Morpheus price resource |
| [morpheus_price_set](docs/resources/price_set.md) | Morpheus price set resource |
| [morpheus_provisiong_workflow](docs/resources/provisioning_workflow.md) | Morpheus provisioning automation workflow resource |
| [morpheus_python_script_task](docs/resources/python_script_task.md) | Morpheus python script automation task resource |
| [morpheus_rest_option_list](docs/resources/rest_option_list.md) | Morpheus REST API option list resource |
| [morpheus_select_list_option_type](docs/resources/select_list_option_type.md) | Morpheus select list option type resource |
| [morpheus_service_plan](docs/resources/service_plan.md) | Morpheus service plan resource |
| [morpheus_task_job](docs/resources/task_job.md) | Morpheus task job resource for scheduling automation tasks |
| [morpheus_tenant](docs/resources/tenant.md) | Morpheus tenant resource |
| [morpheus_terraform_spec_template](docs/resources/terraform_spec_template.md) | Morpheus Terraform spec template resource |
| [morpheus_text_option_type](docs/resources/text_option_type.md) | Morpheus text option type resource |
| [morpheus_typeahead_option_type](docs/resources/typeahead_option_type.md) | Morpheus typeahead option type resource |
| [morpheus_vsphere_cloud](docs/resources/vsphere_cloud.md) | Morpheus VMware vSphere cloud resource |
| [morpheus_vsphere_instance](docs/resources/vsphere_instance.md) | Morpheus VMware vSphere instance resource |
| [morpheus_workflow_policy](docs/resources/workflow_policy.md) | Morpheus workflow policy resource for assigning a workflow to a group, cloud, role, user or globally |

## Supported Data Sources
----------------------

The following list of data sources are supported by the Morpheus Terraform provider:

| Data Source Name | Description |
|------------------|-------------|
| [morpheus_cloud](docs/data-sources/cloud.md) | Morpheus cloud data source |
| [morpheus_contact](docs/data-sources/contact.md) | Morpheus contact data soure |
| [morpheus_environment](docs/data-sources/environment.md) | Morpheus environment data source|
| [morpheus_execute_schedule](docs/data-sources/execute_schedule.md) | Morpheus execute schedule data source |
| [morpheus_group](docs/data-sources/group.md) | Morpheus group data source |
| [morpheus_instance_layout](docs/data-sources/instance_layout.md) | Morpheus isntance layout data source |
| [morpheus_instance_type](docs/data-sources/instance_type.md) | Morpheus instance type data source |
| [morpheus_network](docs/data-sources/network.md) | Morpheus network data source |
| [morpheus_option_type](docs/data-sources/option_type.md) | Morpheus option type data source |
| [morpheus_plan](docs/data-sources/plan.md) | Morpheus plan data source |
| [morpheus_price](docs/data-sources/price.md) | Morpheus price data source |
| [morpheus_price_set](docs/data-sources/price_set.md) | Morpheus price set data source |
| [morpheus_resource_pool](docs/data-sources/resource_pool.md) | Morpheus resources pool data source |
| [morpheus_task](docs/data-sources/task.md) | Morpheus automation task data source |
| [morpheus_workflow](docs/data-sources/workflow.md) | Morpheus workflow data soure |

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
