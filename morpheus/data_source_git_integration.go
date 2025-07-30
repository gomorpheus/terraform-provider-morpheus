package morpheus

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMorpheusGitIntegration() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus git integration data source.",
		ReadContext: dataSourceMorpheusGitIntegrationRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Description:   "The ID of the Morpheus git integration",
				Optional:      true,
				ConflictsWith: []string{"name"},
				Computed:      true,
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the Morpheus git integration.",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
			"repository_ids": {
				Computed:    true,
				Type:        schema.TypeMap,
				Description: "A map of git repository ids for use with integrations that reference a git repository",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func dataSourceMorpheusGitIntegrationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	id := d.Get("id").(int)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == 0 && name != "" {
		resp, err = client.FindIntegrationByName(name)
	} else if id != 0 {
		resp, err = client.GetIntegration(int64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Integration cannot be read without name or id")
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
	result := resp.Result.(*morpheus.GetIntegrationResult)
	integration := result.Integration
	if integration != nil {
		d.SetId(int64ToString(integration.ID))
		d.Set("name", integration.Name)
		resp, err = client.Execute(&morpheus.Request{
			Method:      "GET",
			Path:        fmt.Sprintf("/api/options/codeRepositories?integrationId=%d", integration.ID),
			QueryParams: map[string]string{},
		})
		if err != nil {
			log.Println("API ERROR: ", err)
		}
		log.Println("API RESPONSE:", resp)
		repo_ids := make(map[string]int)

		var itemResponsePayload CodeRepositories
		if err := json.Unmarshal(resp.Body, &itemResponsePayload); err != nil {
			return diag.FromErr(err)
		}
		for _, v := range itemResponsePayload.Data {
			repo_ids[v.Name] = v.Value
		}
		d.Set("repository_ids", repo_ids)
	} else {
		return diag.Errorf("Git integration not found in response data.") // should not happen
	}

	return diags
}
