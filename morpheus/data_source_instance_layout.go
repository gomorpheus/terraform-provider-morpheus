package morpheus

import (
	"context"
	"fmt"
	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMorpheusInstanceLayout() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus instance layout data source.",
		ReadContext: dataSourceMorpheusInstanceLayoutRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"name", "version"},
				Computed:      true,
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the Morpheus instance layout",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
			"version": {
				Type:          schema.TypeString,
				Description:   "The version of the instance layout.",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
			"code": {
				Type:        schema.TypeString,
				Description: "Optional code for use with policies",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the instance layout",
				Computed:    true,
			},
		},
	}
}

func dataSourceMorpheusInstanceLayoutRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	id := d.Get("id").(int)
	version := d.Get("version").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == 0 && name != "" && version == "" {
		resp, err = client.FindInstanceLayoutByName(name)
	} else if id == 0 && name != "" && version != "" {
		resp, err = FindInstanceLayoutByNameAndVersion(client, name, version)
	} else if id != 0 {
		resp, err = client.GetInstanceLayout(int64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Instance layout cannot be read without name or id")
	}

	if err != nil {
		errorPrefix := "API FAILURE"
		if resp != nil && resp.StatusCode == 404 {
			errorPrefix = "API 404"
		}
		log.Printf("%s: %s - %v", errorPrefix, resp, err)
		return diag.FromErr(err)
	}

	log.Printf("API RESPONSE: %s", resp)

	// store resource data
	result := resp.Result.(*morpheus.GetInstanceLayoutResult)
	instanceLayout := result.InstanceLayout
	if instanceLayout != nil {
		d.SetId(int64ToString(instanceLayout.ID))
		d.Set("name", instanceLayout.Name)
		d.Set("code", instanceLayout.Code)
		d.Set("description", instanceLayout.Description)
		d.Set("version", instanceLayout.ContainerVersion)
	} else {
		return diag.Errorf("Instance layout not found in response data.") // should not happen
	}
	return diags
}

func FindInstanceLayoutByNameAndVersion(client *morpheus.Client, name string, version string) (*morpheus.Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListInstanceLayouts(&morpheus.Request{
		QueryParams: map[string]string{
			"name": name,
			"max":  "5000",
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*morpheus.ListInstanceLayoutsResult)
	for _, layout := range *listResult.InstanceLayouts {
		if layout.ContainerVersion == version {
			return client.GetInstanceLayout(layout.ID, &morpheus.Request{})
		}
	}
	return resp, fmt.Errorf("found 0 instance layouts named %v with a version of %v", name, version)
}
