package morpheus

import (
	"context"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSelectListOptionType() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus select list option type resource",
		CreateContext: resourceSelectListOptionTypeCreate,
		ReadContext:   resourceSelectListOptionTypeRead,
		UpdateContext: resourceSelectListOptionTypeUpdate,
		DeleteContext: resourceSelectListOptionTypeDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the select list option type",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the select list option type",
				Optional:    true,
			},
			"field_name": {
				Type:        schema.TypeString,
				Description: "The field name of the select list option type",
				Optional:    true,
				Default:     false,
			},
			"export_meta": {
				Type:        schema.TypeBool,
				Description: "Whether to export the select list option type as a tag",
				Optional:    true,
				Default:     false,
			},
			"dependent_field": {
				Type:        schema.TypeString,
				Description: "The field or code used to trigger the reloading of the field",
				Optional:    true,
			},
			"visibility_field": {
				Type:        schema.TypeString,
				Description: "The field or code used to trigger the visibility of the field",
				Optional:    true,
			},
			"display_value_on_details": {
				Type:        schema.TypeBool,
				Description: "Display the selected value of the text option type on the associated resource's details page",
				Optional:    true,
				Default:     false,
			},
			"field_label": {
				Type:        schema.TypeString,
				Description: "The label associated with the field in the UI",
				Optional:    true,
			},
			"default_value": {
				Type:        schema.TypeString,
				Description: "The default value of the option type",
				Optional:    true,
			},
			"help_block": {
				Type:        schema.TypeString,
				Description: "Text that provides additional details about the use of the option type",
				Optional:    true,
			},
			"required": {
				Type:        schema.TypeBool,
				Description: "Whether the option type is required",
				Optional:    true,
				Default:     false,
			},
			"option_list_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the associated option list",
				Optional:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceSelectListOptionTypeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	req := &morpheus.Request{
		Body: map[string]interface{}{
			"optionType": map[string]interface{}{
				"name":                  name,
				"description":           description,
				"fieldName":             d.Get("field_name").(string),
				"exportMeta":            d.Get("export_meta"),
				"dependsOnCode":         d.Get("dependent_field").(string),
				"visibleOnCode":         d.Get("visibility_field"),
				"displayValueOnDetails": d.Get("display_value_on_details"),
				"type":                  "select",
				"fieldLabel":            d.Get("field_label"),
				"defaultValue":          d.Get("default_value").(string),
				"helpBlock":             d.Get("help_block"),
				"required":              d.Get("required"),
				"optionList": map[string]interface{}{
					"id": d.Get("option_list_id").(int),
				},
			},
		},
	}
	resp, err := client.CreateOptionType(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateOptionTypeResult)
	environment := result.OptionType
	// Successfully created resource, now set id
	d.SetId(int64ToString(environment.ID))

	resourceSelectListOptionTypeRead(ctx, d, meta)
	return diags
}

func resourceSelectListOptionTypeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindOptionTypeByName(name)
	} else if id != "" {
		resp, err = client.GetOptionType(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("OptionType cannot be read without name or id")
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
	result := resp.Result.(*morpheus.GetOptionTypeResult)
	optionType := result.OptionType
	if optionType != nil {
		d.SetId(int64ToString(optionType.ID))
		d.Set("name", optionType.Name)
		d.Set("description", optionType.Description)
		d.Set("field_name", optionType.FieldName)
		d.Set("export_meta", optionType.ExportMeta)
		d.Set("dependent_field", optionType.DependsOnCode)
		d.Set("visibility_field", optionType.VisibleOnCode)
		d.Set("display_value_on_details", optionType.DisplayValueOnDetails)
		d.Set("type", optionType.Type)
		d.Set("field_label", optionType.FieldLabel)
		d.Set("default_value", optionType.DefaultValue)
		d.Set("help_block", optionType.HelpBlock)
		d.Set("required", optionType.Required)
		optionId := optionType.OptionListId
		d.Set("option_list_id", optionId["id"])
	} else {
		log.Println(optionType)
		return diag.Errorf("read operation: option type not found in response data") // should not happen
	}

	return diags
}

func resourceSelectListOptionTypeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	name := d.Get("name").(string)
	description := d.Get("description").(string)

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"optionType": map[string]interface{}{
				"name":                  name,
				"description":           description,
				"fieldName":             d.Get("field_name").(string),
				"exportMeta":            d.Get("export_meta"),
				"dependsOnCode":         d.Get("dependent_field").(string),
				"visibleOnCode":         d.Get("visibility_field"),
				"displayValueOnDetails": d.Get("display_value_on_details"),
				"type":                  "select",
				"fieldLabel":            d.Get("field_label"),
				"defaultValue":          d.Get("default_value").(string),
				"helpBlock":             d.Get("help_block"),
				"required":              d.Get("required"),
				"optionList": map[string]interface{}{
					"id": d.Get("option_list_id").(int),
				},
			},
		},
	}
	log.Printf("API REQUEST: %s", req)
	resp, err := client.UpdateOptionType(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateOptionTypeResult)
	account := result.OptionType
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(account.ID))
	return resourceSelectListOptionTypeRead(ctx, d, meta)
}

func resourceSelectListOptionTypeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteOptionType(toInt64(id), req)
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
