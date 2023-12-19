package morpheus

import (
	"context"

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
			"option_types": {
				Type:        schema.TypeList,
				Description: "List of option type json payloads",
				Optional:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsJSON,
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
