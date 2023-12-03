package morpheus

import (
	"context"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceClusterPackage() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus cluster package resource.",
		CreateContext: resourceClusterPackageCreate,
		ReadContext:   resourceClusterPackageRead,
		UpdateContext: resourceClusterPackageUpdate,
		DeleteContext: resourceClusterPackageDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the cluster package",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the cluster package",
				Required:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code for the cluster package",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the cluster package",
				Optional:    true,
			},
			"package_version": {
				Type:        schema.TypeString,
				Description: "The version of the cluster package",
				Required:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "The package category type (apps, custom, ingress, logging, monitoring, morpheus, network, serviceMesh, storage)",
				Required:    true,
			},
			"package_type": {
				Type:        schema.TypeString,
				Description: "A one word descriptor of package, such as calico, rook, prometheus, etc.",
				Required:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the cluster package is enabled",
				Optional:    true,
				Computed:    true,
			},
			"repeat_install": {
				Type:        schema.TypeBool,
				Description: "Whether to support the reinstallation of the package",
				Optional:    true,
				Computed:    true,
			},
			"spec_template_ids": {
				Type:        schema.TypeList,
				Description: "A list of spec template ids associated with the cluster package",
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Optional:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceClusterPackageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	req := &morpheus.Request{
		Method: "POST",
		Path:   morpheus.ClusterPackagesPath,
		Body: map[string]interface{}{
			"clusterPackage": map[string]interface{}{
				"name":           d.Get("name").(string),
				"code":           d.Get("code").(string),
				"description":    d.Get("description").(string),
				"enabled":        d.Get("enabled").(bool),
				"repeatInstall":  d.Get("repeat_install").(bool),
				"type":           d.Get("type").(string),
				"packageType":    d.Get("package_type").(string),
				"packageVersion": d.Get("package_version").(string),
				"specTemplates":  d.Get("spec_template_ids"),
			},
		},
		Result: &ClusterPackageCreateResult{},
	}

	resp, err := client.Execute(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*ClusterPackageCreateResult)
	// Successfully created resource, now set id
	d.SetId(int64ToString(result.ID))

	resourceClusterPackageRead(ctx, d, meta)
	return diags
}

func resourceClusterPackageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindClusterPackageByName(name)
	} else if id != "" {
		resp, err = client.GetClusterPackage(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Cluster Package cannot be read without name or id")
	}

	if err != nil {
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
	result := resp.Result.(*morpheus.GetClusterPackageResult)
	clusterPackage := result.ClusterPackage
	if clusterPackage != nil {
		d.SetId(int64ToString(clusterPackage.ID))
		d.Set("name", clusterPackage.Name)
		d.Set("description", clusterPackage.Description)
		d.Set("code", clusterPackage.Code)
		d.Set("package_version", clusterPackage.PackageVersion)
		d.Set("type", clusterPackage.Type)
		d.Set("package_type", clusterPackage.PackageType)
		d.Set("enabled", clusterPackage.Enabled)
		d.Set("repeat_install", clusterPackage.RepeatInstall)
		// spec templates
		var specTemplates []int64
		if clusterPackage.SpecTemplates != nil {
			// iterate over the array of spec templates
			for i := 0; i < len(clusterPackage.SpecTemplates); i++ {
				specTemplate := clusterPackage.SpecTemplates[i]
				specTemplates = append(specTemplates, int64(specTemplate.ID))
			}
		}
		d.Set("spec_template_ids", specTemplates)
	} else {
		return diag.Errorf("read operation: cluster package not found in response data") // should not happen
	}

	return diags
}

func resourceClusterPackageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	name := d.Get("name").(string)
	req := &morpheus.Request{
		Body: map[string]interface{}{
			"clusterPackage": map[string]interface{}{
				"name":           name,
				"code":           d.Get("code").(string),
				"description":    d.Get("description").(string),
				"enabled":        d.Get("enabled").(bool),
				"repeatInstall":  d.Get("repeat_install").(bool),
				"type":           d.Get("type").(string),
				"packageType":    d.Get("package_type").(string),
				"packageVersion": d.Get("package_version").(string),
				"specTemplates":  d.Get("spec_template_ids"),
			},
		},
	}

	resp, err := client.UpdateClusterPackage(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	return resourceClusterPackageRead(ctx, d, meta)
}

func resourceClusterPackageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteClusterPackage(toInt64(id), req)
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

type ClusterPackageCreateResult struct {
	ID      int64             `json:"id"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	Success bool              `json:"success"`
}
