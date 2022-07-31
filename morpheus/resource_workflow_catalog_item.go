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
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the workflow catalog item",
				Optional:    true,
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
			},
			"content": {
				Type:        schema.TypeString,
				Description: "The markdown content associated with the workflow catalog item",
				Optional:    true,
			},
			"option_type_ids": {
				Type:        schema.TypeList,
				Description: "The list of option type ids associated with the workflow catalog item",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
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
	resp, err := client.CreateCatalogItem(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateCatalogItemResult)
	catalogItemResult := result.CatalogItem
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
	d.Set("description", catalogItem.Description)
	d.Set("enabled", catalogItem.Enabled)
	d.Set("featured", catalogItem.Featured)
	d.Set("option_type_ids", catalogItem.OptionTypes)
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

	return diags
}

func resourceWorkflowCatalogItemUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	catalogItem := make(map[string]interface{})

	catalogItem["name"] = d.Get("name").(string)
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
	log.Printf("API REQUEST: %s", req)
	resp, err := client.UpdateCatalogItem(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateCatalogItemResult)
	catalogItemResult := result.CatalogItem

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
