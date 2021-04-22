package morpheus

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The URL of the Morpheus Data Appliance where requests will be directed.",
				DefaultFunc: schema.EnvDefaultFunc("MORPHEUS_API_URL", nil),
			},

			"access_token": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				Description:   "Access Token of Morpheus user. This can be used instead of authenticating with Username and Password.",
				DefaultFunc:   schema.EnvDefaultFunc("MORPHEUS_API_TOKEN", nil),
				ConflictsWith: []string{"username", "password"},
			},

			"username": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Username of Morpheus user for authentication",
				DefaultFunc:   schema.EnvDefaultFunc("MORPHEUS_API_USERNAME", nil),
				ConflictsWith: []string{"access_token"},
			},

			"password": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				Description:   "Password of Morpheus user for authentication",
				DefaultFunc:   schema.EnvDefaultFunc("MORPHEUS_API_PASSWORD", nil),
				ConflictsWith: []string{"access_token"},
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"morpheus_group":          resourceMorpheusGroup(),
			"morpheus_cloud":          resourceCloud(),
			"morpheus_instance":       resourceInstance(),
			"morpheus_network_domain": resourceNetworkDomain(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := Config{
		Url:         d.Get("url").(string),
		AccessToken: d.Get("access_token").(string),
		Username:    d.Get("username").(string),
		Password:    d.Get("password").(string),
		//Insecure:                d.Get("insecure").(bool), //.(bool),
	}
	return config.Client()
}
