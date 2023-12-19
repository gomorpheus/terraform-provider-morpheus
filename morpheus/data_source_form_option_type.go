package morpheus

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMorpheusFormOptionType() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus form option type data source.",
		ReadContext: dataSourceMorpheusFormOptionTypeRead,
		Schema: map[string]*schema.Schema{
			"json": {
				Type:        schema.TypeString,
				Description: "JSON form option type rendered based on the arguments defined",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the option type",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the option type",
				Optional:    true,
				Computed:    true,
			},
			"field_name": {
				Type:        schema.TypeString,
				Description: "The field name of the option type",
				Optional:    true,
				Computed:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "The type of the option type",
				Required:    true,
			},
			"field_label": {
				Type:        schema.TypeString,
				Description: "The field label of the option type",
				Optional:    true,
				Computed:    true,
			},
			"localized_label": {
				Type:        schema.TypeString,
				Description: "The field label of the option type",
				Optional:    true,
				Computed:    true,
			},
			"default_value": {
				Type:        schema.TypeString,
				Description: "The default value of the option type",
				Optional:    true,
				Computed:    true,
			},
			"placeholder": {
				Type:        schema.TypeString,
				Description: "The field label of the option type",
				Optional:    true,
				Computed:    true,
			},
			"help_block": {
				Type:        schema.TypeString,
				Description: "The field label of the option type",
				Optional:    true,
				Computed:    true,
			},
			"localized_help_block": {
				Type:        schema.TypeString,
				Description: "The field label of the option type",
				Optional:    true,
				Computed:    true,
			},
			"required": {
				Type:        schema.TypeBool,
				Description: "Whether the option type is required",
				Optional:    true,
				Default:     false,
			},
			"export_meta": {
				Type:        schema.TypeBool,
				Description: "Whether to export the number option type as a tag",
				Optional:    true,
				Default:     false,
			},
			"display_value_on_details": {
				Type:        schema.TypeBool,
				Description: "Display the selected value of the number option type on the associated resource's details page",
				Optional:    true,
				Default:     false,
			},
			"locked": {
				Type:        schema.TypeBool,
				Description: "Display the selected value of the number option type on the associated resource's details page",
				Optional:    true,
				Default:     false,
			},
			"hidden": {
				Type:        schema.TypeBool,
				Description: "Display the selected value of the number option type on the associated resource's details page",
				Optional:    true,
				Default:     false,
			},
			"exclude_from_search": {
				Type:        schema.TypeBool,
				Description: "Display the selected value of the number option type on the associated resource's details page",
				Optional:    true,
				Default:     false,
			},
			"allow_password_peek": {
				Type:        schema.TypeBool,
				Description: "Display the selected value of the number option type on the associated resource's details page",
				Optional:    true,
				Default:     false,
			},
			"row": {
				Type:        schema.TypeString,
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
				Description: "The field or code used to trigger the requirement of this field",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func dataSourceMorpheusFormOptionTypeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var formData FormOptionType
	formData.Name = d.Get("default_group_permission").(string)
	formData.Description = d.Get("description").(string)
	formData.FieldName = d.Get("field_name").(string)
	formData.FieldLabel = d.Get("field_label").(string)
	formData.Type = d.Get("type").(string)

	jsonDoc, err := json.MarshalIndent(formData, "", "  ")
	log.Printf("API RESPONSE: %s", jsonDoc)

	if err != nil {
		return diag.Errorf("writing permission set: formatting JSON: %s", err)
	}
	jsonString := string(jsonDoc)

	d.Set("json", jsonString)
	d.SetId(strconv.Itoa(1))
	return diags
}

type FormOptionType struct {
	Name                  string `json:"name,omitempty"`
	Description           string `json:"description,omitempty"`
	FieldName             string `json:"field_name,omitempty"`
	FieldLabel            string `json:"field_label,omitempty"`
	Type                  string `json:"type,omitempty"`
	DisplayValueOnDetails bool   `json:"display_value_on_details,omitempty"`
	Locked                bool   `json:"locked,omitempty"`
	Hidden                bool   `json:"hidden,omitempty"`
	ExcludeFromSearch     bool   `json:"exclude_from_search,omitempty"`
	AllowPasswordPeak     bool   `json:"allow_password_peek,omitempty"`
	Row                   string `json:"row,omitempty"`
	DependentField        string `json:"dependent_field,omitempty"`
	VisibilityField       string `json:"visibility_field,omitempty"`
	VerifyPattern         string `json:"verify_pattern,omitempty"`
	RequireField          string `json:"require_field,omitempty"`
}
