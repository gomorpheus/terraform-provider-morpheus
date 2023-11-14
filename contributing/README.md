# Contributing to the Morpheus Terraform Provider

This directory contains documentation about the Morpheus Terraform Provider codebase, aimed at readers who are interested in making code contributions.

To learn more about how to create issues and pull requests in this repository, and what happens after they are created, you may refer to the resources below:
- [Pull Request submission and lifecycle](pull-request-lifecycle.md)

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
- [Go](https://golang.org/doc/install) >= 1.20

## Getting started
If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.20+ is required). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to our `$PATH`.

### Developing the SDK

While working on the provider, you may also be working on the [morpheus-go-sdk](https://github.com/gomorpheus/morpheus-go-sdk), which can be found at `$GOPATH/src/github.com/gomorpheus/morpheus-go-sdk`.

## Building the Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using `make dev`. This will place the provider onto your system in a [Terraform 0.13-compliant](https://www.terraform.io/upgrade-guides/0-13.html#in-house-providers) manner.

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

## Testing the Provider

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run. Please read [Writing Acceptance Tests](writing-tests.md) in the contribution guidelines for more information on usage.

```sh
$ make testacc
```

## Generating Docs

To generate or update documentation, run `make gendocs`.
```shell script
$ make gendocs
```

## Checklists

The following checklists are meant to be used for PRs to give developers and reviewers confidence that the proper changes have been made:

* [New resource](checklist-resource.md)
