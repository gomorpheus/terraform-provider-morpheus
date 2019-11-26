package morpheus

import (
	//"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// func Provider() *schema.Provider {
func Provider() terraform.ResourceProvider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{

			// todo: use environment defaults , just uncomment DefaultFunc below...

			"url": {
				Type:     schema.TypeString,
				Required: true,
				Description: "The URL of the Morpheus Data Appliance where requests will be directed.",
				// DefaultFunc: schema.MultiEnvDefaultFunc([]string{
				// 	"MORPHEUS_API_URL",
				// }, nil),
			},

			"access_token": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Access Token of Morpheus user. This can be used instead of authenticating with Username and Password.",
				// DefaultFunc: schema.MultiEnvDefaultFunc([]string{
				// 	"MORPHEUS_API_TOKEN",
				// }, nil),
				// ConflictsWith: []string{"username"},
			},

			"username": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Username of Morpheus user for authentication",
				// DefaultFunc: schema.MultiEnvDefaultFunc([]string{
				// 	"MORPHEUS_API_USERNAME",
				// }, nil),
				// ConflictsWith: []string{"access_token"},
			},

			"password": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Password of Morpheus user for authentication",
				// DefaultFunc: schema.MultiEnvDefaultFunc([]string{
				// 	"MORPHEUS_API_PASSWORD",
				// }, nil),
				// ConflictsWith: []string{"access_token"},
			},

		},

		ResourcesMap: map[string]*schema.Resource{
			"morpheus_group":          resourceMorpheusGroup(),
			"morpheus_cloud":          resourceCloud(),
			"morpheus_instance":       resourceInstance(),
			"morpheus_network_domain": resourceNetworkDomain(),
		},
	}
	
	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(d, terraformVersion)
	}
	return provider
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	config := Config{
		Url:                     d.Get("url").(string),
		AccessToken:             d.Get("access_token").(string),
		Username:                d.Get("username").(string),
		Password:                d.Get("password").(string),
		//Insecure:                d.Get("insecure").(bool), //.(bool),
		terraformVersion:        terraformVersion,
	}
	return config.Client()
}
