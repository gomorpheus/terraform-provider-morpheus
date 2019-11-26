package morpheus

import (
	"github.com/gomorpheus/morpheus-go/morpheusapi"
)

// Config is the configuration structure used to instantiate the Morpheus
// provider.  Only Url and AccessToken are required.
type Config struct {
	Url              string
	AccessToken      string
	RefreshToken     string // optional and unused
	Username         string
	Password         string
	ClientId         string
	// Scope            string // "scope"
	// GrantType            string  // "bearer"

	Insecure         bool

	client           *morpheusapi.Client
	terraformVersion string
	userAgent        string
}

func (c *Config) Client() (*morpheusapi.Client, error) {
	if c.client == nil {
		client := morpheusapi.NewClient(c.Url)
		// should validate url here too, and maybe ping it
		// logging with access token or username and password?
		if c.Username != "" {
			client.SetUsernameAndPassword(c.Username, c.Password)
			// client.Login() // use lazy Login()
		} else {
			var expiresIn int64 = 86400 // lie (unused atm)
			client.SetAccessToken(c.AccessToken, c.RefreshToken, expiresIn, "write")
		}
		c.client = client;
	}
	return c.client, nil
}
