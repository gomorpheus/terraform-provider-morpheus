package morpheus

import (
	"context"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLicense() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus license resource.",
		CreateContext: resourceLicenseCreate,
		ReadContext:   resourceLicenseRead,
		UpdateContext: resourceLicenseUpdate,
		DeleteContext: resourceLicenseDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the license",
				Computed:    true,
			},
			"license": {
				Type:        schema.TypeString,
				Description: "The morpheus license",
				Required:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceLicenseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"license": d.Get("license").(string),
		},
	}

	resp, err := client.InstallLicense(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.GetLicenseResult)
	_ = result.License
	// Successfully created resource, now set id
	d.SetId(int64ToString(1))

	resourceLicenseRead(ctx, d, meta)
	return diags
}

func resourceLicenseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error

	resp, err = client.GetLicense(&morpheus.Request{})
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
	result := resp.Result.(*morpheus.GetLicenseResult)
	_ = result.License
	d.SetId(int64ToString(1))

	return diags
}

func resourceLicenseUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"license": d.Get("license").(string),
		},
	}

	log.Printf("API REQUEST: %s", req)
	resp, err := client.InstallLicense(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.GetLicenseResult)
	_ = result.License
	// Successfully created resource, now set id
	d.SetId(int64ToString(1))

	return resourceLicenseRead(ctx, d, meta)
}

func resourceLicenseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	d.SetId("")
	return diags
}

func decodeLicense(ctx context.Context, d *schema.ResourceData, meta interface{}) (demo *morpheus.GetLicenseResult) {
	client := meta.(*morpheus.Client)
	resp, err := client.TestLicense(&morpheus.Request{})
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.GetLicenseResult)
	return result
}
