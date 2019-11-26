module github.com/gomorpheus/terraform-provider-morpheus

go 1.13

require (
	github.com/gomorpheus/morpheusapi v0.0.0-00010101000000-000000000000
	github.com/hashicorp/terraform v0.12.13
	github.com/hashicorp/terraform-plugin-sdk v1.3.0
//	github.com/gomorpheus/morpheusapi v0.1
)

// voodoo
replace github.com/gomorpheus/morpheusapi => ../morpheusapi

replace github.com/go-resty/resty => gopkg.in/resty.v1 v1.11.0
