package morpheus

import (
	"context"
	"encoding/json"
	"strconv"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceForm() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus form resource",
		CreateContext: resourceFormCreate,
		ReadContext:   resourceFormRead,
		UpdateContext: resourceFormUpdate,
		DeleteContext: resourceFormDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The id of the form",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the form",
				Required:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the form",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the form",
				Optional:    true,
				Computed:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "The organization labels associated with the file template (Only supported on Morpheus 5.5.3 or higher)",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"option_type": {
				Type:        schema.TypeList,
				Description: "Form option type",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Description: "The id of the option type to add to the form",
							Optional:    true,
						},
						"code": {
							Type:        schema.TypeString,
							Description: "The id of the option type to add to the form",
							Optional:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "The id of the option type to add to the form",
							Optional:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "The id of the option type to add to the form",
							Optional:    true,
						},
						"field_name": {
							Type:        schema.TypeString,
							Description: "The id of the option type to add to the form",
							Optional:    true,
						},
						"type": {
							Type:         schema.TypeString,
							Description:  "The id of the option type to add to the form (checkbox, hidden, number, password, radio, text)",
							ValidateFunc: validation.StringInSlice([]string{"checkbox", "hidden", "number", "password", "radio", "text"}, false),
							Optional:     true,
						},
						"option_list_id": {
							Type:        schema.TypeInt,
							Description: "The id of the option type to add to the form",
							Optional:    true,
						},
						"field_label": {
							Type:        schema.TypeString,
							Description: "The id of the option type to add to the form",
							Optional:    true,
						},
						"default_value": {
							Type:        schema.TypeString,
							Description: "The id of the option type to add to the form",
							Optional:    true,
						},
						"placeholder": {
							Type:        schema.TypeString,
							Description: "The id of the option type to add to the form",
							Optional:    true,
						},
						"help_block": {
							Type:        schema.TypeString,
							Description: "The id of the option type to add to the form",
							Optional:    true,
						},
						"required": {
							Type:        schema.TypeBool,
							Description: "The id of the option type to add to the form",
							Optional:    true,
						},
						"export_meta": {
							Type:        schema.TypeBool,
							Description: "Whether to export the text option type as a tag",
							Optional:    true,
							Default:     false,
						},
						"display_value_on_details": {
							Type:        schema.TypeBool,
							Description: "Display the selected value of the text option type on the associated resource's details page",
							Optional:    true,
							Default:     false,
						},
						"locked": {
							Type:        schema.TypeBool,
							Description: "Display the selected value of the text option type on the associated resource's details page",
							Optional:    true,
							Default:     false,
						},
						"hidden": {
							Type:        schema.TypeBool,
							Description: "Display the selected value of the text option type on the associated resource's details page",
							Optional:    true,
							Default:     false,
						},
						"exclude_from_search": {
							Type:        schema.TypeBool,
							Description: "Display the selected value of the text option type on the associated resource's details page",
							Optional:    true,
							Default:     false,
						},
						"allow_password_peek": {
							Type:        schema.TypeBool,
							Description: "The field or code used to trigger the reloading of the field",
							Optional:    true,
							Computed:    true,
						},
						"min_value": {
							Type:        schema.TypeInt,
							Description: "The field or code used to trigger the reloading of the field",
							Optional:    true,
							Computed:    true,
						},
						"max_value": {
							Type:        schema.TypeInt,
							Description: "The field or code used to trigger the reloading of the field",
							Optional:    true,
							Computed:    true,
						},
						"step": {
							Type:        schema.TypeInt,
							Description: "The field or code used to trigger the reloading of the field",
							Optional:    true,
							Computed:    true,
						},
						"dependent_field": {
							Type:        schema.TypeString,
							Description: "The field or code used to trigger the reloading of the field",
							Optional:    true,
							Computed:    true,
						},
						"visibility_field": {
							Type:        schema.TypeString,
							Description: "The field or code used to trigger the visibility of the field",
							Optional:    true,
							Computed:    true,
						},
						"verify_pattern": {
							Type:        schema.TypeString,
							Description: "The regex pattern used to validate the entered",
							Optional:    true,
							Computed:    true,
						},
						"require_field": {
							Type:        schema.TypeString,
							Description: "The field or code used to determine whether the field is required or not",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"field_group": {
				Type:        schema.TypeList,
				Description: "Field group to add to the form",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "The name of the field group",
							Required:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "Whether to mark the cloud datastore as a default store for this tenant",
							Optional:    true,
							Computed:    true,
						},
						"collapsible": {
							Type:        schema.TypeBool,
							Description: "Whether to mark the cloud datastore as an image target for this tenant",
							Optional:    true,
							Computed:    true,
						},
						"collapsed_by_deafult": {
							Type:        schema.TypeBool,
							Description: "Whether to mark the cloud datastore as an image target for this tenant",
							Optional:    true,
							Computed:    true,
						},
						"visibility_field": {
							Type:        schema.TypeString,
							Description: "Whether to mark the cloud datastore as an image target for this tenant",
							Optional:    true,
							Computed:    true,
						},
						"option_type": {
							Type:        schema.TypeList,
							Description: "Form option type",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Description: "The id of the option type to add to the form",
										Optional:    true,
										Computed:    true,
									},
									"code": {
										Type:        schema.TypeString,
										Description: "The id of the option type to add to the form",
										Optional:    true,
									},
									"name": {
										Type:        schema.TypeString,
										Description: "The id of the option type to add to the form",
										Optional:    true,
									},
									"description": {
										Type:        schema.TypeString,
										Description: "The id of the option type to add to the form",
										Optional:    true,
									},
									"field_name": {
										Type:        schema.TypeString,
										Description: "The id of the option type to add to the form",
										Optional:    true,
									},
									"type": {
										Type:        schema.TypeString,
										Description: "The id of the option type to add to the form",
										Optional:    true,
									},
									"label": {
										Type:        schema.TypeString,
										Description: "The id of the option type to add to the form",
										Optional:    true,
									},
								},
							},
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceFormCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)

	// optiontypes
	var optionTypes []map[string]interface{}
	if d.Get("option_type") != nil {
		optionTypeList := d.Get("option_type").([]interface{})
		// iterate over the array of optionTypes
		for i := 0; i < len(optionTypeList); i++ {
			row := make(map[string]interface{})
			optionTypeConfig := optionTypeList[i].(map[string]interface{})
			if optionTypeConfig["id"].(int) > 0 {
				row["id"] = optionTypeConfig["id"]
			} else {
				row["name"] = optionTypeConfig["name"]
				row["code"] = optionTypeConfig["code"]
				row["type"] = optionTypeConfig["type"]
				row["description"] = optionTypeConfig["description"]
				row["fieldName"] = optionTypeConfig["field_name"]
				row["fieldLabel"] = optionTypeConfig["field_label"]
				row["placeHolder"] = optionTypeConfig["placeholder"]
				row["helpBlock"] = optionTypeConfig["help_block"]
				switch optionTypeConfig["type"] {
				case "checkbox":
					row["defaultValue"] = optionTypeConfig["default_checked"]
				case "number":
					n, err := strconv.Atoi(optionTypeConfig["default_value"].(string))
					if err != nil {
						return diag.Errorf("The default_value attribute must be a number string when the type attribute is set to number")
					}
					row["defaultValue"] = n
				case "password":
					configStep := make(map[string]interface{})
					configStep["canPeek"] = optionTypeConfig["allow_password_peek"]
					row["config"] = configStep
				default:
					row["defaultValue"] = optionTypeConfig["default_value"]
				}
				row["required"] = optionTypeConfig["required"]
				row["exportMeta"] = optionTypeConfig["export_meta"]
				row["editable"] = optionTypeConfig["editable"]
				if optionTypeConfig["option_list_id"].(int) > 0 {
					optionList := make(map[string]interface{})
					optionList["id"] = optionTypeConfig["option_list_id"]
					row["optionList"] = optionList
				}
				row["displayValueOnDetails"] = optionTypeConfig["display_value_on_details"]
				row["isLocked"] = optionTypeConfig["locked"]
				row["isHidden"] = optionTypeConfig["hidden"]
				row["excludeFromSearch"] = optionTypeConfig["exclude_from_search"]
				row["minVal"] = optionTypeConfig["min_value"]
				row["maxVal"] = optionTypeConfig["max_value"]
				if optionTypeConfig["step"].(int) > 0 {
					configStep := make(map[string]interface{})
					configStep["step"] = optionTypeConfig["step"]
					row["config"] = configStep
				}
				row["dependsOnCode"] = optionTypeConfig["dependent_field"]
				row["visibleOnCode"] = optionTypeConfig["visibility_field"]
				row["verifyPattern"] = optionTypeConfig["verify_pattern"]
				row["requireOnCode"] = optionTypeConfig["require_field"]
			}
			optionTypes = append(optionTypes, row)
		}
	}

	// fieldGroups
	var fieldGroups []map[string]interface{}
	if d.Get("field_group") != nil {
		fieldGroupList := d.Get("field_group").([]interface{})
		// iterate over the array of fieldGroups
		for i := 0; i < len(fieldGroupList); i++ {
			row := make(map[string]interface{})
			fieldGroupConfig := fieldGroupList[i].(map[string]interface{})
			row["name"] = fieldGroupConfig["name"]
			row["description"] = fieldGroupConfig["description"]
			row["collapsible"] = fieldGroupConfig["collapsible"]
			row["defaultCollapsed"] = fieldGroupConfig["collapsed_by_deafult"]
			row["visibleOnCode"] = fieldGroupConfig["visibility_field"]
			// optiontypes
			var optionTypes []map[string]interface{}
			if fieldGroupConfig["option_type"] != nil {
				optionTypeList := fieldGroupConfig["option_type"].([]interface{})
				// iterate over the array of optionTypes
				for i := 0; i < len(optionTypeList); i++ {
					row := make(map[string]interface{})
					optionTypeConfig := optionTypeList[i].(map[string]interface{})
					row["id"] = optionTypeConfig["id"]
					optionTypes = append(optionTypes, row)
				}
			}
			row["options"] = optionTypes
			fieldGroups = append(fieldGroups, row)
		}
	}

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		for _, s := range attr.(*schema.Set).List() {
			labelsPayload = append(labelsPayload, s.(string))
		}
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"optionTypeForm": map[string]interface{}{
				"name":        name,
				"code":        d.Get("code").(string),
				"description": d.Get("description").(string),
				"labels":      labelsPayload,
				"fieldGroups": fieldGroups,
				"options":     optionTypes,
			},
		},
	}
	jsonRequest, _ := json.Marshal(req.Body)
	log.Printf("API JSON REQUEST: %s", string(jsonRequest))

	resp, err := client.CreateForm(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateFormResult)
	formResult := result.Form
	// Successfully created resource, now set id
	d.SetId(int64ToString(formResult.ID))

	resourceFormRead(ctx, d, meta)
	return diags
}

func resourceFormRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error

	if id == "" && name != "" {
		resp, err = client.FindFormByName(name)
	} else if id != "" {
		resp, err = client.GetForm(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Form cannot be read without name or id")
	}

	if err != nil {
		// 404 is ok?
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
			d.SetId("")
			return diags
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)
			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)

	// store resource data
	result := resp.Result.(*morpheus.GetFormResult)
	form := result.Form

	d.SetId(int64ToString(form.ID))
	d.Set("name", form.Name)
	d.Set("code", form.Code)
	d.Set("description", form.Description)
	d.Set("labels", form.Labels)

	// Option Types

	var optionTypes []map[string]interface{}
	if len(form.Options) != 0 {
		for _, optionType := range form.Options {
			row := make(map[string]interface{})
			switch optionType.Type {
			case "number":
				row["min_value"] = optionType.MinVal
				row["max_value"] = optionType.MaxVal
			}
			row["id"] = optionType.ID
			row["name"] = optionType.Name
			row["code"] = optionType.Code
			row["type"] = optionType.Type
			row["field_label"] = optionType.FieldLabel
			row["field_name"] = optionType.FieldName
			row["default_value"] = optionType.DefaultValue
			row["placeholder"] = optionType.PlaceHolder
			row["help_block"] = optionType.HelpBlock
			row["required"] = optionType.Required
			row["export_meta"] = optionType.ExportMeta
			row["display_value_on_details"] = optionType.DisplayValueOnDetails
			row["locked"] = optionType.IsLocked
			row["hidden"] = optionType.IsHidden
			row["exclude_from_search"] = optionType.ExcludeFromSearch

			optionTypes = append(optionTypes, row)
		}
	}
	d.Set("option_type", optionTypes)

	// Field Groups
	var fieldGroups []map[string]interface{}
	if len(form.FieldGroups) != 0 {
		for _, fieldGroup := range form.FieldGroups {
			row := make(map[string]interface{})
			row["name"] = fieldGroup.Name
			row["description"] = fieldGroup.Description
			row["collapsible"] = fieldGroup.Collapsible
			row["collapsed_by_deafult"] = fieldGroup.DefaultCollapsed
			row["visibility_field"] = fieldGroup.VisibleOnCode
			fieldGroups = append(fieldGroups, row)
		}
	}
	d.Set("field_group", fieldGroups)

	return diags
}

func resourceFormUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	name := d.Get("name").(string)

	// fieldGroups
	var fieldGroups []map[string]interface{}
	if d.Get("field_group") != nil {
		fieldGroupList := d.Get("field_group").([]interface{})
		// iterate over the array of fieldGroups
		for i := 0; i < len(fieldGroupList); i++ {
			row := make(map[string]interface{})
			fieldGroupConfig := fieldGroupList[i].(map[string]interface{})
			row["name"] = fieldGroupConfig["name"]
			row["description"] = fieldGroupConfig["description"]
			row["collapsible"] = fieldGroupConfig["collapsible"]
			row["defaultCollapsed"] = fieldGroupConfig["collapsed_by_deafult"]
			row["visibleOnCode"] = fieldGroupConfig["visibility_field"]
			fieldGroups = append(fieldGroups, row)
		}
	}

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		for _, s := range attr.(*schema.Set).List() {
			labelsPayload = append(labelsPayload, s.(string))
		}
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"optionTypeForm": map[string]interface{}{
				"name":        name,
				"code":        d.Get("code").(string),
				"description": d.Get("description").(string),
				"labels":      labelsPayload,
				"fieldGroups": fieldGroups,
			},
		},
	}

	resp, err := client.UpdateForm(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.UpdateFormResult)
	formResult := result.Form
	// Successfully created resource, now set id
	d.SetId(int64ToString(formResult.ID))
	return resourceFormRead(ctx, d, meta)
}

func resourceFormDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteForm(toInt64(id), req)
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
