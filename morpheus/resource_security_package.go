package morpheus

import (
	"context"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSecurityPackage() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus security package resource",
		CreateContext: resourceSecurityPackageCreate,
		ReadContext:   resourceSecurityPackageRead,
		UpdateContext: resourceSecurityPackageUpdate,
		DeleteContext: resourceSecurityPackageDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the security package",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the security package",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the security package",
				Optional:    true,
				Computed:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "The organization labels associated with the security package (Only supported on Morpheus 5.5.3 or higher)",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the security package is enabled",
				Optional:    true,
				Default:     true,
			},
			"url": {
				Type:        schema.TypeString,
				Description: "The source url of the security package",
				Required:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceSecurityPackageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	securityPackage := make(map[string]interface{})

	securityPackage["name"] = d.Get("name").(string)
	securityPackage["type"] = "scap-package"
	securityPackage["description"] = d.Get("description").(string)
	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		for _, s := range attr.(*schema.Set).List() {
			labelsPayload = append(labelsPayload, s.(string))
		}
	}
	securityPackage["labels"] = labelsPayload
	securityPackage["enabled"] = d.Get("enabled").(bool)
	securityPackage["url"] = d.Get("url").(string)

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"securityPackage": securityPackage,
		},
	}
	resp, err := client.CreateSecurityPackage(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateSecurityPackageResult)
	securityPackageResult := result.SecurityPackage
	// Successfully created resource, now set id
	d.SetId(int64ToString(securityPackageResult.ID))

	resourceSecurityPackageRead(ctx, d, meta)
	return diags
}

func resourceSecurityPackageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindSecurityPackageByName(name)
	} else if id != "" {
		resp, err = client.GetSecurityPackage(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Security package cannot be read without name or id")
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
	result := resp.Result.(*morpheus.GetSecurityPackageResult)
	securityPackage := result.SecurityPackage

	d.SetId(intToString(int(securityPackage.ID)))
	d.Set("name", securityPackage.Name)
	d.Set("description", securityPackage.Description)
	d.Set("labels", securityPackage.Labels)
	d.Set("enabled", securityPackage.Enabled)
	d.Set("url", securityPackage.Url)

	return diags
}

func resourceSecurityPackageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	securityPackage := make(map[string]interface{})

	securityPackage["name"] = d.Get("name").(string)
	securityPackage["description"] = d.Get("description").(string)
	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		for _, s := range attr.(*schema.Set).List() {
			labelsPayload = append(labelsPayload, s.(string))
		}
	}
	securityPackage["labels"] = labelsPayload
	securityPackage["enabled"] = d.Get("enabled").(bool)
	securityPackage["url"] = d.Get("url").(string)

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"securityPackage": securityPackage,
		},
	}
	resp, err := client.UpdateSecurityPackage(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateSecurityPackageResult)
	securityPackageResult := result.SecurityPackage

	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(securityPackageResult.ID))
	return resourceSecurityPackageRead(ctx, d, meta)
}

func resourceSecurityPackageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteSecurityPackage(toInt64(id), req)
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
