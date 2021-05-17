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
			"morpheus_checkbox_option_type": resourceCheckboxOptionType(),
			//			"morpheus_cloud":                   resourceCloud(),
			"morpheus_environment":             resourceEnvironment(),
			"morpheus_group":                   resourceMorpheusGroup(),
			"morpheus_hidden_option_type":      resourceHiddenOptionType(),
			"morpheus_manual_option_list":      resourceManualOptionList(),
			"morpheus_network_domain":          resourceNetworkDomain(),
			"morpheus_number_option_type":      resourceNumberOptionType(),
			"morpheus_operational_workflow":    resourceOperationalWorkflow(),
			"morpheus_password_option_type":    resourcePasswordOptionType(),
			"morpheus_provisioning_workflow":   resourceProvisioningWorkflow(),
			"morpheus_rest_option_list":        resourceRestOptionList(),
			"morpheus_select_list_option_type": resourceSelectListOptionType(),
			"morpheus_tenant":                  resourceTenant(),
			"morpheus_text_option_type":        resourceTextOptionType(),
			"morpheus_typeahead_option_type":   resourceTypeAheadOptionType(),
			"morpheus_vsphere_cloud":           resourceVsphereCloud(),
			"morpheus_vsphere_instance":        resourceVsphereInstance(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"morpheus_cloud":           dataSourceMorpheusCloud(),
			"morpheus_environment":     dataSourceMorpheusEnvironment(),
			"morpheus_group":           dataSourceMorpheusGroup(),
			"morpheus_instance_type":   dataSourceMorpheusInstanceType(),
			"morpheus_instance_layout": dataSourceMorpheusInstanceLayout(),
			"morpheus_network":         dataSourceMorpheusNetwork(),
			"morpheus_option_type":     dataSourceMorpheusOptionType(),
			"morpheus_plan":            dataSourceMorpheusPlan(),
			"morpheus_resource_pool":   dataSourceMorpheusResourcePool(),
			"morpheus_task":            dataSourceMorpheusTask(),
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
