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

func dataSourceMorpheusCypherSecret() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus cypher secret data source.",
		ReadContext: dataSourceMorpheusCypherSecretRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"key": {
				Type:        schema.TypeString,
				Description: "The path of the cypher secret, excluding the secret prefix",
				Required:    true,
			},
			"value": {
				Type:        schema.TypeString,
				Description: "The cypher secret value",
				Computed:    true,
			},
			"ttl": {
				Type:        schema.TypeInt,
				Description: "The time to live of the cypher secret",
				Computed:    true,
			},
		},
	}
}

func dataSourceMorpheusCypherSecretRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var resp *morpheus.Response
	var err error
	secretPath := fmt.Sprintf("secret/%s", d.Get("key").(string))

	resp, err = client.Execute(&morpheus.Request{
		Method: "GET",
		Path:   fmt.Sprintf("%s/%s", "/api/cypher", secretPath),
		Result: &LocalGetCypherResult{},
	})

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
	cypher := resp.Result.(*LocalGetCypherResult)
	if cypher != nil {
		d.SetId(int64ToString(cypher.Cypher.ID))
		if cypher.Type == "object" {
			jsonPayload, _ := json.Marshal(cypher.Data)
			d.Set("value", string(jsonPayload))
		} else {
			d.Set("value", cypher.Data.(string))
		}
		d.Set("ttl", cypher.LeaseDuration)
	} else {
		return diag.Errorf("Cypher secret not found in response data.") // should not happen
	}
	return diags
}

type LocalGetCypherResult struct {
	Success       bool              `json:"success"`
	Data          interface{}       `json:"data"`
	Type          string            `json:"type"`
	LeaseDuration int64             `json:"lease_duration"`
	Cypher        *morpheus.Cypher  `json:"cypher"`
	Message       string            `json:"msg"`
	Errors        map[string]string `json:"errors"`
}
