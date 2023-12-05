package morpheus

import (
	"context"
	"strings"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePreseedScript() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus preseed script resource",
		CreateContext: resourcePreseedScriptCreate,
		ReadContext:   resourcePreseedScriptRead,
		UpdateContext: resourcePreseedScriptUpdate,
		DeleteContext: resourcePreseedScriptDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the preseed script",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the preseed script",
				Required:    true,
			},
			"content": {
				Type:        schema.TypeString,
				Description: "The content of the preseed script",
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

func resourcePreseedScriptCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"preseedScript": map[string]interface{}{
				"fileName": d.Get("name").(string),
				"content":  d.Get("content").(string),
			},
		},
	}

	resp, err := client.CreatePreseedScript(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreatePreseedScriptResult)
	preseedScript := result.PreseedScript
	// Successfully created resource, now set id
	d.SetId(int64ToString(preseedScript.ID))

	resourcePreseedScriptRead(ctx, d, meta)
	return diags
}

func resourcePreseedScriptRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindPreseedScriptByName(name)
	} else if id != "" {
		resp, err = client.GetPreseedScript(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Preseed script cannot be read without name or id")
	}

	if err != nil {
		// 404 is ok?
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
	result := resp.Result.(*morpheus.GetPreseedScriptResult)
	preseedScript := result.PreseedScript
	d.SetId(int64ToString(preseedScript.ID))
	d.Set("name", preseedScript.FileName)
	d.Set("content", preseedScript.Content)
	return diags
}

func resourcePreseedScriptUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"preseedScript": map[string]interface{}{
				"fileName": d.Get("name").(string),
				"content":  d.Get("content").(string),
			},
		},
	}

	resp, err := client.UpdatePreseedScript(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdatePreseedScriptResult)
	preseedScript := result.PreseedScript
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(preseedScript.ID))
	return resourcePreseedScriptRead(ctx, d, meta)
}

func resourcePreseedScriptDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeletePreseedScript(toInt64(id), req)
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
