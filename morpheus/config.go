package morpheus

import (
	"fmt"
	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"os"
)

const (
	InsecureEnabledWarning = "Insecure mode enabled, this is NOT RECOMMENDED"
	InsecureChangeWarning  = "To enable TLS verification, unset environment variable MORPHEUS_INSECURE (insecure mode defaults to false) or set MORPHEUS_INSECURE to false"
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

	insecure bool

	client *morpheus.Client
}

func (c *Config) Client() (*morpheus.Client, diag.Diagnostics) {

	debug := logging.IsDebugOrHigher() && os.Getenv("MORPHEUS_API_HTTPTRACE") == "true"

	var diags diag.Diagnostics

	if c.client == nil {
		client := morpheus.NewClient(c.Url, morpheus.WithDebug(debug), morpheus.WithInsecure(c.insecure))
		// should validate url here too, and maybe ping it
		// logging with access token or username and password?

		if c.insecure {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  InsecureEnabledWarning,
				Detail:   InsecureChangeWarning,
			})
		}

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
