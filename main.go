// Copyright (c) 2019 Morpheus Data https://www.morpheusdata.com, All rights reserved.
// terraform-provider-morpheus source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	// "github.com/hashicorp/terraform/plugin"
	// "github.com/hashicorp/terraform/terraform"
	// "github.com/gomorpheus/terraform-provider-morpheus/morpheus"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/gomorpheus/terraform-provider-morpheus/morpheus"
)

func main() {

	// plugin.Serve(&plugin.ServeOpts{
	// 	ProviderFunc: func() terraform.ResourceProvider {
	// 		return morpheus.Provider()
	// 	},
	// })

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: morpheus.Provider})

}
