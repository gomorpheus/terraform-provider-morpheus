// Copyright (c) 2019 Morpheus Data https://www.morpheusdata.com, All rights reserved.
// terraform-provider-morpheus source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/gomorpheus/terraform-provider-morpheus/morpheus"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return morpheus.Provider()
		},
	})
}
