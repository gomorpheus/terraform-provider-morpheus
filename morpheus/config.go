package morpheus

import (
	"crypto/x509"
	"errors"
	"fmt"
	"os"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
)

// Config is the configuration structure used to instantiate the Morpheus
// provider.  Only Url and AccessToken are required.
type Config struct {
	Url             string
	AccessToken     string
	RefreshToken    string // optional and unused
	Username        string
	Password        string
	ClientId        string
	TenantSubdomain string
	// Scope            string // "scope"
	// GrantType            string  // "bearer"

	Insecure bool

	client *morpheus.Client
}

const sslCertErrorMsg = `
You have enabled SSL certificate verification, but the certificate presented by
the Morpheus server is not trusted. This could be due to a self-signed
certificate or an internal certificate authority.

We recommend fixing the certificate issue. If you need to bypass this check,
proceed with caution and understand the security implications of doing so. You can
disable certificate verification by either setting the provider argument
"secure = false" in your provider configuration or by setting the environment
variable MORPHEUS_API_SECURE to false.
provider "morpheus" {
	url = "https://..."
	.
	.
	.
	secure = false <-- set to false to disable certificate verification
}
`

func certErrCallback(err error) error {
	var certErr x509.UnknownAuthorityError
	if errors.As(err, &certErr) {
		return errors.New(certErr.Error() + sslCertErrorMsg)
	}
	return nil
}

func (c *Config) Client() (*morpheus.Client, diag.Diagnostics) {

	debug := logging.IsDebugOrHigher() && os.Getenv("MORPHEUS_API_HTTPTRACE") == "true"
	diags := diag.Diagnostics{}

	if c.client == nil {
		var client *morpheus.Client
		if c.Insecure {
			client = morpheus.NewClient(c.Url, morpheus.WithDebug(debug), morpheus.Insecure())
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "SSL Certificate Verification Disabled",
				Detail: `
SSL certificate verification is disabled. To enable verification, set "secure = true" in your
provider configuration or set environment variable MORPHEUS_API_SECURE to true.
provider "morpheus" {
	url = "https://..."
	.
	.
	.
	secure = true <-- set to true to enable certificate verification
}
`,
			})
		} else {
			client = morpheus.NewClient(c.Url, morpheus.WithDebug(debug), morpheus.WithErrCallbackFunc(certErrCallback))
		}

		// should validate url here too, and maybe ping it
		// logging with access token or username and password?
		if c.Username != "" {
			if c.TenantSubdomain != "" {
				username := fmt.Sprintf(`%s\\%s`, c.TenantSubdomain, c.Username)
				client.SetUsernameAndPassword(username, c.Password)
			} else {
				client.SetUsernameAndPassword(c.Username, c.Password)
			}
		} else {
			var expiresIn int64 = 86400 // lie (unused atm)
			client.SetAccessToken(c.AccessToken, c.RefreshToken, expiresIn, "write")
		}
		c.client = client
	}

	return c.client, diags
}
