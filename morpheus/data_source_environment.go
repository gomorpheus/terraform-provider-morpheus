package morpheus

import (
	"context"
	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMorpheusEnvironment() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus environment data source.",
		ReadContext: dataSourceMorpheusEnvironmentRead,
		Schema: map[string]*schema.Schema{
			"active": {
				Type:        schema.TypeBool,
				Description: "Whether the environment is active",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Morpheus environment",
				Optional:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the Morpheus environment",
				Computed:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "Optional code for use with policies",
				Computed:    true,
			},
			"visibility": {
				Type:        schema.TypeString,
				Description: "Whether the environment is visible in sub-tenants or not",
				Computed:    true,
			},
		},
	}
}

func dataSourceMorpheusEnvironmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindEnvironmentByName(name)
	} else if id != "" {
		resp, err = client.GetEnvironment(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Environment cannot be read without name or id")
	}
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %v", resp, err)
			return nil
		} else {
			log.Printf("API FAILURE: %s - %v", resp, err)
			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)

	// store resource data
	result := resp.Result.(*morpheus.GetEnvironmentResult)
	environment := result.Environment
	if environment != nil {
		d.SetId(int64ToString(environment.ID))
		d.Set("active", environment.Active)
		d.Set("name", environment.Name)
		d.Set("code", environment.Code)
		d.Set("description", environment.Description)
		d.Set("visibility", environment.Visibility)
	} else {
		return diag.Errorf("Environment not found in response data.") // should not happen
	}
	return diags
}
