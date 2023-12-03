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

func resourceKubernetesAppBlueprint() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus kubernetes app blueprint resource",
		CreateContext: resourceKubernetesAppBlueprintCreate,
		ReadContext:   resourceKubernetesAppBlueprintRead,
		UpdateContext: resourceKubernetesAppBlueprintUpdate,
		DeleteContext: resourceKubernetesAppBlueprintDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the kubernetes app blueprint",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the kubernetes app blueprint",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the kubernetes app blueprint",
				Optional:    true,
			},
			"category": {
				Type:        schema.TypeString,
				Description: "The category of the kubernetes app blueprint",
				Optional:    true,
			},
			"source_type": {
				Type:         schema.TypeString,
				Description:  "The source of the kubernetes app blueprint (yaml, spec or repository)",
				ValidateFunc: validation.StringInSlice([]string{"yaml", "spec", "repository"}, false),
				Required:     true,
			},
			"blueprint_content": {
				Type:        schema.TypeString,
				Description: "The content of the kubernetes app blueprint. Used when the yaml source type is specified",
				Optional:    true,
			},
			"working_path": {
				Type:        schema.TypeString,
				Description: "The path of the kubernetes app blueprint in the git repository",
				Optional:    true,
			},
			"integration_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the git integration",
				Optional:    true,
			},
			"repository_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the git repository",
				Optional:    true,
			},
			"version_ref": {
				Type:        schema.TypeString,
				Description: "The git reference of the repository to pull (main, master, etc.)",
				Optional:    true,
			},
			"spec_template_ids": {
				Type:        schema.TypeList,
				Description: "A list of kubernetes spec template ids associated with the app blueprint",
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Optional:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceKubernetesAppBlueprintCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	blueprint_type := "kubernetes"
	description := d.Get("description").(string)
	category := d.Get("category").(string)

	config := make(map[string]interface{})
	config["name"] = name
	config["description"] = description
	config["category"] = category
	config["type"] = "kubernetes"

	kubernetesConfig := make(map[string]interface{})
	config["kubernetes"] = kubernetesConfig

	switch d.Get("source_type") {
	case "yaml":
		kubernetesConfig["configType"] = "yaml"
		kubernetesConfig["yaml"] = d.Get("blueprint_content").(string)

	case "spec":
		kubernetesConfig["configType"] = "spec"
		var spec_templates []map[string]interface{}
		if d.Get("spec_template_ids") != nil {
			specTemplateList := d.Get("spec_template_ids").([]interface{})
			// iterate over the array of spec templates
			for i := 0; i < len(specTemplateList); i++ {
				row := make(map[string]interface{})
				row["id"] = specTemplateList[i]
				row["value"] = specTemplateList[i]
				spec_templates = append(spec_templates, row)
			}
		}
		specConfig := make(map[string]interface{})
		config["config"] = specConfig
		specConfig["specs"] = spec_templates
	case "repository":
		kubernetesConfig["configType"] = "git"
		kubernetesGitConfig := make(map[string]interface{})
		kubernetesGitConfig["integrationId"] = d.Get("integration_id")
		kubernetesGitConfig["repoId"] = d.Get("repository_id")
		kubernetesGitConfig["branch"] = d.Get("version_ref").(string)
		kubernetesGitConfig["path"] = d.Get("working_path").(string)
		kubernetesConfig["git"] = kubernetesGitConfig
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"blueprint": map[string]interface{}{
				"name":        name,
				"type":        blueprint_type,
				"description": description,
				"category":    category,
				"config":      config,
			},
		},
	}

	resp, err := client.CreateBlueprint(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateBlueprintResult)
	blueprint := result.Blueprint
	// Successfully created resource, now set id
	d.SetId(int64ToString(blueprint.ID))

	resourceKubernetesAppBlueprintRead(ctx, d, meta)
	return diags
}

func resourceKubernetesAppBlueprintRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindBlueprintByName(name)
	} else if id != "" {
		resp, err = client.GetBlueprint(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Blueprint cannot be read without name or id")
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
	var kubernetesBlueprint KubernetesAppBlueprint
	json.Unmarshal(resp.Body, &kubernetesBlueprint)
	d.SetId(intToString(kubernetesBlueprint.Blueprint.ID))
	d.Set("name", kubernetesBlueprint.Blueprint.Name)
	d.Set("description", kubernetesBlueprint.Blueprint.Description)
	d.Set("category", kubernetesBlueprint.Blueprint.Category)

	switch kubernetesBlueprint.Blueprint.Config.Kubernetes.Configtype {
	case "yaml":
		d.Set("source_type", "yaml")
		d.Set("blueprint_content", kubernetesBlueprint.Blueprint.Config.Kubernetes)
	case "git":
		d.Set("source_type", "repository")
		d.Set("working_path", kubernetesBlueprint.Blueprint.Config.Kubernetes.Git.Path)
		d.Set("integration_id", kubernetesBlueprint.Blueprint.Config.Kubernetes.Git.IntegrationId)
		d.Set("repository_id", kubernetesBlueprint.Blueprint.Config.Kubernetes.Git.RepoId)
		d.Set("version_ref", kubernetesBlueprint.Blueprint.Config.Kubernetes.Git.Branch)
	case "spec":
		d.Set("source_type", "spec")
		// spec templates
		var specTemplates []int64
		if kubernetesBlueprint.Blueprint.Config.Config.Specs != nil {
			// iterate over the array of tasks
			for i := 0; i < len(kubernetesBlueprint.Blueprint.Config.Config.Specs); i++ {
				specTemplate := kubernetesBlueprint.Blueprint.Config.Config.Specs[i]
				specTemplates = append(specTemplates, int64(specTemplate.ID))
			}
		}
		d.Set("spec_templates_ids", specTemplates)
	}

	return diags
}

func resourceKubernetesAppBlueprintUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	name := d.Get("name").(string)
	blueprint_type := "kubernetes"
	description := d.Get("description").(string)
	category := d.Get("category").(string)

	config := make(map[string]interface{})
	config["name"] = name
	config["description"] = description
	config["category"] = category
	config["type"] = "kubernetes"

	kubernetesConfig := make(map[string]interface{})
	config["kubernetes"] = kubernetesConfig

	switch d.Get("source_type") {
	case "yaml":
		kubernetesConfig["configType"] = "yaml"
		kubernetesConfig["yaml"] = d.Get("blueprint_content").(string)

	case "spec":
		kubernetesConfig["configType"] = "spec"
		var spec_templates []map[string]interface{}
		if d.Get("spec_template_ids") != nil {
			specTemplateList := d.Get("spec_template_ids").([]interface{})
			// iterate over the array of spec templates
			for i := 0; i < len(specTemplateList); i++ {
				row := make(map[string]interface{})
				row["id"] = specTemplateList[i]
				row["value"] = specTemplateList[i]
				spec_templates = append(spec_templates, row)
			}
		}
		specConfig := make(map[string]interface{})
		config["config"] = specConfig
		specConfig["specs"] = spec_templates
	case "repository":
		kubernetesConfig["configType"] = "git"
		kubernetesGitConfig := make(map[string]interface{})
		kubernetesGitConfig["integrationId"] = d.Get("integration_id")
		kubernetesGitConfig["repoId"] = d.Get("repository_id")
		kubernetesGitConfig["branch"] = d.Get("version_ref").(string)
		kubernetesGitConfig["path"] = d.Get("working_path").(string)
		kubernetesConfig["git"] = kubernetesGitConfig
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"blueprint": map[string]interface{}{
				"name":        name,
				"type":        blueprint_type,
				"description": description,
				"category":    category,
				"config":      config,
			},
		},
	}
	log.Printf("API REQUEST: %s", req)
	resp, err := client.UpdateBlueprint(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateBlueprintResult)
	blueprint := result.Blueprint
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(blueprint.ID))
	return resourceKubernetesAppBlueprintRead(ctx, d, meta)
}

func resourceKubernetesAppBlueprintDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteBlueprint(toInt64(id), req)
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

type KubernetesAppBlueprint struct {
	Blueprint struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Type        string `json:"type"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Config      struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Kubernetes  struct {
				Configtype string `json:"configType"`
				Git        struct {
					Path          string `json:"path"`
					RepoId        int    `json:"repoId"`
					IntegrationId int    `json:"integrationId"`
					Branch        string `json:"branch"`
				} `json:"git"`
			} `json:"kubernetes"`
			Config struct {
				Specs []struct {
					ID    int    `json:"id"`
					Value string `json:"value"`
					Name  string `json:"name"`
				} `json:"specs"`
			} `json:"config"`
			Type     string `json:"type"`
			Category string `json:"category"`
			Image    string `json:"image"`
		} `json:"config"`
		Visibility         string `json:"visibility"`
		Resourcepermission struct {
			All      bool          `json:"all"`
			Sites    []interface{} `json:"sites"`
			AllPlans bool          `json:"allPlans"`
			Plans    []interface{} `json:"plans"`
		} `json:"resourcePermission"`
		Owner struct {
			ID       int    `json:"id"`
			Username string `json:"username"`
		} `json:"owner"`
		Tenant struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"tenant"`
	} `json:"blueprint"`
}
