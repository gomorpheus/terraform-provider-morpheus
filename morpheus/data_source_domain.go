package morpheus

import (
	"context"
	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMorpheusDomain() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus domain data source.",
		ReadContext: dataSourceMorpheusDomainRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"name"},
				Computed:      true,
			},
			"active": {
				Type:        schema.TypeBool,
				Description: "Whether the domain is active",
				Computed:    true,
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the Morpheus domain",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the Morpheus domain",
				Computed:    true,
			},
			"visibility": {
				Type:        schema.TypeString,
				Description: "The visibility of the Morpheus domain",
				Computed:    true,
			},
		},
	}
}

func dataSourceMorpheusDomainRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	id := d.Get("id").(int)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == 0 && name != "" {
		resp, err = client.FindNetworkDomainByName(name)
	} else if id != 0 {
		resp, err = client.GetNetworkDomain(int64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Domain cannot be read without name or id")
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
	result := resp.Result.(*morpheus.GetNetworkDomainResult)
	domain := result.NetworkDomain
	if domain != nil {
		d.SetId(int64ToString(domain.ID))
		d.Set("active", domain.Active)
		d.Set("name", domain.Name)
		d.Set("description", domain.Description)
		d.Set("visibility", domain.Visibility)
	} else {
		return diag.Errorf("Domain not found in response data.") // should not happen
	}
	return diags
}
