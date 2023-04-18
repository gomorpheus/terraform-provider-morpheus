package morpheus

import (
	"context"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceRestOptionList() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus rest option list resource.",
		CreateContext: resourceRestOptionListCreate,
		ReadContext:   resourceRestOptionListRead,
		UpdateContext: resourceRestOptionListUpdate,
		DeleteContext: resourceRestOptionListDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the rest option list",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the option list",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the option list",
				Optional:    true,
				Computed:    true,
			},
			"visibility": {
				Type:         schema.TypeString,
				Description:  "Whether the option list is visible in sub-tenants or not",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"private", "public", ""}, false),
				Default:      "private",
			},
			"source_url": {
				Type:        schema.TypeString,
				Description: "The HTTP URL used for the API request",
				Optional:    true,
				Computed:    true,
			},
			"source_method": {
				Type:         schema.TypeString,
				Description:  "The HTTP method used for the API request",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"GET", "POST", ""}, false),
			},
			"source_headers": {
				Type:        schema.TypeList,
				Description: "An array of source headers to use when requesting data",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "The name of the source header",
							Optional:    true,
						},
						"value": {
							Type:        schema.TypeString,
							Description: "The value of the source header",
							Optional:    true,
						},
						// Masked is currently unsupported as there is no way to
						// evaluate the state difference as the API returns a masked value payload
						//"masked": {
						//	Type:        schema.TypeBool,
						//	Description: "Whether the source header value is masked or not",
						//	Optional:    true,
						//},
					},
				},
			},
			"real_time": {
				Type:        schema.TypeBool,
				Description: "Whether the list is refreshed every time an associated option type is requested",
				Optional:    true,
				Default:     false,
			},
			"ignore_ssl_errors": {
				Type:        schema.TypeBool,
				Description: "Whether to ignore SSL errors with the REST API endpoint",
				Optional:    true,
				Default:     false,
			},
			"initial_dataset": {
				Type:        schema.TypeString,
				Description: "The initial dataset used to populate the option list",
				Optional:    true,
				Computed:    true,
			},
			"translation_script": {
				Type:        schema.TypeString,
				Description: "A js script to translate the result data object into an Array containing objects with properties 'name’ and 'value’.",
				Optional:    true,
				Computed:    true,
			},
			"request_script": {
				Type:        schema.TypeString,
				Description: "A js script to prepare the API request",
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceRestOptionListCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	headers := d.Get("source_headers").([]interface{})
	var sourceHeaders []map[string]interface{}
	// iterate over the array of sourceHeaders
	for i := 0; i < len(headers); i++ {
		row := make(map[string]interface{})
		evarconfig := headers[i].(map[string]interface{})
		for k, v := range evarconfig {
			switch k {
			case "name":
				row["name"] = v.(string)
			case "value":
				row["value"] = v.(string)
				//case "masked":
				//	row["masked"] = v
			}
		}
		sourceHeaders = append(sourceHeaders, row)
		log.Printf("source headers payload: %s", sourceHeaders)
	}
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	req := &morpheus.Request{
		Body: map[string]interface{}{
			"optionTypeList": map[string]interface{}{
				"name":              name,
				"description":       description,
				"type":              "rest",
				"visibility":        d.Get("visibility"),
				"sourceUrl":         d.Get("source_url"),
				"realTime":          d.Get("real_time").(bool),
				"ignoreSSLErrors":   d.Get("ignore_ssl_errors"),
				"sourceMethod":      d.Get("source_method"),
				"initialDataset":    d.Get("initial_dataset").(string),
				"translationScript": d.Get("translation_script").(string),
				"requestScript":     d.Get("request_script").(string),
				"config": map[string]interface{}{
					"sourceHeaders": sourceHeaders,
				},
			},
		},
	}
	resp, err := client.CreateOptionList(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateOptionListResult)
	optionList := result.OptionList
	// Successfully created resource, now set id
	d.SetId(int64ToString(optionList.ID))

	resourceRestOptionListRead(ctx, d, meta)
	return diags
}

func resourceRestOptionListRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindOptionListByName(name)
	} else if id != "" {
		resp, err = client.GetOptionList(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Option list cannot be read without name or id")
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
	result := resp.Result.(*morpheus.GetOptionListResult)
	optionList := result.OptionList
	if optionList != nil {
		d.SetId(int64ToString(optionList.ID))
		d.Set("name", optionList.Name)
		d.Set("description", optionList.Description)
		d.Set("visibility", optionList.Visibility)
		d.Set("initial_dataset", optionList.InitialDataset)
		d.Set("real_time", optionList.RealTime)
		d.Set("translation_script", optionList.TranslationScript)
		d.Set("request_script", optionList.RequestScript)
		d.Set("source_url", optionList.SourceURL)
		d.Set("source_method", optionList.SourceMethod)
	} else {
		log.Println(optionList)
		return diag.Errorf("read operation: option list not found in response data") // should not happen
	}

	return diags
}

func resourceRestOptionListUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	name := d.Get("name").(string)
	description := d.Get("description").(string)

	headers := d.Get("source_headers").([]interface{})
	var sourceHeaders []map[string]interface{}
	// iterate over the array of sourceHeaders
	for i := 0; i < len(headers); i++ {
		row := make(map[string]interface{})
		evarconfig := headers[i].(map[string]interface{})
		for k, v := range evarconfig {
			switch k {
			case "name":
				row["name"] = v.(string)
			case "value":
				row["value"] = v.(string)
				//case "masked":
				//	row["masked"] = v
			}
		}
		sourceHeaders = append(sourceHeaders, row)
		log.Printf("source headers payload: %s", sourceHeaders)
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"optionTypeList": map[string]interface{}{
				"name":              name,
				"description":       description,
				"type":              "rest",
				"visibility":        d.Get("visibility"),
				"sourceUrl":         d.Get("source_url"),
				"realTime":          d.Get("real_time").(bool),
				"ignoreSSLErrors":   d.Get("ignore_ssl_errors"),
				"sourceMethod":      d.Get("source_method"),
				"initialDataset":    d.Get("initial_dataset").(string),
				"translationScript": d.Get("translation_script").(string),
				"requestScript":     d.Get("request_script").(string),
				"config": map[string]interface{}{
					"sourceHeaders": sourceHeaders,
				},
			},
		},
	}
	resp, err := client.UpdateOptionList(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateOptionListResult)
	optionList := result.OptionList
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(optionList.ID))
	return resourceRestOptionListRead(ctx, d, meta)
}

func resourceRestOptionListDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteOptionList(toInt64(id), req)
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
