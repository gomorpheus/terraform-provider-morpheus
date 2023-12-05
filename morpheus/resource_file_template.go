package morpheus

import (
	"context"
	"strings"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceFileTemplate() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus file template resource",
		CreateContext: resourceFileTemplateCreate,
		ReadContext:   resourceFileTemplateRead,
		UpdateContext: resourceFileTemplateUpdate,
		DeleteContext: resourceFileTemplateDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the file template",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the file template",
				Required:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "The organization labels associated with the file template (Only supported on Morpheus 5.5.3 or higher)",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"file_name": {
				Type:        schema.TypeString,
				Description: "The name of the file deployed by the file template",
				Required:    true,
			},
			"file_path": {
				Type:        schema.TypeString,
				Description: "The system path of the file deployed by the file template",
				Optional:    true,
			},
			"phase": {
				Type:         schema.TypeString,
				Description:  "The phase that the file template should be run during (preProvision, provision, postProvision, preDeploy, deploy)",
				ValidateFunc: validation.StringInSlice([]string{"preProvision", "provision", "postProvision", "preDeploy", "deploy"}, false),
				Required:     true,
			},
			"file_content": {
				Type:        schema.TypeString,
				Description: "The content of the file template",
				Optional:    true,
				StateFunc: func(v interface{}) string {
					payload := strings.TrimSuffix(v.(string), "\n")
					return payload
				},
			},
			"file_owner": {
				Type:        schema.TypeString,
				Description: "The file template file owner",
				Optional:    true,
			},
			"setting_name": {
				Type:        schema.TypeString,
				Description: "The file template setting name",
				Optional:    true,
			},
			"setting_category": {
				Type:        schema.TypeString,
				Description: "The file template setting category",
				Optional:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceFileTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		for _, s := range attr.(*schema.Set).List() {
			labelsPayload = append(labelsPayload, s.(string))
		}
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"containerTemplate": map[string]interface{}{
				"name":            name,
				"labels":          labelsPayload,
				"fileName":        d.Get("file_name").(string),
				"filePath":        d.Get("file_path").(string),
				"templatePhase":   d.Get("phase").(string),
				"template":        d.Get("file_content").(string),
				"fileOwner":       d.Get("file_owner").(string),
				"settingName":     d.Get("setting_name").(string),
				"settingCategory": d.Get("setting_category").(string),
			},
		},
	}

	resp, err := client.CreateFileTemplate(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateFileTemplateResult)
	fileTemplate := result.FileTemplate
	// Successfully created resource, now set id
	d.SetId(int64ToString(fileTemplate.ID))

	resourceFileTemplateRead(ctx, d, meta)
	return diags
}

func resourceFileTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindFileTemplateByName(name)
	} else if id != "" {
		resp, err = client.GetFileTemplate(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("File template cannot be read without name or id")
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
	result := resp.Result.(*morpheus.GetFileTemplateResult)
	fileTemplate := result.FileTemplate
	d.SetId(int64ToString(fileTemplate.ID))
	d.Set("name", fileTemplate.Name)
	d.Set("labels", fileTemplate.Labels)
	d.Set("file_name", fileTemplate.FileName)
	d.Set("file_path", fileTemplate.FilePath)
	d.Set("phase", fileTemplate.TemplatePhase)
	d.Set("file_content", fileTemplate.Template)
	d.Set("file_owner", fileTemplate.FileOwner)
	d.Set("setting_name", fileTemplate.SettingName)
	d.Set("setting_category", fileTemplate.SettingCategory)
	return diags
}

func resourceFileTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	name := d.Get("name").(string)

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		for _, s := range attr.(*schema.Set).List() {
			labelsPayload = append(labelsPayload, s.(string))
		}
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"containerTemplate": map[string]interface{}{
				"name":            name,
				"labels":          labelsPayload,
				"fileName":        d.Get("file_name").(string),
				"filePath":        d.Get("file_path").(string),
				"templatePhase":   d.Get("phase").(string),
				"template":        d.Get("file_content").(string),
				"fileOwner":       d.Get("file_owner").(string),
				"settingName":     d.Get("setting_name").(string),
				"settingCategory": d.Get("setting_category").(string),
			},
		},
	}

	resp, err := client.UpdateFileTemplate(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateFileTemplateResult)
	fileTemplate := result.FileTemplate
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(fileTemplate.ID))
	return resourceFileTemplateRead(ctx, d, meta)
}

func resourceFileTemplateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteFileTemplate(toInt64(id), req)
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
