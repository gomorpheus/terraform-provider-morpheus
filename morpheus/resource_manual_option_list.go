package morpheus

import (
	"context"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceManualOptionList() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus manual option list resource.",
		CreateContext: resourceManualOptionListCreate,
		ReadContext:   resourceManualOptionListRead,
		UpdateContext: resourceManualOptionListUpdate,
		DeleteContext: resourceManualOptionListDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the manual option list",
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
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "The organization labels associated with the option list (Only supported on Morpheus 5.5.3 or higher)",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"visibility": {
				Type:         schema.TypeString,
				Description:  "Whether the option list is visible in sub-tenants or not",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"private", "public", ""}, false),
				Default:      "private",
			},
			"dataset": {
				Type:             schema.TypeString,
				Description:      "The dataset for the manual option list",
				Optional:         true,
				DiffSuppressFunc: suppressEquivalentJsonDiffs,
			},
			"real_time": {
				Type:        schema.TypeBool,
				Description: "Whether the list is refreshed every time an associated option type is requested",
				Optional:    true,
				Default:     false,
			},
			"translation_script": {
				Type:             schema.TypeString,
				Description:      "A js script to translate the result data object into an Array containing objects with properties 'name’ and 'value’.",
				DiffSuppressFunc: supressOptionListScripts,
				Optional:         true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceManualOptionListCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		for _, s := range attr.(*schema.Set).List() {
			labelsPayload = append(labelsPayload, s.(string))
		}
	}
	req := &morpheus.Request{
		Body: map[string]interface{}{
			"optionTypeList": map[string]interface{}{
				"name":              name,
				"description":       description,
				"labels":            labelsPayload,
				"type":              "manual",
				"visibility":        d.Get("visibility"),
				"initialDataset":    d.Get("dataset").(string),
				"realTime":          d.Get("real_time").(bool),
				"translationScript": d.Get("translation_script").(string),
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

	resourceManualOptionListRead(ctx, d, meta)
	return diags
}

func resourceManualOptionListRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	result := resp.Result.(*morpheus.GetOptionListResult)
	optionList := result.OptionList
	if optionList != nil {
		d.SetId(int64ToString(optionList.ID))
		d.Set("name", optionList.Name)
		d.Set("description", optionList.Description)
		d.Set("labels", optionList.Labels)
		d.Set("type", optionList.Type)
		d.Set("visibility", optionList.Visibility)
		d.Set("dataset", optionList.InitialDataset)
		d.Set("real_time", optionList.RealTime)
		d.Set("translation_script", optionList.TranslationScript)
	} else {
		log.Println(optionList)
		return diag.Errorf("read operation: option list not found in response data") // should not happen
	}

	return diags
}

func resourceManualOptionListUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		for _, s := range attr.(*schema.Set).List() {
			labelsPayload = append(labelsPayload, s.(string))
		}
	}
	req := &morpheus.Request{
		Body: map[string]interface{}{
			"optionTypeList": map[string]interface{}{
				"name":              name,
				"description":       description,
				"labels":            labelsPayload,
				"type":              "manual",
				"visibility":        d.Get("visibility"),
				"initialDataset":    d.Get("dataset").(string),
				"realTime":          d.Get("real_time").(bool),
				"translationScript": d.Get("translation_script").(string),
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
	return resourceManualOptionListRead(ctx, d, meta)
}

func resourceManualOptionListDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
