package morpheus

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceTerraformSpecTemplate() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus terraform spec template resource",
		CreateContext: resourceTerraformSpecTemplateCreate,
		ReadContext:   resourceTerraformSpecTemplateRead,
		UpdateContext: resourceTerraformSpecTemplateUpdate,
		DeleteContext: resourceTerraformSpecTemplateDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the terraform spec template",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the terraform spec template",
				Required:    true,
			},
			"source_type": {
				Type:         schema.TypeString,
				Description:  "The source of the terraform spec template (local, url or repository)",
				ValidateFunc: validation.StringInSlice([]string{"local", "url", "repository"}, false),
				Required:     true,
			},
			"spec_content": {
				Type:        schema.TypeString,
				Description: "The content of the terraform spec template. Used when the local source type is specified",
				Optional:    true,
				Computed:    true,
				StateFunc: func(val interface{}) string {
					return strings.TrimSuffix(val.(string), "\n")
				},
			},
			"spec_path": {
				Type:        schema.TypeString,
				Description: "The path of the terraform spec template, either the url or the path in the repository",
				Optional:    true,
				Computed:    true,
			},
			"repository_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the git repository integration",
				Optional:    true,
				Computed:    true,
			},
			"version_ref": {
				Type:        schema.TypeString,
				Description: "The git reference of the repository to pull (main, master, etc.)",
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTerraformSpecTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)

	sourceOptions := make(map[string]interface{})
	if d.Get("spec_content") != "" {
		sourceOptions["content"] = d.Get("spec_content")
	}
	if d.Get("spec_path") != "" {
		sourceOptions["contentPath"] = d.Get("spec_path")
	}
	sourceOptions["contentRef"] = d.Get("version_ref")
	sourceOptions["repository"] = map[string]interface{}{
		"id": d.Get("repository_id"),
	}
	sourceOptions["sourceType"] = d.Get("source_type")

	specTemplateType := make(map[string]interface{})
	specTemplateType["code"] = "terraform"

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"specTemplate": map[string]interface{}{
				"name": name,
				"file": sourceOptions,
				"type": specTemplateType,
			},
		},
	}
	resp, err := client.CreateSpecTemplate(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateSpecTemplateResult)
	specTemplate := result.SpecTemplate
	// Successfully created resource, now set id
	d.SetId(int64ToString(specTemplate.ID))

	resourceTerraformSpecTemplateRead(ctx, d, meta)
	return diags
}

func resourceTerraformSpecTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindSpecTemplateByName(name)
	} else if id != "" {
		resp, err = client.GetSpecTemplate(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Spec template cannot be read without name or id")
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
	var terraformSpecTemplate TerraformSpecTemplate
	json.Unmarshal(resp.Body, &terraformSpecTemplate)
	d.SetId(intToString(terraformSpecTemplate.Spectemplate.ID))
	d.Set("name", terraformSpecTemplate.Spectemplate.Name)
	d.Set("source_type", terraformSpecTemplate.Spectemplate.File.Sourcetype)

	switch terraformSpecTemplate.Spectemplate.File.Sourcetype {
	case "local":
		d.Set("source_type", "local")
		d.Set("spec_content", terraformSpecTemplate.Spectemplate.File.Content)
	case "url":
		d.Set("source_type", "url")
		d.Set("spec_path", terraformSpecTemplate.Spectemplate.File.Contentpath)
	case "git":
		d.Set("source_type", "repository")
		d.Set("spec_path", terraformSpecTemplate.Spectemplate.File.Contentpath)
		d.Set("repository_id", terraformSpecTemplate.Spectemplate.File.Repository.ID)
		d.Set("version_ref", terraformSpecTemplate.Spectemplate.File.Contentref)
	}
	return diags
}

func resourceTerraformSpecTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	name := d.Get("name").(string)

	sourceOptions := make(map[string]interface{})
	if d.Get("spec_content") != "" {
		sourceOptions["content"] = d.Get("spec_content")
	}
	if d.Get("spec_path") != "" {
		sourceOptions["contentPath"] = d.Get("spec_path")
	}
	sourceOptions["contentRef"] = d.Get("version_ref")
	sourceOptions["repository"] = map[string]interface{}{
		"id": d.Get("repository_id"),
	}
	sourceOptions["sourceType"] = d.Get("source_type")

	specTemplateType := make(map[string]interface{})
	specTemplateType["code"] = "terraform"

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"specTemplate": map[string]interface{}{
				"name": name,
				"file": sourceOptions,
				"type": specTemplateType,
			},
		},
	}
	resp, err := client.UpdateSpecTemplate(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateSpecTemplateResult)
	specTemplate := result.SpecTemplate
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(specTemplate.ID))
	return resourceTerraformSpecTemplateRead(ctx, d, meta)
}

func resourceTerraformSpecTemplateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteSpecTemplate(toInt64(id), req)
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

type TerraformSpecTemplate struct {
	Spectemplate struct {
		ID      int `json:"id"`
		Account struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"account"`
		Name string      `json:"name"`
		Code interface{} `json:"code"`
		Type struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"type"`
		Externalid   interface{} `json:"externalId"`
		Externaltype interface{} `json:"externalType"`
		Deploymentid interface{} `json:"deploymentId"`
		Status       interface{} `json:"status"`
		File         struct {
			ID          int         `json:"id"`
			Sourcetype  string      `json:"sourceType"`
			Contentref  interface{} `json:"contentRef"`
			Contentpath interface{} `json:"contentPath"`
			Repository  struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"repository"`
			Content string `json:"content"`
		} `json:"file"`
		Config struct {
		} `json:"config"`
		Createdby   string      `json:"createdBy"`
		Updatedby   interface{} `json:"updatedBy"`
		Datecreated time.Time   `json:"dateCreated"`
		Lastupdated time.Time   `json:"lastUpdated"`
	} `json:"specTemplate"`
}
