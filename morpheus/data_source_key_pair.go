package morpheus

import (
	"context"
	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMorpheusKeyPair() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus key pair data source.",
		ReadContext: dataSourceMorphesKeyPairRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Description:   "The ID of the key pair",
				Optional:      true,
				ConflictsWith: []string{"name"},
				Computed:      true,
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the integration",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
		},
	}
}

func dataSourceMorphesKeyPairRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	id := d.Get("id").(int)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == 0 && name != "" {
		resp, err = client.ListKeyPairs(&morpheus.Request{})
	} else {
		return diag.Errorf("Key pair cannot be read without name or id")
	}

	if err != nil {
		return diag.FromErr(err)
	}

	listResult := resp.Result.(*morpheus.ListKeyPairsResult)
	keyPairCount := len(*listResult.KeyPairs)
	if keyPairCount == 0 {
		return diag.Errorf("found %d key pairs", keyPairCount) // should not happen
	}

	var keyPairID int64
	for _, v := range *listResult.KeyPairs {
		if v.Name == name {
			keyPairID = v.ID
		}
	}
	log.Printf("KEY PAIR ID: %d", keyPairID)
	if keyPairID > 0 {
		resp, err = client.GetKeyPair(keyPairID, &morpheus.Request{})
	} else {
		return diag.Errorf("found %d key pairs for %s", keyPairID, name) // should not happen
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
	result := resp.Result.(*morpheus.GetKeyPairResult)
	keyPair := result.KeyPair
	if keyPair != nil {
		d.SetId(int64ToString(keyPair.ID))
		d.Set("name", keyPair.Name)
	} else {
		return diag.Errorf("Key pair not found in response data.") // should not happen
	}
	return diags
}
