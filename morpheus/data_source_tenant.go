package morpheus

import (
	"context"
	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMorpheusTenant() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus tenant data source.",
		ReadContext: dataSourceMorpheusTenantRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"name"},
				Computed:      true,
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the Morpheus tenant.",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
			"account_number": {
				Type:        schema.TypeString,
				Description: "An optional field that can be used for billing and accounting",
				Computed:    true,
			},
			"account_name": {
				Type:        schema.TypeString,
				Description: "An optional field that can be used for billing and accounting",
				Computed:    true,
			},
			"customer_number": {
				Type:        schema.TypeString,
				Description: "An optional field that can be used for billing and accounting",
				Computed:    true,
			},
		},
	}
}

func dataSourceMorpheusTenantRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	id := d.Get("id").(int)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == 0 && name != "" {
		resp, err = client.FindTenantByName(name)
	} else if id != 0 {
		resp, err = client.GetTenant(int64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Tenant cannot be read without name or id")
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
	result := resp.Result.(*morpheus.GetTenantResult)
	tenant := result.Tenant
	if tenant != nil {
		d.SetId(int64ToString(tenant.ID))
		d.Set("name", tenant.Name)
		d.Set("account_number", tenant.AccountNumber)
		d.Set("account_name", tenant.AccountName)
		d.Set("customer_number", tenant.CustomerNumber)
	} else {
		return diag.Errorf("Tenant not found in response data.") // should not happen
	}
	return diags
}
