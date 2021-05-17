package morpheus

import (
	"context"
	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMorpheusGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus group data source.",
		ReadContext: dataSourceMorpheusGroupRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Morpheus group.",
				Optional:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "Optional code for use with policies",
				Computed:    true,
			},
			"location": {
				Type:        schema.TypeString,
				Description: "Optional location argument for your group",
				Computed:    true,
			},
		},
	}
}

func dataSourceMorpheusGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindGroupByName(name)
	} else if id != "" {
		resp, err = client.GetGroup(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Group cannot be read without name or id")
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
	result := resp.Result.(*morpheus.GetGroupResult)
	group := result.Group
	if group != nil {
		d.SetId(int64ToString(group.ID))
		d.Set("name", group.Name)
		d.Set("code", group.Code)
		d.Set("location", group.Location)
	} else {
		return diag.Errorf("Group not found in response data.") // should not happen
	}
	return diags
}
