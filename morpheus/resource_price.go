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

func resourcePrice() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a price resource",
		CreateContext: resourcePriceCreate,
		ReadContext:   resourcePriceRead,
		UpdateContext: resourcePriceUpdate,
		DeleteContext: resourcePriceDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the price",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the price",
				Required:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the price",
				Required:    true,
				ForceNew:    true,
			},
			"tenant_id": {
				Type:        schema.TypeInt,
				Description: "The id of the tenant to assign the price to",
				Optional:    true,
				ForceNew:    true,
			},
			"price_type": {
				Type:         schema.TypeString,
				Description:  "The price type (fixed, compute, memory, cores, storage, datastore, platform, software, load_balancer, load_balancer_virtual_server)",
				ValidateFunc: validation.StringInSlice([]string{"fixed", "compute", "memory", "cores", "storage", "datastore", "platform", "software", "load_balancer", "load_balancer_virtual_server"}, false),
				Required:     true,
			},
			"platform": {
				Type:         schema.TypeString,
				Description:  "The name of the platform (canonical, centos, debian, fedora, opensuse, redhat, suse, xen, linux, windows)",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"canonical", "centos", "debian", "fedora", "opensuse", "redhat", "suse", "xen", "linux", "windows"}, false),
			},
			"volume_type_id": {
				Type:        schema.TypeInt,
				Description: "The id of the volume type",
				Optional:    true,
				Computed:    true,
			},
			"software": {
				Type:        schema.TypeString,
				Description: "The name of the software",
				Optional:    true,
				Computed:    true,
			},
			"datastore_id": {
				Type:        schema.TypeInt,
				Description: "The id of the datastore to associate the price with",
				Optional:    true,
				Computed:    true,
			},
			"apply_price_accross_clouds": {
				Type:        schema.TypeBool,
				Description: "Whether to apply the datastore price across clouds",
				Optional:    true,
				Computed:    true,
			},
			"price_unit": {
				Type:         schema.TypeString,
				Description:  "The price unit (minute, hour, day, month, year, two year, three year, four year, five year)",
				ValidateFunc: validation.StringInSlice([]string{"minute", "hour", "day", "month", "year", "two year", "three year", "four year", "five year"}, false),
				Required:     true,
			},
			"incur_charges": {
				Type:         schema.TypeString,
				Description:  "When charges will be incurred (running, stopped, always)",
				ValidateFunc: validation.StringInSlice([]string{"running", "stopped", "always"}, false),
				Required:     true,
			},
			"currency": {
				Type:        schema.TypeString,
				Description: "The currency of the price",
				Required:    true,
			},
			"cost": {
				Type:        schema.TypeFloat,
				Description: "The cost of the price",
				Required:    true,
			},
			"markup_type": {
				Type:         schema.TypeString,
				Description:  "The type of markup applied to the cost (fixed, percent, custom)",
				ValidateFunc: validation.StringInSlice([]string{"fixed", "percent", "custom"}, false),
				Optional:     true,
				Computed:     true,
			},
			"markup_cost": {
				Type:          schema.TypeFloat,
				Description:   "The fixed cost at which the base cost is marked up",
				Optional:      true,
				ConflictsWith: []string{"markup_percent", "custom_price"},
			},
			"markup_percent": {
				Type:          schema.TypeFloat,
				Description:   "The percentage at which the base cost is marked up",
				Optional:      true,
				ConflictsWith: []string{"markup_cost", "custom_price"},
			},
			"custom_price": {
				Type:          schema.TypeFloat,
				Description:   "The custom price",
				Optional:      true,
				ConflictsWith: []string{"markup_cost", "markup_percent"},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourcePriceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	price := make(map[string]interface{})

	price["name"] = d.Get("name").(string)
	price["code"] = d.Get("code").(string)
	price["priceType"] = d.Get("price_type").(string)
	price["priceUnit"] = d.Get("price_unit").(string)
	price["incurCharges"] = d.Get("incur_charges").(string)
	price["currency"] = d.Get("currency").(string)
	price["cost"] = d.Get("cost").(float64)

	// Evaluate different markup types
	switch d.Get("markup_type") {
	case "fixed":
		price["markupType"] = "fixed"
		price["markup"] = d.Get("markup_cost").(float64)
	case "percent":
		price["markupType"] = "percent"
		price["markupPercent"] = d.Get("markup_percent").(float64)
	case "custom":
		price["markupType"] = "custom"
		price["customPrice"] = d.Get("custom_price")
	}

	if d.Get("tenant_id") != nil {
		price["account"] = map[string]interface{}{
			"id": d.Get("tenant_id"),
		}
	}

	// Evaluate different price types
	switch d.Get("price_type") {
	case "platform":
		if d.Get("platform") == "" {
			return diag.Errorf("A platform must be specified")
		} else {
			price["platform"] = d.Get("platform").(string)
		}
	case "software":
		price["software"] = d.Get("software").(string)
	case "storage":
		price["volumeType"] = map[string]interface{}{
			"id": d.Get("volume_type_id").(int64),
		}
	case "datastore":
		price["datastore"] = map[string]interface{}{
			"id": d.Get("datastore_id").(int64),
		}
		price["crossCloudApply"] = d.Get("apply_price_accross_clouds").(bool)
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"price": price,
		},
	}
	resp, err := client.CreatePrice(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreatePriceResult)
	// Successfully created resource, now set id
	d.SetId(int64ToString(result.ID))
	resourcePriceRead(ctx, d, meta)
	return diags
}

func resourcePriceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindPriceByName(name)
	} else if id != "" {
		resp, err = client.GetPrice(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Price cannot be read without name or id")
	}

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

	// store resource data
	var price MorpheusPrice
	json.Unmarshal(resp.Body, &price)

	if !price.Price.Active {
		d.SetId("")
		return diags
	}

	d.SetId(intToString(price.Price.ID))
	d.Set("name", price.Price.Name)
	d.Set("code", price.Price.Code)
	if _, ok := d.GetOk("tenant_id"); ok {
		d.Set("tenant_id", price.Price.Account.ID)
	}
	d.Set("price_type", price.Price.Pricetype)
	if _, ok := d.GetOk("platform"); ok {
		d.Set("platform", price.Price.Platform)
	}
	if _, ok := d.GetOk("volume_type_id"); ok {
		d.Set("volume_type_id", price.Price.Volumetype.ID)
	}
	if _, ok := d.GetOk("software"); ok {
		d.Set("software", price.Price.Software)
	}
	if _, ok := d.GetOk("datastore_id"); ok {
		d.Set("datastore_id", price.Price.Datastore.ID)
	}
	if _, ok := d.GetOk("apply_price_accross_clouds"); ok {
		d.Set("apply_price_accross_clouds", price.Price.Crosscloudapply)
	}
	d.Set("price_unit", price.Price.Priceunit)
	d.Set("incur_charges", price.Price.Incurcharges)
	d.Set("currency", price.Price.Currency)
	d.Set("cost", price.Price.Cost)
	if _, ok := d.GetOk("markup_type"); ok {
		d.Set("markup_type", price.Price.Markuptype)
	}
	if _, ok := d.GetOk("markup_cost"); ok {
		d.Set("markup_cost", price.Price.Markup)
	}
	if _, ok := d.GetOk("markup_percent"); ok {
		d.Set("markup_percent", price.Price.Markuppercent)
	}
	if _, ok := d.GetOk("custom_price"); ok {
		d.Set("custom_price", price.Price.Customprice)
	}

	return diags
}

func resourcePriceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	price := make(map[string]interface{})

	price["name"] = d.Get("name").(string)
	price["code"] = d.Get("code").(string)
	price["priceType"] = d.Get("price_type").(string)
	price["priceUnit"] = d.Get("price_unit").(string)
	price["incurCharges"] = d.Get("incur_charges").(string)
	price["currency"] = d.Get("currency").(string)
	price["cost"] = d.Get("cost").(float64)

	switch d.Get("markup_type") {
	case "fixed":
		price["markupType"] = "fixed"
		price["markup"] = d.Get("markup_cost").(float64)
	case "percent":
		price["markupType"] = "percent"
		price["markupPercent"] = d.Get("markup_percent").(float64)
	case "custom":
		price["markupType"] = "custom"
		price["customPrice"] = d.Get("custom_price")
	}

	// Evaluate different price types
	switch d.Get("price_type") {
	case "platform":
		price["platform"] = d.Get("platform").(string)
	case "software":
		price["software"] = d.Get("software").(string)
	case "storage":
		price["volumeType"] = map[string]interface{}{
			"id": d.Get("volume_type_id").(int64),
		}
	case "datastore":
		price["datastore"] = map[string]interface{}{
			"id": d.Get("datastore_id").(int64),
		}
		price["crossCloudApply"] = d.Get("apply_price_accross_clouds").(bool)
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"price": price,
		},
	}
	log.Printf("API REQUEST: %s", req)
	resp, err := client.UpdatePrice(toInt64(id), req)
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
	return resourcePriceRead(ctx, d, meta)
}

func resourcePriceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeletePrice(toInt64(id), req)
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

type MorpheusPrice struct {
	Price struct {
		ID                  int     `json:"id"`
		Name                string  `json:"name"`
		Code                string  `json:"code"`
		Active              bool    `json:"active"`
		Pricetype           string  `json:"priceType"`
		Priceunit           string  `json:"priceUnit"`
		Additionalpriceunit string  `json:"additionalPriceUnit"`
		Price               float64 `json:"price"`
		Customprice         float64 `json:"customPrice"`
		Markuptype          string  `json:"markupType"`
		Markup              float64 `json:"markup"`
		Markuppercent       float64 `json:"markupPercent"`
		Cost                float64 `json:"cost"`
		Currency            string  `json:"currency"`
		Incurcharges        string  `json:"incurCharges"`
		Platform            string  `json:"platform"`
		Software            string  `json:"software"`
		Volumetype          struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"volumeType"`
		Datastore struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"datastore"`
		Crosscloudapply bool `json:"crossCloudApply"`
		RestartUsage    bool `json:"restartUsage"`
		Account         struct {
			ID int `json:"id"`
		} `json:"account"`
	} `json:"price"`
}
