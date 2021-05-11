package morpheus

import (
	"context"
	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNetworkDomain() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus network domain resource.",
		CreateContext: resourceNetworkDomainCreate,
		ReadContext:   resourceNetworkDomainRead,
		UpdateContext: resourceNetworkDomainUpdate,
		DeleteContext: resourceNetworkDomainDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the network domain",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The name of the network domain",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The user friendly description of the network domain",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"public_zone": {
				Description: "Whether the domain will be public or private",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"auto_join_domain": {
				Description: "Whether to automatically join machines to the domain",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"domain_controller": {
				Description: "The domain controller used to facilitate an automated domain join operation",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"domain_username": {
				Description: "The username of the account used to facilitate an automated domain join operation",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"domain_password": {
				Description: "The password of the account used to facilitate an automated domain join operation",
				Type:        schema.TypeString,
				Sensitive:   true,
				Optional:    true,
			},
			"active": {
				Description: "The state of the network domain",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"visibility": {
				Description:  "Determines whether the resource is visible in sub-tenants or not",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"private", "public", ""}, false),
				Default:      "private",
			},
			"tenant_id": {
				Description: "The tenant to assign the network domain",
				Type:        schema.TypeInt,
				Optional:    true,
			},
		},
	}
}

func resourceNetworkDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	// domainController := d.Get("domain_controller").(bool) // .(bool)
	//active := d.Get("active").(bool)
	req := &morpheus.Request{
		Body: map[string]interface{}{
			"networkDomain": map[string]interface{}{
				"name":        name,
				"description": description,
				"publicZone":  d.Get("public_zone").(bool),
				"visibility":  d.Get("visibility").(string),
				// "domainController": domainController,
				// "active":active,
			},
		},
	}
	resp, err := client.CreateNetworkDomain(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateNetworkDomainResult)
	networkDomain := result.NetworkDomain
	// Successfully created resource, now set id
	d.SetId(int64ToString(networkDomain.ID))
	resourceNetworkDomainRead(ctx, d, meta)
	return diags
}

func resourceNetworkDomainRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindNetworkDomainByName(name)
	} else if id != "" {
		resp, err = client.GetNetworkDomain(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("NetworkDomain cannot be read without name or id")
	}
	if err != nil {
		// 404 is ok?
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
	result := resp.Result.(*morpheus.GetNetworkDomainResult)
	networkDomain := result.NetworkDomain
	if networkDomain != nil {
		d.SetId(int64ToString(networkDomain.ID))
		d.Set("name", networkDomain.Name)
		d.Set("description", networkDomain.Description)
		d.Set("active", networkDomain.Active)
		d.Set("public_zone", networkDomain.PublicZone)
		d.Set("domain_controller", networkDomain.DomainController)
		d.Set("visibility", networkDomain.Visibility)
		// d.Set("fqdn", networkDomain.Fqdn)
	} else {
		return diag.Errorf("NetworkDomain not found in response data.") // should not happen
	}

	return diags
}

func resourceNetworkDomainUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	// publicZone := d.Get("public_zone").(bool) // .(bool)
	// domainController := d.Get("domain_controller").(bool) // .(bool)
	//active := d.Get("active").(bool)

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"networkDomain": map[string]interface{}{
				"name":        name,
				"description": description,
				// "publicZone": publicZone,
				// "domainController": domainController,
				//"active":active,
			},
		},
	}
	resp, err := client.UpdateNetworkDomain(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateNetworkDomainResult)
	networkDomain := result.NetworkDomain
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(networkDomain.ID))
	return resourceNetworkDomainRead(ctx, d, meta)
}

func resourceNetworkDomainDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteNetworkDomain(toInt64(id), req)
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
