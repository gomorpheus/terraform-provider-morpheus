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

func resourceArmAppBlueprint() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus arm app blueprint resource",
		CreateContext: resourceArmAppBlueprintCreate,
		ReadContext:   resourceArmAppBlueprintRead,
		UpdateContext: resourceArmAppBlueprintUpdate,
		DeleteContext: resourceArmAppBlueprintDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the arm app blueprint",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the arm app blueprint",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the arm app blueprint",
				Optional:    true,
				Computed:    true,
			},
			"category": {
				Type:        schema.TypeString,
				Description: "The category of the arm app blueprint",
				Optional:    true,
				Computed:    true,
			},
			"install_agent": {
				Type:        schema.TypeBool,
				Description: "Whether to install the Morpheus agent",
				Optional:    true,
			},
			"cloud_init_enabled": {
				Type:        schema.TypeBool,
				Description: "Whether cloud init is enabled",
				Optional:    true,
			},
			"os_type": {
				Type:         schema.TypeString,
				Description:  "The workload operating system type (linux, windows)",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"linux", "windows"}, false),
			},
			"source_type": {
				Type:         schema.TypeString,
				Description:  "The source of the arm app blueprint (json, repository)",
				ValidateFunc: validation.StringInSlice([]string{"json", "repository"}, false),
				Required:     true,
			},
			"blueprint_content": {
				Type:        schema.TypeString,
				Description: "The content of the arm app blueprint. Used when the json source type is specified",
				Optional:    true,
			},
			"working_path": {
				Type:        schema.TypeString,
				Description: "The path of the arm app blueprint in the git repository",
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
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceArmAppBlueprintCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	blueprint_type := "arm"
	description := d.Get("description").(string)
	category := d.Get("category").(string)

	config := make(map[string]interface{})
	config["name"] = name
	config["description"] = description
	config["category"] = category
	config["type"] = "arm"

	armConfig := make(map[string]interface{})
	config["arm"] = armConfig
	armConfig["osType"] = d.Get("os_type").(string)
	armConfig["installAgent"] = d.Get("install_agent").(bool)
	armConfig["cloudInitEnabled"] = d.Get("cloud_init_enabled").(bool)

	switch d.Get("source_type") {
	case "json":
		armConfig["configType"] = "json"
		armConfig["json"] = d.Get("blueprint_content").(string)

	case "repository":
		armConfig["configType"] = "git"
		armGitConfig := make(map[string]interface{})
		armGitConfig["integrationId"] = d.Get("integration_id")
		armGitConfig["repoId"] = d.Get("repository_id")
		armGitConfig["branch"] = d.Get("version_ref").(string)
		armGitConfig["path"] = d.Get("working_path").(string)
		armConfig["git"] = armGitConfig
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

	resourceArmAppBlueprintRead(ctx, d, meta)
	return diags
}

func resourceArmAppBlueprintRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
			return diag.FromErr(err)
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)
			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)

	// store resource data
	var armBlueprint ArmAppBlueprint
	json.Unmarshal(resp.Body, &armBlueprint)
	d.SetId(intToString(armBlueprint.Blueprint.ID))
	d.Set("name", armBlueprint.Blueprint.Name)
	d.Set("description", armBlueprint.Blueprint.Description)
	d.Set("category", armBlueprint.Blueprint.Category)
	d.Set("install_agent", armBlueprint.Blueprint.Config.Arm.InstallAgent)
	d.Set("cloud_init_enabled", armBlueprint.Blueprint.Config.Arm.CloudInitEnabled)
	d.Set("os_type", armBlueprint.Blueprint.Config.Arm.OsType)
	switch armBlueprint.Blueprint.Config.Arm.Configtype {
	case "json":
		d.Set("source_type", "json")
		d.Set("blueprint_content", armBlueprint.Blueprint.Config.Arm.JSON)
	case "git":
		d.Set("source_type", "repository")
		d.Set("working_path", armBlueprint.Blueprint.Config.Arm.Git.Path)
		d.Set("integration_id", armBlueprint.Blueprint.Config.Arm.Git.IntegrationId)
		d.Set("repository_id", armBlueprint.Blueprint.Config.Arm.Git.RepoId)
		d.Set("version_ref", armBlueprint.Blueprint.Config.Arm.Git.Branch)
	}
	return diags
}

func resourceArmAppBlueprintUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	name := d.Get("name").(string)
	blueprint_type := "arm"
	description := d.Get("description").(string)
	category := d.Get("category").(string)

	config := make(map[string]interface{})
	config["name"] = name
	config["description"] = description
	config["category"] = category
	config["type"] = "arm"

	armConfig := make(map[string]interface{})
	config["arm"] = armConfig
	armConfig["osType"] = d.Get("os_type").(string)
	armConfig["installAgent"] = d.Get("install_agent").(bool)
	armConfig["cloudInitEnabled"] = d.Get("cloud_init_enabled").(bool)

	switch d.Get("source_type") {
	case "json":
		armConfig["configType"] = "json"
		armConfig["json"] = d.Get("blueprint_content").(string)

	case "repository":
		armConfig["configType"] = "git"
		armGitConfig := make(map[string]interface{})
		armGitConfig["integrationId"] = d.Get("integration_id")
		armGitConfig["repoId"] = d.Get("repository_id")
		armGitConfig["branch"] = d.Get("version_ref").(string)
		armGitConfig["path"] = d.Get("working_path").(string)
		armConfig["git"] = armGitConfig
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
	return resourceArmAppBlueprintRead(ctx, d, meta)
}

func resourceArmAppBlueprintDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

type ArmAppBlueprint struct {
	Blueprint struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Type        string `json:"type"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Config      struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Arm         struct {
				Configtype       string `json:"configType"`
				OsType           string `json:"osType"`
				CloudInitEnabled bool   `json:"cloudInitEnabled"`
				InstallAgent     bool   `json:"installAgent"`
				JSON             string `json:"json"`
				Git              struct {
					Path          string `json:"path"`
					RepoId        int    `json:"repoId"`
					IntegrationId int    `json:"integrationId"`
					Branch        string `json:"branch"`
				} `json:"git"`
			} `json:"arm"`
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
