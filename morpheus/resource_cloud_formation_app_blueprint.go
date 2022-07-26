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

func resourceCloudFormationAppBlueprint() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus cloud formation app blueprint resource",
		CreateContext: resourceCloudFormationAppBlueprintCreate,
		ReadContext:   resourceCloudFormationAppBlueprintRead,
		UpdateContext: resourceCloudFormationAppBlueprintUpdate,
		DeleteContext: resourceCloudFormationAppBlueprintDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the cloud formation app blueprint",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the cloud formation app blueprint",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the cloud formation app blueprint",
				Optional:    true,
			},
			"category": {
				Type:        schema.TypeString,
				Description: "The category of the cloud formation app blueprint",
				Optional:    true,
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
			"capability_iam": {
				Type:        schema.TypeBool,
				Description: "Whether the iam capability is added to the cloud formation",
				Optional:    true,
			},
			"capability_named_iam": {
				Type:        schema.TypeBool,
				Description: "Whether the named iam capability is added to the cloud formation",
				Optional:    true,
			},
			"capability_auto_expand": {
				Type:        schema.TypeBool,
				Description: "Whether the auto expand capability is added to the cloud formation",
				Optional:    true,
			},
			"source_type": {
				Type:         schema.TypeString,
				Description:  "The source of the cloud formation app blueprint (yaml, json, repository)",
				ValidateFunc: validation.StringInSlice([]string{"yaml", "json", "repository"}, false),
				Required:     true,
			},
			"blueprint_content": {
				Type:        schema.TypeString,
				Description: "The content of the cloud formation app blueprint. Used when the hcl or json source types are specified",
				Optional:    true,
			},
			"working_path": {
				Type:        schema.TypeString,
				Description: "The path of the cloud formation chart in the git repository",
				Optional:    true,
				Default:     "./",
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
				Default:     "master",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCloudFormationAppBlueprintCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	blueprint_type := "cloudFormation"
	description := d.Get("description").(string)
	category := d.Get("category").(string)
	visibility := d.Get("visibility")

	config := make(map[string]interface{})
	config["name"] = name
	config["description"] = description
	config["category"] = category
	config["type"] = "cloudFormation"

	cloudformationConfig := make(map[string]interface{})
	config["cloudFormation"] = cloudformationConfig
	cloudformationConfig["installAgent"] = d.Get("install_agent").(bool)
	cloudformationConfig["cloudInitEnabled"] = d.Get("cloud_init_enabled").(bool)
	cloudformationConfig["IAM"] = d.Get("capability_iam").(bool)
	cloudformationConfig["CAPABILITY_NAMED_IAM"] = d.Get("capability_named_iam").(bool)
	cloudformationConfig["CAPABILITY_AUTO_EXPAND"] = d.Get("capability_auto_expand").(bool)

	switch d.Get("source_type") {
	case "json":
		cloudformationConfig["configType"] = "json"
		cloudformationConfig["json"] = d.Get("blueprint_content").(string)

	case "yaml":
		cloudformationConfig["configType"] = "yaml"
		cloudformationConfig["yaml"] = d.Get("blueprint_content").(string)

	case "repository":
		cloudformationConfig["configType"] = "git"
		cloudformationGitConfig := make(map[string]interface{})
		cloudformationGitConfig["integrationId"] = d.Get("integration_id")
		cloudformationGitConfig["repoId"] = d.Get("repository_id")
		cloudformationGitConfig["branch"] = d.Get("version_ref").(string)
		cloudformationGitConfig["path"] = d.Get("working_path").(string)
		cloudformationConfig["git"] = cloudformationGitConfig
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"blueprint": map[string]interface{}{
				"name":        name,
				"type":        blueprint_type,
				"description": description,
				"category":    category,
				"config":      config,
				"visibility":  visibility,
			},
		},
	}
	jsonRequest, _ := json.Marshal(req.Body)
	log.Printf("API JSON REQUEST: %s", string(jsonRequest))
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

	resourceCloudFormationAppBlueprintRead(ctx, d, meta)
	return diags
}

func resourceCloudFormationAppBlueprintRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	var cloudformationBlueprint CloudFormationAppBlueprint
	json.Unmarshal(resp.Body, &cloudformationBlueprint)
	d.SetId(intToString(cloudformationBlueprint.Blueprint.ID))
	d.Set("name", cloudformationBlueprint.Blueprint.Name)
	d.Set("description", cloudformationBlueprint.Blueprint.Description)
	d.Set("category", cloudformationBlueprint.Blueprint.Category)
	d.Set("visibility", cloudformationBlueprint.Blueprint.Visibility)
	d.Set("install_agent", cloudformationBlueprint.Blueprint.Config.CloudFormation.InstallAgent)
	d.Set("cloud_init_enabled", cloudformationBlueprint.Blueprint.Config.CloudFormation.CloudInitEnabled)
	d.Set("capability_iam", cloudformationBlueprint.Blueprint.Config.CloudFormation.IAM)
	d.Set("capability_named_iam", cloudformationBlueprint.Blueprint.Config.CloudFormation.IAMNamed)
	d.Set("capability_auto_expand", cloudformationBlueprint.Blueprint.Config.CloudFormation.AutoExpand)

	switch cloudformationBlueprint.Blueprint.Config.CloudFormation.Configtype {
	case "json":
		d.Set("source_type", "json")
		d.Set("blueprint_content", cloudformationBlueprint.Blueprint.Config.CloudFormation.JSON)
	case "yaml":
		d.Set("source_type", "yaml")
		d.Set("blueprint_content", cloudformationBlueprint.Blueprint.Config.CloudFormation.YAML)
	case "git":
		d.Set("source_type", "repository")
		d.Set("working_path", cloudformationBlueprint.Blueprint.Config.CloudFormation.Git.Path)
		d.Set("integration_id", cloudformationBlueprint.Blueprint.Config.CloudFormation.Git.IntegrationId)
		d.Set("repository_id", cloudformationBlueprint.Blueprint.Config.CloudFormation.Git.RepoId)
		d.Set("version_ref", cloudformationBlueprint.Blueprint.Config.CloudFormation.Git.Branch)
	}
	d.Set("visibility", cloudformationBlueprint.Blueprint.Visibility)
	return diags
}

func resourceCloudFormationAppBlueprintUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	name := d.Get("name").(string)
	blueprint_type := "cloudFormation"
	description := d.Get("description").(string)
	category := d.Get("category").(string)
	visibility := d.Get("visibility")

	config := make(map[string]interface{})
	config["name"] = name
	config["description"] = description
	config["category"] = category
	config["type"] = "cloudFormation"

	cloudformationConfig := make(map[string]interface{})
	config["cloudFormation"] = cloudformationConfig
	cloudformationConfig["installAgent"] = d.Get("install_agent").(bool)
	cloudformationConfig["cloudInitEnabled"] = d.Get("cloud_init_enabled").(bool)
	cloudformationConfig["IAM"] = d.Get("capability_iam").(bool)
	cloudformationConfig["CAPABILITY_NAMED_IAM"] = d.Get("capability_named_iam").(bool)
	cloudformationConfig["CAPABILITY_AUTO_EXPAND"] = d.Get("capability_auto_expand").(bool)

	switch d.Get("source_type") {
	case "json":
		cloudformationConfig["configType"] = "json"
		cloudformationConfig["json"] = d.Get("blueprint_content").(string)

	case "yaml":
		cloudformationConfig["configType"] = "yaml"
		cloudformationConfig["yaml"] = d.Get("blueprint_content").(string)

	case "repository":
		cloudformationConfig["configType"] = "git"
		cloudformationGitConfig := make(map[string]interface{})
		cloudformationGitConfig["integrationId"] = d.Get("integration_id")
		cloudformationGitConfig["repoId"] = d.Get("repository_id")
		cloudformationGitConfig["branch"] = d.Get("version_ref").(string)
		cloudformationGitConfig["path"] = d.Get("working_path").(string)
		cloudformationConfig["git"] = cloudformationGitConfig
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"blueprint": map[string]interface{}{
				"name":        name,
				"type":        blueprint_type,
				"description": description,
				"category":    category,
				"config":      config,
				"visibility":  visibility,
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
	return resourceCloudFormationAppBlueprintRead(ctx, d, meta)
}

func resourceCloudFormationAppBlueprintDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

type CloudFormationAppBlueprint struct {
	Blueprint struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Type        string `json:"type"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Config      struct {
			Name           string `json:"name"`
			Description    string `json:"description"`
			CloudFormation struct {
				Configtype       string `json:"configType"`
				CloudInitEnabled bool   `json:"cloudInitEnabled"`
				InstallAgent     bool   `json:"installAgent"`
				JSON             string `json:"json"`
				YAML             string `json:"yaml"`
				IAM              bool   `json:"IAM"`
				IAMNamed         bool   `json:"CAPABILITY_NAMED_IAM"`
				AutoExpand       bool   `json:"CAPABILITY_AUTO_EXPAND"`
				Git              struct {
					Path          string `json:"path"`
					RepoId        int    `json:"repoId"`
					IntegrationId int    `json:"integrationId"`
					Branch        string `json:"branch"`
				} `json:"git"`
			} `json:"cloudformation"`
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
