module github.com/gomorpheus/terraform-provider-morpheus

go 1.14

require (
	github.com/gomorpheus/morpheus-go-sdk v0.0.0-00010101000000-000000000000
	github.com/hashicorp/terraform-plugin-docs v0.4.0
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.4.3
)

// voodoo
replace github.com/gomorpheus/morpheus-go-sdk => ../morpheus-go-sdk
