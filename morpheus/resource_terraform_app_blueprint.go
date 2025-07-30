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

func resourceTerraformAppBlueprint() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus terraform app blueprint resource",
		CreateContext: resourceTerraformAppBlueprintCreate,
		ReadContext:   resourceTerraformAppBlueprintRead,
		UpdateContext: resourceTerraformAppBlueprintUpdate,
		DeleteContext: resourceTerraformAppBlueprintDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the terraform app blueprint",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the terraform app blueprint",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the terraform app blueprint",
				Optional:    true,
				Computed:    true,
			},
			"category": {
				Type:        schema.TypeString,
				Description: "The category of the terraform app blueprint",
				Optional:    true,
				Computed:    true,
			},
			"source_type": {
				Type:         schema.TypeString,
				Description:  "The source of the terraform app blueprint (hcl, json, spec or repository)",
				ValidateFunc: validation.StringInSlice([]string{"hcl", "json", "spec", "repository"}, false),
				Required:     true,
			},
			"blueprint_content": {
				Type:        schema.TypeString,
				Description: "The content of the terraform app blueprint. Used when the hcl or json source types are specified",
				Optional:    true,
				Computed:    true,
			},
			"working_path": {
				Type:          schema.TypeString,
				Description:   "The path of the terraform code in the git repository",
				Optional:      true,
				ConflictsWith: []string{"blueprint_content"},
			},
			"integration_id": {
				Type:          schema.TypeInt,
				Description:   "The ID of the git integration",
				Optional:      true,
				ConflictsWith: []string{"blueprint_content"},
			},
			"repository_id": {
				Type:          schema.TypeInt,
				Description:   "The ID of the git repository",
				Optional:      true,
				ConflictsWith: []string{"blueprint_content"},
				RequiredWith:  []string{"integration_id"},
			},
			"version_ref": {
				Type:        schema.TypeString,
				Description: "The git reference of the repository to pull (main, master, etc.)",
				Optional:    true,
				Computed:    true,
			},
			"spec_template_ids": {
				Type:        schema.TypeList,
				Description: "A list of terraform spec template ids associated with the app blueprint",
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Optional:    true,
				Computed:    true,
			},
			"terraform_version": {
				Type:        schema.TypeString,
				Description: "The terraform version associated with the app blueprint",
				Optional:    true,
				Computed:    true,
			},
			"terraform_options": {
				Type:        schema.TypeString,
				Description: "The additional terraform options to add to the app blueprint",
				Optional:    true,
				Computed:    true,
			},
			"tfvar_secret": {
				Type:        schema.TypeString,
				Description: "The name of the tfvar cypher secret to associate with the app blueprint",
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTerraformAppBlueprintCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	blueprint_type := "terraform"
	description := d.Get("description").(string)
	category := d.Get("category").(string)

	config := make(map[string]interface{})
	config["name"] = name
	config["description"] = description
	config["category"] = category
	config["type"] = "terraform"

	terraformConfig := make(map[string]interface{})
	terraformConfig["tfVersion"] = d.Get("terraform_version").(string)
	terraformConfig["commandOptions"] = d.Get("terraform_options").(string)
	terraformConfig["tfvarSecret"] = d.Get("tfvar_secret").(string)
	config["terraform"] = terraformConfig

	switch d.Get("source_type") {
	case "hcl":
		terraformConfig["configType"] = "tf"
		terraformConfig["tf"] = d.Get("blueprint_content").(string)
	case "json":
		terraformConfig["configType"] = "json"
		terraformConfig["json"] = d.Get("blueprint_content").(string)
	case "repository":
		terraformConfig["configType"] = "git"
		terraformGitConfig := make(map[string]interface{})
		terraformGitConfig["integrationId"] = d.Get("integration_id")
		terraformGitConfig["repoId"] = d.Get("repository_id")
		terraformGitConfig["branch"] = d.Get("version_ref").(string)
		terraformGitConfig["path"] = d.Get("working_path").(string)
		terraformConfig["git"] = terraformGitConfig
	case "spec":
		terraformConfig["configType"] = "spec"
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

	resourceTerraformAppBlueprintRead(ctx, d, meta)
	return diags
}

func resourceTerraformAppBlueprintRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	var terraformBlueprint TerraformAppBlueprint
	if err := json.Unmarshal(resp.Body, &terraformBlueprint); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(intToString(terraformBlueprint.Blueprint.ID))
	d.Set("name", terraformBlueprint.Blueprint.Name)
	d.Set("description", terraformBlueprint.Blueprint.Description)
	d.Set("category", terraformBlueprint.Blueprint.Category)
	d.Set("terraform_version", terraformBlueprint.Blueprint.Config.Terraform.Tfversion)
	d.Set("terraform_options", terraformBlueprint.Blueprint.Config.Terraform.Commandoptions)
	d.Set("tfvar_secret", terraformBlueprint.Blueprint.Config.Terraform.Tfvarsecret)

	switch terraformBlueprint.Blueprint.Config.Terraform.Configtype {
	case "tf":
		d.Set("source_type", "hcl")
		d.Set("blueprint_content", terraformBlueprint.Blueprint.Config.Terraform.Tf)
	case "json":
		d.Set("source_type", "json")
		d.Set("blueprint_content", terraformBlueprint.Blueprint.Config.Terraform.JSON)
	case "git":
		d.Set("source_type", "repository")
		d.Set("working_path", terraformBlueprint.Blueprint.Config.Terraform.Git.Path)
		d.Set("integration_id", terraformBlueprint.Blueprint.Config.Terraform.Git.IntegrationId)
		d.Set("repository_id", terraformBlueprint.Blueprint.Config.Terraform.Git.RepoId)
		d.Set("version_ref", terraformBlueprint.Blueprint.Config.Terraform.Git.Branch)
	case "spec":
		d.Set("source_type", "spec")
		// spec templates
		var specTemplates []int64
		if terraformBlueprint.Blueprint.Config.Config.Specs != nil {
			// iterate over the array of tasks
			for i := 0; i < len(terraformBlueprint.Blueprint.Config.Config.Specs); i++ {
				specTemplate := terraformBlueprint.Blueprint.Config.Config.Specs[i]
				specTemplates = append(specTemplates, int64(specTemplate.ID))
			}
		}
		d.Set("spec_templates_ids", specTemplates)
	}
	return diags
}

func resourceTerraformAppBlueprintUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	name := d.Get("name").(string)
	blueprint_type := "terraform"
	description := d.Get("description").(string)
	category := d.Get("category").(string)

	config := make(map[string]interface{})
	config["name"] = name
	config["description"] = description
	config["category"] = category
	config["type"] = "terraform"

	terraformConfig := make(map[string]interface{})
	terraformConfig["tfVersion"] = d.Get("terraform_version").(string)
	terraformConfig["commandOptions"] = d.Get("terraform_options").(string)
	terraformConfig["tfvarSecret"] = d.Get("tfvar_secret").(string)
	config["terraform"] = terraformConfig

	switch d.Get("source_type") {
	case "hcl":
		terraformConfig["configType"] = "tf"
		terraformConfig["tf"] = d.Get("blueprint_content").(string)
	case "json":
		terraformConfig["configType"] = "json"
		terraformConfig["json"] = d.Get("blueprint_content").(string)
	case "repository":
		terraformConfig["configType"] = "git"
		terraformGitConfig := make(map[string]interface{})
		terraformGitConfig["integrationId"] = d.Get("integration_id")
		terraformGitConfig["repoId"] = d.Get("repository_id")
		terraformGitConfig["branch"] = d.Get("version_ref").(string)
		terraformGitConfig["path"] = d.Get("working_path").(string)
		terraformConfig["git"] = terraformGitConfig
	case "spec":
		terraformConfig["configType"] = "spec"
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
	return resourceTerraformAppBlueprintRead(ctx, d, meta)
}

func resourceTerraformAppBlueprintDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

type TerraformAppBlueprint struct {
	Blueprint struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Type        string `json:"type"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Config      struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Terraform   struct {
				Tfversion      string `json:"tfVersion"`
				Tf             string `json:"tf"`
				Tfvarsecret    string `json:"tfvarSecret"`
				Commandoptions string `json:"commandOptions"`
				Configtype     string `json:"configType"`
				JSON           string `json:"json"`
				Git            struct {
					Path          string `json:"path"`
					RepoId        int    `json:"repoId"`
					IntegrationId int    `json:"integrationId"`
					Branch        string `json:"branch"`
				} `json:"git"`
			} `json:"terraform"`
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
