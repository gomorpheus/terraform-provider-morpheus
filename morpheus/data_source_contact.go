package morpheus

import (
	"context"
	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMorpheusContact() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus contact data source.",
		ReadContext: dataSourceMorpheusContactRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Morpheus contact.",
				Required:    true,
			},
		},
	}
}

func dataSourceMorpheusContactRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
			log.Printf("API 404: %s - %v", resp, err)
			return nil
		} else {
			log.Printf("API FAILURE: %s - %v", resp, err)
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
	} else {
		return diag.Errorf("Contact not found in response data.") // should not happen
	}
	return diags
}
