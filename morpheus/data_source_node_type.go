package morpheus

import (
	"context"
	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMorpheusNodeType() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus node type data source.",
		ReadContext: dataSourceMorpheusNodeTypeRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Description:   "The ID of the node type",
				Optional:      true,
				ConflictsWith: []string{"name"},
				Computed:      true,
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the node type",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
			"provisioning_type": {
				Type:          schema.TypeString,
				Description:   "The provisioning type",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
		},
	}
}

func dataSourceMorpheusNodeTypeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	id := d.Get("id").(int)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == 0 && name != "" {
		resp, err = client.FindNodeType(name, d.Get("provisioning_type").(string))
	} else if id != 0 {
		resp, err = client.GetNodeType(int64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Node type cannot be read without name or id")
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
	result := resp.Result.(*morpheus.GetNodeTypeResult)
	nodeType := result.NodeType
	if nodeType != nil {
		d.SetId(int64ToString(nodeType.ID))
		d.Set("name", nodeType.Name)
		d.Set("provisioning_type", nodeType.ProvisionType.Code)
	} else {
		return diag.Errorf("Node type not found in response data.") // should not happen
	}
	return diags
}
