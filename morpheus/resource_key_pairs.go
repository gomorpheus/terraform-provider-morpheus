package morpheus

import (
	"context"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceKeyPair() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus key pair resource.",
		CreateContext: resourceKeyPairCreate,
		ReadContext:   resourceKeyPairRead,
		DeleteContext: resourceKeyPairDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the KeyPair",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the KeyPair",
				ForceNew:    true,
				Required:    true,
			},
			"publickey": {
				Type:        schema.TypeString,
				Description: "PublicKey of the KeyPair",
				ForceNew:    true,
				Required:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}
func resourceKeyPairCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	publicKey := d.Get("publickey").(string)

	resp, err := client.CreateKeyPair(name, publicKey)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateKeyPairResult)
	keyPair := result.KeyPair
	// Successfully created resource, now set id
	d.SetId(int64ToString(keyPair.ID))

	resourceKeyPairRead(ctx, d, meta)
	return diags
}

func resourceKeyPairRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	//read name and ID from Terraform
	name := d.Get("name").(string)
	id := toInt64(d.Id())

	var resp *morpheus.Response
	var err error
	if (id != 0 && name == "") || (id != 0 && name != "") { //if an ID is defined, use ID to retrieve KeyPair
		resp, err = client.GetKeyPair(id)
	} else if id == 0 && name != "" { // if no ID is defined search by name
		resp, err = client.GetKeyPairByName(name)
	} else if id == 0 && name == "" { // in case neither ID nore name is defined throw an error
		return diag.Errorf("Key pair cannot be read without name or id")
	}

	if err != nil {
		return diag.FromErr(err)
	}
	var keyPair *morpheus.KeyPair
	if id != 0 {
		result := resp.Result.(*morpheus.GetKeyPairResult) //read KeyPair from response, retireving an KeyPair via ID returns a pointer of type KeyPair, find by name return an pointer to an Array of KeyPair
		keyPair = result.KeyPair
	} else if name != "" {
		listResult := resp.Result.(*morpheus.ListKeyPairsResult)
		keyPairs := listResult.KeyPairs
		keyPair = &(*keyPairs)[0]
	}
	if keyPair != nil {
		d.SetId(int64ToString((*keyPair).ID))
		d.Set("name", (*keyPair).Name)
		d.Set("publickey", (*keyPair).PublicKey)
	} else {
		return diag.Errorf("Key pair not found in response data.") // should not happen
	}
	return diags
}

func resourceKeyPairDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	//name := d.Get("name").(string)
	id := toInt64(d.Get("id").(string))

	resp, err := client.DeleteKeyPair(id)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	resourceKeyPairRead(ctx, d, meta)
	return diags
}
