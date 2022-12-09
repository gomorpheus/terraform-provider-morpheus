package morpheus

import (
	"context"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTenant() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus tenant resource.",
		CreateContext: resourceTenantCreate,
		ReadContext:   resourceTenantRead,
		UpdateContext: resourceTenantUpdate,
		DeleteContext: resourceTenantDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the tenant",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the tenant",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the tenant",
				Optional:    true,
				Computed:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the tenant is enabled or not",
				Optional:    true,
				Default:     true,
			},
			"subdomain": {
				Type:        schema.TypeString,
				Description: "Sets the custom login url or login prefix for logging into a sub-tenant user",
				Optional:    true,
				Computed:    true,
			},
			"base_role_id": {
				Type:        schema.TypeInt,
				Description: "The default base role for the account",
				Required:    true,
			},
			"currency": {
				Type:        schema.TypeString,
				Description: "Currency ISO Code to be used for the account",
				Optional:    true,
				Default:     "USD",
			},
			"account_number": {
				Type:        schema.TypeString,
				Description: "An optional field that can be used for billing and accounting",
				Optional:    true,
				Computed:    true,
			},
			"account_name": {
				Type:        schema.TypeString,
				Description: "An optional field that can be used for billing and accounting",
				Optional:    true,
				Computed:    true,
			},
			"customer_number": {
				Type:        schema.TypeString,
				Description: "An optional field that can be used for billing and accounting",
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTenantCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	req := &morpheus.Request{
		Body: map[string]interface{}{
			"account": map[string]interface{}{
				"name":        name,
				"description": description,
				"active":      d.Get("enabled").(bool),
				"subdomain":   d.Get("subdomain").(string),
				"role": map[string]interface{}{
					"id": d.Get("base_role_id").(int),
				},
				"currency":       d.Get("currency").(string),
				"accountNumber":  d.Get("account_number").(string),
				"accountName":    d.Get("account_name").(string),
				"customerNumber": d.Get("customer_number").(string),
			},
		},
	}

	resp, err := client.CreateTenant(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateTenantResult)
	tenant := result.Tenant
	// Successfully created resource, now set id
	d.SetId(int64ToString(tenant.ID))

	resourceTenantRead(ctx, d, meta)
	return diags
}

func resourceTenantRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindTenantByName(name)
	} else if id != "" {
		resp, err = client.GetTenant(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Tenant cannot be read without name or id")
	}

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
			return diag.FromErr(err)
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)
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
		d.Set("description", tenant.Description)
		d.Set("enabled", tenant.Active)
		d.Set("subdomain", tenant.Subdomain)
		d.Set("base_role_id", tenant.Role.ID)
		d.Set("currency", tenant.Currency)
		d.Set("account_number", tenant.AccountNumber)
		d.Set("account_name", tenant.AccountName)
		d.Set("customer_number", tenant.CustomerNumber)
	} else {
		log.Println(tenant)
		return diag.Errorf("read operation: option type not found in response data") // should not happen
	}

	return diags
}

func resourceTenantUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	name := d.Get("name").(string)
	description := d.Get("description").(string)

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"account": map[string]interface{}{
				"name":        name,
				"description": description,
				"active":      d.Get("enabled").(bool),
				"subdomain":   d.Get("subdomain").(string),
				"role": map[string]interface{}{
					"id": d.Get("base_role_id").(int),
				},
				"currency":       d.Get("currency").(string),
				"accountNumber":  d.Get("account_number").(string),
				"accountName":    d.Get("account_name").(string),
				"customerNumber": d.Get("customer_number").(string),
			},
		},
	}
	log.Printf("API REQUEST: %s", req)
	resp, err := client.UpdateTenant(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateTenantResult)
	account := result.Tenant
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(account.ID))
	return resourceTenantRead(ctx, d, meta)
}

func resourceTenantDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteTenant(toInt64(id), req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
			return diag.FromErr(err)
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)
			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)
	d.SetId("")
	return diags
}
