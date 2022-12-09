package morpheus

import (
	"context"
	"strings"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBootScript() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus boot script resource",
		CreateContext: resourceBootScriptCreate,
		ReadContext:   resourceBootScriptRead,
		UpdateContext: resourceBootScriptUpdate,
		DeleteContext: resourceBootScriptDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the boot script",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the boot script",
				Required:    true,
			},
			"content": {
				Type:        schema.TypeString,
				Description: "The content of the boot script",
				Optional:    true,
				StateFunc: func(v interface{}) string {
					payload := strings.TrimSuffix(v.(string), "\n")
					return payload
				},
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceBootScriptCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"bootScript": map[string]interface{}{
				"fileName": d.Get("name").(string),
				"content":  d.Get("content").(string),
			},
		},
	}

	resp, err := client.CreateBootScript(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateBootScriptResult)
	bootScript := result.BootScript
	// Successfully created resource, now set id
	d.SetId(int64ToString(bootScript.ID))

	resourceBootScriptRead(ctx, d, meta)
	return diags
}

func resourceBootScriptRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindBootScriptByName(name)
	} else if id != "" {
		resp, err = client.GetBootScript(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("File template cannot be read without name or id")
	}

	if err != nil {
		// 404 is ok?
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
	result := resp.Result.(*morpheus.GetBootScriptResult)
	bootScript := result.BootScript
	d.SetId(int64ToString(bootScript.ID))
	d.Set("name", bootScript.FileName)
	d.Set("content", bootScript.Content)
	return diags
}

func resourceBootScriptUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"bootScript": map[string]interface{}{
				"fileName": d.Get("name").(string),
				"content":  d.Get("content").(string),
			},
		},
	}

	resp, err := client.UpdateBootScript(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateBootScriptResult)
	bootScript := result.BootScript
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(bootScript.ID))
	return resourceBootScriptRead(ctx, d, meta)
}

func resourceBootScriptDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteBootScript(toInt64(id), req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
			return nil
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)
			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)
	d.SetId("")
	return diags
}
