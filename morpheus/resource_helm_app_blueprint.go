package morpheus

import (
	"context"
	"encoding/json"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceHelmAppBlueprint() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus helm app blueprint resource",
		CreateContext: resourceHelmAppBlueprintCreate,
		ReadContext:   resourceHelmAppBlueprintRead,
		UpdateContext: resourceHelmAppBlueprintUpdate,
		DeleteContext: resourceHelmAppBlueprintDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the helm app blueprint",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the helm app blueprint",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the helm app blueprint",
				Optional:    true,
			},
			"category": {
				Type:        schema.TypeString,
				Description: "The category of the helm app blueprint",
				Optional:    true,
			},
			"working_path": {
				Type:        schema.TypeString,
				Description: "The path of the helm chart in the git repository",
				Optional:    true,
				Default:     "./",
			},
			"integration_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the git integration",
				Required:    true,
			},
			"repository_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the git repository",
				Required:    true,
			},
			"version_ref": {
				Type:        schema.TypeString,
				Description: "The git reference of the repository to pull (main, master, etc.)",
				Optional:    true,
				Default:     "master",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceHelmAppBlueprintCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	blueprint_type := "helm"
	description := d.Get("description").(string)
	category := d.Get("category").(string)

	config := make(map[string]interface{})
	config["name"] = name
	config["description"] = description
	config["category"] = category
	config["type"] = "helm"

	helmConfig := make(map[string]interface{})
	config["helm"] = helmConfig

	helmConfig["configType"] = "git"
	helmGitConfig := make(map[string]interface{})
	helmGitConfig["integrationId"] = d.Get("integration_id")
	helmGitConfig["repoId"] = d.Get("repository_id")
	helmGitConfig["branch"] = d.Get("version_ref").(string)
	helmGitConfig["path"] = d.Get("working_path").(string)
	helmConfig["git"] = helmGitConfig

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

	resourceHelmAppBlueprintRead(ctx, d, meta)
	return diags
}

func resourceHelmAppBlueprintRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	var helmBlueprint HelmAppBlueprint
	json.Unmarshal(resp.Body, &helmBlueprint)
	d.SetId(intToString(helmBlueprint.Blueprint.ID))
	d.Set("name", helmBlueprint.Blueprint.Name)
	d.Set("description", helmBlueprint.Blueprint.Description)
	d.Set("category", helmBlueprint.Blueprint.Category)
	d.Set("working_path", helmBlueprint.Blueprint.Config.Helm.Git.Path)
	d.Set("integration_id", helmBlueprint.Blueprint.Config.Helm.Git.IntegrationId)
	d.Set("repository_id", helmBlueprint.Blueprint.Config.Helm.Git.RepoId)
	d.Set("version_ref", helmBlueprint.Blueprint.Config.Helm.Git.Branch)

	return diags
}

func resourceHelmAppBlueprintUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	name := d.Get("name").(string)
	blueprint_type := "helm"
	description := d.Get("description").(string)
	category := d.Get("category").(string)

	config := make(map[string]interface{})
	config["name"] = name
	config["description"] = description
	config["category"] = category
	config["type"] = "helm"

	helmConfig := make(map[string]interface{})
	config["helm"] = helmConfig

	helmConfig["configType"] = "git"
	helmGitConfig := make(map[string]interface{})
	helmGitConfig["integrationId"] = d.Get("integration_id")
	helmGitConfig["repoId"] = d.Get("repository_id")
	helmGitConfig["branch"] = d.Get("version_ref").(string)
	helmGitConfig["path"] = d.Get("working_path").(string)
	helmConfig["git"] = helmGitConfig

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
	return resourceHelmAppBlueprintRead(ctx, d, meta)
}

func resourceHelmAppBlueprintDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

type HelmAppBlueprint struct {
	Blueprint struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Type        string `json:"type"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Config      struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Helm        struct {
				Configtype string `json:"configType"`
				Git        struct {
					Path          string `json:"path"`
					RepoId        int    `json:"repoId"`
					IntegrationId int    `json:"integrationId"`
					Branch        string `json:"branch"`
				} `json:"git"`
			} `json:"helm"`
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
