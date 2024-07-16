package morpheus

import (
	"context"
	"encoding/json"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourcePriceSet() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a price set resource",
		CreateContext: resourcePriceSetCreate,
		ReadContext:   resourcePriceSetRead,
		UpdateContext: resourcePriceSetUpdate,
		DeleteContext: resourcePriceSetDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the price set",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the price set",
				Required:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the price set",
				Required:    true,
				ForceNew:    true,
			},
			"region_code": {
				Type:        schema.TypeString,
				Description: "The region code of the price set",
				Required:    true,
				ForceNew:    true,
			},
			"cloud_id": {
				Type:        schema.TypeInt,
				Description: "The id of the cloud",
				Optional:    true,
				ForceNew:    true,
			},
			"resource_pool_id": {
				Type:        schema.TypeInt,
				Description: "The resource pool to assign the price set to",
				Optional:    true,
				ForceNew:    true,
			},
			"type": {
				Type:         schema.TypeString,
				Description:  "The price type (fixed, compute, memory, cores, storage, datastore, platform, software_or_service, load_balancer, load_balancer_virtual_server)",
				ValidateFunc: validation.StringInSlice([]string{"fixed", "compute", "memory", "cores", "storage", "datastore", "platform", "software_or_service", "load_balancer", "load_balancer_virtual_server"}, false),
				Required:     true,
				ForceNew:     true,
			},
			"price_unit": {
				Type:         schema.TypeString,
				Description:  "The price unit (minute, hour, day, month, year, two year, three year, four year, five year)",
				ValidateFunc: validation.StringInSlice([]string{"minute", "hour", "day", "month", "year", "two year", "three year", "four year", "five year"}, false),
				Required:     true,
				ForceNew:     true,
			},
			"price_ids": {
				Type:        schema.TypeList,
				Description: "The list of price ids associated with the price set",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourcePriceSetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	priceSet := make(map[string]interface{})

	priceSet["name"] = d.Get("name").(string)
	priceSet["code"] = d.Get("code").(string)
	priceSet["regionCode"] = d.Get("region_code").(string)
	priceSet["zone"] = map[string]interface{}{
		"id": d.Get("cloud_id").(int),
	}
	priceSet["zonePool"] = map[string]interface{}{
		"id": d.Get("resource_pool_id").(int),
	}
	priceSet["priceUnit"] = d.Get("price_unit").(string)
	priceSet["type"] = d.Get("type").(string)
	//price["restartUsage"] = d.Get("restart_usage").(bool)

	var priceIds []map[string]interface{}
	if d.Get("price_ids") != nil {
		priceIdList := d.Get("price_ids").([]interface{})
		// iterate over the array of price ids
		for i := 0; i < len(priceIdList); i++ {
			row := make(map[string]interface{})
			row["id"] = priceIdList[i]
			priceIds = append(priceIds, row)
		}
	}
	priceSet["prices"] = priceIds

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"priceSet": priceSet,
		},
	}
	resp, err := client.CreatePriceSet(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreatePriceSetResult)
	// Successfully created resource, now set id
	d.SetId(int64ToString(result.ID))
	resourcePriceSetRead(ctx, d, meta)
	return diags
}

func resourcePriceSetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindPriceSetByName(name)
	} else if id != "" {
		resp, err = client.GetPriceSet(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Price cannot be read without name or id")
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
	log.Printf("API RESPONSE: %s", resp)

	// store resource data
	var priceSet MorpheusPriceSet
	json.Unmarshal(resp.Body, &priceSet)

	// Remove the resource from state if the
	// price set has
	if !priceSet.Priceset.Active {
		d.SetId("")
		return diags
	}

	d.SetId(intToString(priceSet.Priceset.ID))
	d.Set("name", priceSet.Priceset.Name)
	d.Set("code", priceSet.Priceset.Code)
	d.Set("region_code", priceSet.Priceset.Regioncode)
	d.Set("cloud_id", priceSet.Priceset.Zone.ID)

	if _, ok := d.GetOk("resource_pool_id"); ok {
		d.Set("resource_pool_id", priceSet.Priceset.Zonepool.ID)
	}

	d.Set("price_unit", priceSet.Priceset.Priceunit)
	d.Set("type", priceSet.Priceset.Type)
	var priceIds []int
	if len(priceSet.Priceset.Prices) > 0 {
		for _, v := range priceSet.Priceset.Prices {
			priceIds = append(priceIds, int(v.ID))
		}
	}

	// Adjust the price set order returned from the API (numerical low to high)
	// to match the order defined in the Terraform code.
	statePricePayload := matchPricesWithSchema(priceIds, d.Get("price_ids").([]interface{}))
	d.Set("price_ids", statePricePayload)

	return diags
}

func resourcePriceSetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	priceSet := make(map[string]interface{})

	priceSet["name"] = d.Get("name").(string)
	priceSet["code"] = d.Get("code").(string)
	priceSet["regionCode"] = d.Get("region_code").(string)
	priceSet["zone"] = map[string]interface{}{
		"id": d.Get("cloud_id").(int),
	}
	priceSet["zonePool"] = map[string]interface{}{
		"id": d.Get("resource_pool_id").(int),
	}
	priceSet["priceUnit"] = d.Get("price_unit").(string)
	priceSet["type"] = d.Get("type").(string)
	//price["restartUsage"] = d.Get("restart_usage").(bool)

	var priceIds []map[string]interface{}
	if d.Get("price_ids") != nil {
		priceIdList := d.Get("price_ids").([]interface{})
		// iterate over the array of price ids
		for i := 0; i < len(priceIdList); i++ {
			row := make(map[string]interface{})
			row["id"] = priceIdList[i]
			priceIds = append(priceIds, row)
		}
	}
	priceSet["prices"] = priceIds

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"priceSet": priceSet,
		},
	}

	resp, err := client.UpdatePriceSet(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		log.Fatal(err)
	}

	d.SetId(id)
	return resourcePriceSetRead(ctx, d, meta)
}

func resourcePriceSetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeletePriceSet(toInt64(id), req)
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

// This cannot currently be handled efficiently by a DiffSuppressFunc.
// See: https://github.com/hashicorp/terraform-plugin-sdk/issues/477
func matchPricesWithSchema(prices []int, declaredPrices []interface{}) []int {
	result := make([]int, len(declaredPrices))

	rMap := make(map[int]int, len(prices))
	for _, price := range prices {
		rMap[price] = price
	}

	for i, declaredPrice := range declaredPrices {
		declaredPrice := declaredPrice.(int)

		if v, ok := rMap[declaredPrice]; ok {
			// matched price declared by ID
			result[i] = v
			delete(rMap, v)
		}
	}
	// append unmatched price set to the result
	for _, rcpt := range rMap {
		result = append(result, rcpt)
	}
	return result
}

type MorpheusPriceSet struct {
	Priceset struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		Code          string `json:"code"`
		Active        bool   `json:"active"`
		Priceunit     string `json:"priceUnit"`
		Type          string `json:"type"`
		Regioncode    string `json:"regionCode"`
		Systemcreated bool   `json:"systemCreated"`
		Zone          struct {
			ID int `json:"id"`
		} `json:"zone"`
		Zonepool struct {
			ID int `json:"id"`
		} `json:"zonePool"`
		Account interface{} `json:"account"`
		Prices  []struct {
			ID                  int         `json:"id"`
			Name                string      `json:"name"`
			Code                string      `json:"code"`
			Pricetype           string      `json:"priceType"`
			Priceunit           string      `json:"priceUnit"`
			Additionalpriceunit string      `json:"additionalPriceUnit"`
			Price               float64     `json:"price"`
			Customprice         float64     `json:"customPrice"`
			Markuptype          interface{} `json:"markupType"`
			Markup              float64     `json:"markup"`
			Markuppercent       float64     `json:"markupPercent"`
			Cost                float64     `json:"cost"`
			Currency            string      `json:"currency"`
			Incurcharges        string      `json:"incurCharges"`
			Platform            interface{} `json:"platform"`
			Software            interface{} `json:"software"`
			Volumetype          struct {
				ID   int    `json:"id"`
				Code string `json:"code"`
				Name string `json:"name"`
			} `json:"volumeType"`
			Datastore       interface{} `json:"datastore"`
			Crosscloudapply interface{} `json:"crossCloudApply"`
			Account         interface{} `json:"account"`
		} `json:"prices"`
	} `json:"priceSet"`
}
