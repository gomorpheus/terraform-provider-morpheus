package morpheus

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCypherSecret() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus cypher secret resource.",
		CreateContext: resourceCypherSecretCreate,
		ReadContext:   resourceCypherSecretRead,
		DeleteContext: resourceCypherSecretDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the cypher secret",
				Computed:    true,
			},
			"key": {
				Type:        schema.TypeString,
				Description: "The path of the cypher secret, excluding the secret prefix",
				Required:    true,
				ForceNew:    true,
			},
			"value": {
				Type:        schema.TypeString,
				Description: "The value of the cypher secret",
				Required:    true,
				Sensitive:   true,
				ForceNew:    true,
			},
			"ttl": {
				Type:        schema.TypeInt,
				Description: "The time to live of the cypher secret",
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCypherSecretCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"value": d.Get("value").(string),
		},
		QueryParams: map[string]string{
			"ttl":  strconv.Itoa(d.Get("ttl").(int)),
			"type": "string",
		},
	}

	secretPath := fmt.Sprintf("secret/%s", d.Get("key").(string))
	resp, err := client.CreateCypher(secretPath, req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	// Masking to avoid credential exposure
	// log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateCypherResult)
	// Successfully created resource, now set id
	d.SetId(int64ToString(result.Cypher.ID))

	resourceCypherSecretRead(ctx, d, meta)
	return diags
}

func resourceCypherSecretRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()

	var resp *morpheus.Response
	var err error
	if id != "" {
		secretPath := fmt.Sprintf("secret/%s", d.Get("key").(string))
		resp, err = client.GetCypher(secretPath, &morpheus.Request{})
	} else {
		return diag.Errorf("Cypher cannot be read without id")
	}

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
			log.Printf("Forcing recreation of resource")
			d.SetId("")
			return diags
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)
			return diag.FromErr(err)
		}
	}
	// Masking to avoid credential exposure
	// log.Printf("API RESPONSE: %s", resp)

	// store resource data
	result := resp.Result.(*morpheus.GetCypherResult)
	if result.Cypher != nil {
		d.SetId(int64ToString(result.Cypher.ID))
		keyData := strings.Split(result.Cypher.ItemKey, "/")
		keyData = keyData[1:]
		d.Set("key", strings.Join(keyData, "/"))
		d.Set("ttl", result.LeaseDuration)
	} else {
		return diag.Errorf("read operation: contact not found in response data") // should not happen
	}

	return diags
}

func resourceCypherSecretDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	req := &morpheus.Request{}
	secretPath := fmt.Sprintf("secret/%s", d.Get("key").(string))
	resp, err := client.DeleteCypher(secretPath, req)
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
