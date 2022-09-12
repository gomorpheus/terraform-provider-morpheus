package morpheus

import (
	"context"
	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMorpheusScriptTemplate() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus script template data source.",
		ReadContext: dataSourceMorpheusScriptTemplateRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"name"},
				Computed:      true,
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the Morpheus script template.",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
		},
	}
}

func dataSourceMorpheusScriptTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	id := d.Get("id").(int)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == 0 && name != "" {
		resp, err = client.FindScriptTemplateByName(name)
	} else if id != 0 {
		resp, err = client.GetScriptTemplate(int64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Script template cannot be read without name or id")
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
	result := resp.Result.(*morpheus.GetScriptTemplateResult)
	scriptTemplate := result.ScriptTemplate
	if scriptTemplate != nil {
		d.SetId(int64ToString(scriptTemplate.ID))
		d.Set("name", scriptTemplate.Name)
	} else {
		return diag.Errorf("Script template not found in response data.") // should not happen
	}
	return diags
}
