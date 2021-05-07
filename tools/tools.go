// +build tools

package tools

//go:generate go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

import (
	// document generation
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
)
