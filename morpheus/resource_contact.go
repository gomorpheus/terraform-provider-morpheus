package morpheus

import (
	"context"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceContact() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus contact resource.",
		CreateContext: resourceContactCreate,
		ReadContext:   resourceContactRead,
		UpdateContext: resourceContactUpdate,
		DeleteContext: resourceContactDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the contact",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the contact",
				Required:    true,
			},
			"email_address": {
				Type:        schema.TypeString,
				Description: "The email address associated with the contact",
				Optional:    true,
			},
			"mobile_number": {
				Type:        schema.TypeString,
				Description: "The mobile phone number associated with the contact",
				Optional:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceContactCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	req := &morpheus.Request{
		Body: map[string]interface{}{
			"contact": map[string]interface{}{
				"name":         name,
				"emailAddress": d.Get("email_address").(string),
				"smsAddress":   d.Get("mobile_number").(string),
			},
		},
	}

	resp, err := client.CreateContact(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateContactResult)
	contact := result.Contact
	// Successfully created resource, now set id
	d.SetId(int64ToString(contact.ID))

	resourceContactRead(ctx, d, meta)
	return diags
}

func resourceContactRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindContactByName(name)
	} else if id != "" {
		resp, err = client.GetContact(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Contact cannot be read without name or id")
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
	result := resp.Result.(*morpheus.GetContactResult)
	contact := result.Contact
	if contact != nil {
		d.SetId(int64ToString(contact.ID))
		d.Set("name", contact.Name)
		d.Set("email_address", contact.EmailAddress)
		d.Set("mobile_number", contact.SmsAddress)
	} else {
		return diag.Errorf("read operation: contact not found in response data") // should not happen
	}

	return diags
}

func resourceContactUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	name := d.Get("name").(string)
	req := &morpheus.Request{
		Body: map[string]interface{}{
			"contact": map[string]interface{}{
				"name":         name,
				"emailAddress": d.Get("email_address").(string),
				"smsAddress":   d.Get("mobile_number").(string),
			},
		},
	}
	log.Printf("API REQUEST: %s", req)
	resp, err := client.UpdateContact(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateContactResult)
	contact := result.Contact
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(contact.ID))
	return resourceTenantRead(ctx, d, meta)
}

func resourceContactDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteContact(toInt64(id), req)
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
