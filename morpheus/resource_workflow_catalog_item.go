package morpheus

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceWorkflowCatalogItem() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus workflow catalog item resource",
		CreateContext: resourceWorkflowCatalogItemCreate,
		ReadContext:   resourceWorkflowCatalogItemRead,
		UpdateContext: resourceWorkflowCatalogItemUpdate,
		DeleteContext: resourceWorkflowCatalogItemDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the workflow catalog item",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the workflow catalog item",
				Required:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "The organization labels associated with the catalog item (Only supported on Morpheus 5.5.3 or higher)",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the workflow catalog item",
				Optional:    true,
				Computed:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the workflow catalog item is enabled",
				Optional:    true,
				Default:     true,
			},
			"featured": {
				Type:        schema.TypeBool,
				Description: "Whether the workflow catalog item is featured",
				Optional:    true,
				Computed:    true,
			},
			"workflow_id": {
				Type:        schema.TypeInt,
				Description: "The id of the workflow associated with the workflow catalog item",
				Required:    true,
			},
			"context_type": {
				Type:         schema.TypeString,
				Description:  "The Morpheus context type of the operational workflow",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"instance", "server", "appliance"}, false),
				Computed:     true,
			},
			"content": {
				Type:        schema.TypeString,
				Description: "The markdown content associated with the workflow catalog item",
				Optional:    true,
				Computed:    true,
			},
			"option_type_ids": {
				Type:        schema.TypeList,
				Description: "The list of option type ids associated with the workflow catalog item",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == old
				},
				Computed: true,
			},
			"logo_image_name": {
				Type:        schema.TypeString,
				Description: "The file name of the workflow catalog item logo image",
				Optional:    true,
				Computed:    true,
			},
			"logo_image_path": {
				Type:        schema.TypeString,
				Description: "The file path of the workflow catalog item logo image including the file name",
				Optional:    true,
				Computed:    true,
			},
			"dark_logo_image_name": {
				Type:        schema.TypeString,
				Description: "The file name of the workflow catalog item dark mode logo image",
				Optional:    true,
				Computed:    true,
			},
			"dark_logo_image_path": {
				Type:        schema.TypeString,
				Description: "The file path of the workflow catalog item dark mode logo image including the file name",
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceWorkflowCatalogItemCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	catalogItem := make(map[string]interface{})

	catalogItem["name"] = d.Get("name").(string)
	catalogItem["description"] = d.Get("description").(string)
	catalogItem["enabled"] = d.Get("enabled").(bool)
	catalogItem["featured"] = d.Get("featured").(bool)
	catalogItem["type"] = "workflow"
	catalogItem["iconPath"] = "custom"
	catalogItem["context"] = d.Get("context_type").(string)
	catalogItem["optionTypes"] = d.Get("option_type_ids")
	catalogItem["content"] = d.Get("content").(string)

	catalogItem["workflow"] = map[string]interface{}{
		"id": d.Get("workflow_id").(int),
	}

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		for _, s := range attr.(*schema.Set).List() {
			labelsPayload = append(labelsPayload, s.(string))
		}
	}
	catalogItem["labels"] = labelsPayload

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"catalogItemType": catalogItem,
		},
	}
	resp, err := client.CreateCatalogItem(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateCatalogItemResult)
	catalogItemResult := result.CatalogItem

	var filePayloads []*morpheus.FilePayload

	if d.Get("logo_image_path") != "" && d.Get("logo_image_name") != "" {
		data, err := os.ReadFile(d.Get("logo_image_path").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		filePayload := &morpheus.FilePayload{
			ParameterName: "logo",
			FileName:      d.Get("logo_image_name").(string),
			FileContent:   data,
		}
		filePayloads = append(filePayloads, filePayload)
	}
	if d.Get("dark_logo_image_path") != "" && d.Get("dark_logo_image_name") != "" {
		darkLogoData, err := os.ReadFile(d.Get("dark_logo_image_path").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		darkLogoPayload := &morpheus.FilePayload{
			ParameterName: "darkLogo",
			FileName:      d.Get("dark_logo_image_name").(string),
			FileContent:   darkLogoData,
		}
		filePayloads = append(filePayloads, darkLogoPayload)
	}

	response, err := client.UpdateCatalogItemLogo(catalogItemResult.ID, filePayloads, &morpheus.Request{})
	if err != nil {
		log.Printf("API FAILURE: %s - %s", response, err)
	}
	log.Printf("API RESPONSE: %s", response)

	// Successfully created resource, now set id
	d.SetId(int64ToString(catalogItemResult.ID))

	resourceWorkflowCatalogItemRead(ctx, d, meta)
	return diags
}

func resourceWorkflowCatalogItemRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindCatalogItemByName(name)
	} else if id != "" {
		resp, err = client.GetCatalogItem(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Catalog Item cannot be read without name or id")
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
	result := resp.Result.(*morpheus.GetCatalogItemResult)
	catalogItem := result.CatalogItem

	d.SetId(intToString(int(catalogItem.ID)))
	d.Set("name", catalogItem.Name)
	d.Set("labels", catalogItem.Labels)
	d.Set("description", catalogItem.Description)
	d.Set("enabled", catalogItem.Enabled)
	d.Set("featured", catalogItem.Featured)
	// option types
	var optionTypes []int64
	if catalogItem.OptionTypes != nil {
		// iterate over the array of tasks
		for i := 0; i < len(catalogItem.OptionTypes); i++ {
			option := catalogItem.OptionTypes[i].(map[string]interface{})
			optionID := int64(option["id"].(float64))
			optionTypes = append(optionTypes, optionID)
		}
	}
	d.Set("option_type_ids", optionTypes)
	d.Set("content", catalogItem.Content)
	d.Set("context_type", catalogItem.Context)

	// Parse workflow ID
	var data map[string]interface{}
	err = json.Unmarshal([]byte(resp.Body), &data)
	if err != nil {
		panic(err)
	}
	catalogItemData := data["catalogItemType"].(map[string]interface{})
	workflowData := catalogItemData["workflow"].(map[string]interface{})
	workflowId := int(workflowData["id"].(float64))
	d.Set("workflow_id", workflowId)
	imagePath := strings.Split(catalogItem.ImagePath, "/")
	opt := strings.Replace(imagePath[len(imagePath)-1], "_original", "", 1)
	d.Set("logo_image_path", opt)
	darkImagePath := strings.Split(catalogItem.DarkImagePath, "/")
	darkOpt := strings.Replace(darkImagePath[len(darkImagePath)-1], "_original", "", 1)
	d.Set("dark_logo_image_path", darkOpt)
	return diags
}

func resourceWorkflowCatalogItemUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	catalogItem := make(map[string]interface{})

	catalogItem["name"] = d.Get("name").(string)
	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		for _, s := range attr.(*schema.Set).List() {
			labelsPayload = append(labelsPayload, s.(string))
		}
	}
	catalogItem["labels"] = labelsPayload
	catalogItem["description"] = d.Get("description").(string)
	catalogItem["enabled"] = d.Get("enabled").(bool)
	catalogItem["featured"] = d.Get("featured").(bool)
	catalogItem["type"] = "workflow"
	catalogItem["context"] = d.Get("context_type").(string)
	catalogItem["optionTypes"] = d.Get("option_type_ids")
	catalogItem["content"] = d.Get("content").(string)

	catalogItem["workflow"] = map[string]interface{}{
		"id": d.Get("workflow_id").(int),
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"catalogItemType": catalogItem,
		},
	}

	resp, err := client.UpdateCatalogItem(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateCatalogItemResult)
	catalogItemResult := result.CatalogItem

	var filePayloads []*morpheus.FilePayload

	if d.HasChange("logo_image_path") || d.HasChange("logo_image_name") {
		data, err := os.ReadFile(d.Get("logo_image_path").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		filePayload := &morpheus.FilePayload{
			ParameterName: "logo",
			FileName:      d.Get("logo_image_name").(string),
			FileContent:   data,
		}
		filePayloads = append(filePayloads, filePayload)
	}
	if d.HasChange("dark_logo_image_path") || d.HasChange("dark_logo_image_name") {
		darkLogoData, err := os.ReadFile(d.Get("dark_logo_image_path").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		darkLogoPayload := &morpheus.FilePayload{
			ParameterName: "darkLogo",
			FileName:      d.Get("dark_logo_image_name").(string),
			FileContent:   darkLogoData,
		}
		filePayloads = append(filePayloads, darkLogoPayload)
	}

	response, err := client.UpdateCatalogItemLogo(catalogItemResult.ID, filePayloads, &morpheus.Request{})
	if err != nil {
		log.Printf("API FAILURE: %s - %s", response, err)
	}
	log.Printf("API RESPONSE: %s", response)

	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(catalogItemResult.ID))
	return resourceWorkflowCatalogItemRead(ctx, d, meta)
}

func resourceWorkflowCatalogItemDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteCatalogItem(toInt64(id), req)
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
