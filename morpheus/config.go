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

const sslCertErrorMsg = `
If you understand the potential security risks of accepting an untrusted server
certificate, you can bypass this error by setting "insecure = true" in your
provider configuration or by setting environment variable MORPHEUS_API_INSECURE to true. Use this option with caution.

provider "morpheus" {
	url = "https://..."
	.
	.
	.
	insecure = true <-- set to true to ignore SSL certificate errors
}
`

func certErrCallback(err error) error {
	var certErr x509.UnknownAuthorityError
	if errors.As(err, &certErr) {
		return errors.New(certErr.Error() + sslCertErrorMsg)
	}
	return nil
}

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

	insecure bool

	diags diag.Diagnostics

	client *morpheus.Client
}

func (c *Config) Client() (*morpheus.Client, diag.Diagnostics) {

	debug := logging.IsDebugOrHigher() && os.Getenv("MORPHEUS_API_HTTPTRACE") == "true"

	var diags diag.Diagnostics

	if c.client == nil {
		var client *morpheus.Client
		if c.insecure {
			client = morpheus.NewClient(c.Url, morpheus.WithDebug(debug), morpheus.Insecure())
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
