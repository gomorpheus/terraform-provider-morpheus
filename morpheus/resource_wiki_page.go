package morpheus

import (
	"context"
	"strings"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceWikiPage() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus wiki page resource",
		CreateContext: resourceWikiPageCreate,
		ReadContext:   resourceWikiPageRead,
		UpdateContext: resourceWikiPageUpdate,
		DeleteContext: resourceWikiPageDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the wiki page",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the wiki page",
				Required:    true,
			},
			"category": {
				Type:        schema.TypeString,
				Description: "The category of the wiki page",
				Optional:    true,
				Computed:    true,
			},
			"content": {
				Type:        schema.TypeString,
				Description: "The content of the wiki page",
				Optional:    true,
				StateFunc: func(v interface{}) string {
					payload := strings.TrimSuffix(v.(string), "\n")
					return payload
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceWikiPageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	wikiPage := make(map[string]interface{})

	wikiPage["name"] = d.Get("name").(string)
	wikiPage["category"] = d.Get("category").(string)
	wikiPage["content"] = d.Get("content").(string)

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"page": wikiPage,
		},
	}
	resp, err := client.CreateWiki(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateWikiResult)
	wikiPageResult := result.Wiki
	// Successfully created resource, now set id
	d.SetId(int64ToString(wikiPageResult.ID))

	resourceWikiPageRead(ctx, d, meta)
	return diags
}

func resourceWikiPageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindWikiByName(name)
	} else if id != "" {
		resp, err = client.GetWiki(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Wiki Page cannot be read without name or id")
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
	result := resp.Result.(*morpheus.GetWikiResult)
	wikiPage := result.Wiki

	d.SetId(intToString(int(wikiPage.ID)))
	d.Set("name", wikiPage.Name)
	d.Set("category", wikiPage.Category)
	d.Set("content", wikiPage.Content)

	return diags
}

func resourceWikiPageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	wikiPage := make(map[string]interface{})

	wikiPage["name"] = d.Get("name").(string)
	wikiPage["category"] = d.Get("category").(string)
	wikiPage["content"] = d.Get("content").(string)

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"page": wikiPage,
		},
	}
	resp, err := client.UpdateWiki(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateWikiResult)
	wikiPageResult := result.Wiki

	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(wikiPageResult.ID))
	return resourceWikiPageRead(ctx, d, meta)
}

func resourceWikiPageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteWiki(toInt64(id), req)
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
