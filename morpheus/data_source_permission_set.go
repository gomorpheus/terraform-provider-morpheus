package morpheus

import (
	"context"
	"encoding/json"
	"log"
	"sort"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceMorpheusPermissionSet() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus permission set data source.",
		ReadContext: dataSourceMorpheusPermissionSetRead,
		Schema: map[string]*schema.Schema{
			"json": {
				Type:        schema.TypeString,
				Description: "JSON permission set rendered based on the arguments defined",
				Computed:    true,
			},
			"default_group_permission": {
				Type:         schema.TypeString,
				Description:  "The default role permission for groups (none, read, full)",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"none", "read", "full"}, true),
			},
			"default_instance_type_permission": {
				Type:         schema.TypeString,
				Description:  "The default role permission for instance types (none, full)",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"none", "full"}, true),
			},
			"default_blueprint_permission": {
				Type:         schema.TypeString,
				Description:  "The default role permission for blueprints (none, full)",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"none", "full"}, true),
			},
			"default_report_type_permission": {
				Type:         schema.TypeString,
				Description:  "The default role permission for report types (none, full)",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"none", "full"}, true),
			},
			"default_persona_permission": {
				Type:         schema.TypeString,
				Description:  "The default role permission for personas (none, full)",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"none", "full"}, true),
			},
			"default_persona": {
				Type:         schema.TypeString,
				Description:  "The default role persona (standard, serviceCatalog, vdi)",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"standard", "serviceCatalog", "vdi"}, true),
			},
			"default_catalog_item_type_permission": {
				Type:         schema.TypeString,
				Description:  "The default role permission for catalog item types (none, full)",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"none", "full"}, true),
			},
			"default_vdi_pool_permission": {
				Type:         schema.TypeString,
				Description:  "The default role permission for vdi pools (none, full)",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"none", "full"}, true),
			},
			"default_workflow_permission": {
				Type:         schema.TypeString,
				Description:  "The default role permission for workflows (none, full)",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"none", "full"}, true),
			},
			"default_task_permission": {
				Type:         schema.TypeString,
				Description:  "The default role permission for tasks (none, full)",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"none", "full"}, true),
			},
			"feature_permission": {
				Type:        schema.TypeList,
				Description: "The feature permissions associated with the user role",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:        schema.TypeString,
							Description: "The code of the feature permission",
							Optional:    true,
						},
						"access": {
							Type:         schema.TypeString,
							Description:  "The level of access granted to the feature permission (full, full_decrypted, group, listfiles, managerules, no, none, provision, read, rolemappings, user, view, yes)",
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"full", "full_decrypted", "group", "listfiles", "managerules", "no", "none", "provision", "read", "rolemappings", "user", "view", "yes"}, true),
						},
					},
				},
			},
			"group_permission": {
				Type:        schema.TypeList,
				Description: "The group permissions associated with the user role",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Description: "The id of the group",
							Optional:    true,
						},
						"access": {
							Type:        schema.TypeString,
							Description: "The level of access granted to the group (default, full, none)",
							Optional:    true,
						},
					},
				},
			},
			"instance_type_permission": {
				Type:        schema.TypeList,
				Description: "The instance type permissions associated with the user role",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:        schema.TypeString,
							Description: "The code of the instance type",
							Optional:    true,
						},
						"access": {
							Type:        schema.TypeString,
							Description: "The level of access granted to the workflow (default, full, none)",
							Optional:    true,
						},
					},
				},
			},
			"blueprint_permission": {
				Type:        schema.TypeList,
				Description: "The blueprint permissions associated with the user role",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Description: "The id of the blueprint",
							Optional:    true,
						},
						"access": {
							Type:        schema.TypeString,
							Description: "The level of access granted to the blueprint (default, full, none)",
							Optional:    true,
						},
					},
				},
			},
			"report_type_permission": {
				Type:        schema.TypeList,
				Description: "The report type permissions associated with the user role",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "The name of report",
							Optional:    true,
						},
						"access": {
							Type:        schema.TypeString,
							Description: "The level of access granted to the report type (default, full, none)",
							Optional:    true,
						},
					},
				},
			},
			"persona_permission": {
				Type:        schema.TypeList,
				Description: "The persona permissions associated with the user role",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:        schema.TypeString,
							Description: "The name of the environment variable",
							Optional:    true,
						},
						"access": {
							Type:        schema.TypeString,
							Description: "The level of access granted to the persona (default, full, none)",
							Optional:    true,
						},
					},
				},
			},
			"catalog_item_type_permission": {
				Type:        schema.TypeList,
				Description: "The catalog item permissions associated with the user role",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "The name of the environment variable",
							Optional:    true,
						},
						"access": {
							Type:        schema.TypeString,
							Description: "The level of access granted to the catalog item (default, full, none)",
							Optional:    true,
						},
					},
				},
			},
			"vdi_permission": {
				Type:        schema.TypeList,
				Description: "The vdi pool permissions associated with the user role",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "The name of the environment variable",
							Optional:    true,
						},
						"access": {
							Type:        schema.TypeString,
							Description: "The level of access granted to the vdi pool (default, full, none)",
							Optional:    true,
						},
					},
				},
			},
			"workflow_permission": {
				Type:        schema.TypeList,
				Description: "The workflow permissions associated with the user role",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Description: "The id of the workflow",
							Optional:    true,
						},
						"access": {
							Type:        schema.TypeString,
							Description: "The level of access granted to the workflow (default, full, none)",
							Optional:    true,
						},
					},
				},
			},
			"task_permission": {
				Type:        schema.TypeList,
				Description: "The task permissions associated with the user role",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Description: "The id of the task",
							Optional:    true,
						},
						"access": {
							Type:        schema.TypeString,
							Description: "The level of access granted to the task (default, full, none)",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceMorpheusPermissionSetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var permissionData PermissionSet

	permissionData.DefaultGroupPermission = d.Get("default_group_permission").(string)
	//	demo.DefaultCloudPermission = d.Get("default_cloud_permission").(string)
	permissionData.DefaultInstanceTypePermission = d.Get("default_instance_type_permission").(string)
	permissionData.DefaultBlueprintPermission = d.Get("default_blueprint_permission").(string)
	permissionData.DefaultReportTypePermission = d.Get("default_report_type_permission").(string)
	permissionData.DefaultPersonaPermission = d.Get("default_persona_permission").(string)
	permissionData.DefaultPersona = d.Get("default_persona").(string)
	permissionData.DefaultCatalogItemTypePermission = d.Get("default_catalog_item_type_permission").(string)
	permissionData.DefaultVdiPoolPermission = d.Get("default_vdi_pool_permission").(string)
	permissionData.DefaultWorkflowPermission = d.Get("default_workflow_permission").(string)
	permissionData.DefaultTaskPermission = d.Get("default_task_permission").(string)

	var features []featurePermission
	if d.Get("feature_permission") != nil {
		taskList := d.Get("feature_permission").([]interface{})
		// iterate over the array of tasks
		for i := 0; i < len(taskList); i++ {
			var row featurePermission
			taskconfig := taskList[i].(map[string]interface{})
			row.Code = taskconfig["code"].(string)
			row.Access = taskconfig["access"].(string)
			features = append(features, row)
		}
	}
	sort.Slice(features, func(i, j int) bool { return features[i].Code < features[j].Code })

	permissionData.FeaturePermissions = features

	var groups []groupPermission
	if d.Get("group_permission") != nil {
		taskList := d.Get("group_permission").([]interface{})
		// iterate over the array of groups
		for i := 0; i < len(taskList); i++ {
			var row groupPermission
			taskconfig := taskList[i].(map[string]interface{})
			row.Id = taskconfig["id"].(int)
			row.Access = taskconfig["access"].(string)
			groups = append(groups, row)
		}
	}

	permissionData.GroupPermissions = groups

	var instanceTypes []instanceTypePermission
	if d.Get("instance_type_permission") != nil {
		taskList := d.Get("instance_type_permission").([]interface{})
		// iterate over the array of groups
		for i := 0; i < len(taskList); i++ {
			var row instanceTypePermission
			taskconfig := taskList[i].(map[string]interface{})
			row.Code = taskconfig["code"].(string)
			row.Access = taskconfig["access"].(string)
			instanceTypes = append(instanceTypes, row)
		}
	}

	sort.Slice(instanceTypes, func(i, j int) bool { return instanceTypes[i].Code < instanceTypes[j].Code })

	permissionData.InstanceTypePermissions = instanceTypes

	// Personas
	var personas []personaPermission
	if d.Get("instance_type_permission") != nil {
		taskList := d.Get("persona_permission").([]interface{})
		// iterate over the array of groups
		for i := 0; i < len(taskList); i++ {
			var row personaPermission
			taskconfig := taskList[i].(map[string]interface{})
			row.Code = taskconfig["code"].(string)
			row.Access = taskconfig["access"].(string)
			personas = append(personas, row)
		}
	}

	sort.Slice(personas, func(i, j int) bool { return personas[i].Code < personas[j].Code })

	permissionData.PersonaPermissions = personas

	// Tasks
	var tasks []taskPermission
	if d.Get("task_permission") != nil {
		taskList := d.Get("task_permission").([]interface{})
		// iterate over the array of tasks
		for i := 0; i < len(taskList); i++ {
			var row taskPermission
			taskconfig := taskList[i].(map[string]interface{})
			row.Id = taskconfig["id"].(int)
			row.Access = taskconfig["access"].(string)
			tasks = append(tasks, row)
		}
	}

	sort.Slice(tasks, func(i, j int) bool { return tasks[i].Id < tasks[j].Id })

	permissionData.TaskPermissions = tasks

	var workflows []workflowPermission
	if d.Get("workflow_permission") != nil {
		taskList := d.Get("workflow_permission").([]interface{})
		// iterate over the array of workflows
		for i := 0; i < len(taskList); i++ {
			var row workflowPermission
			taskconfig := taskList[i].(map[string]interface{})
			row.Id = taskconfig["id"].(int)
			row.Access = taskconfig["access"].(string)
			workflows = append(workflows, row)
		}
	}

	sort.Slice(workflows, func(i, j int) bool { return workflows[i].Id < workflows[j].Id })

	permissionData.WorkflowPermissions = workflows
	jsonDoc, err := json.MarshalIndent(permissionData, "", "  ")
	log.Printf("API RESPONSE: %s", jsonDoc)

	if err != nil {
		log.Println("error")
		// should never happen if the above code is correct
		//		return diags.AppendErrorf(diags, "writing IAM Policy Document: formatting JSON: %s", err)
	}
	jsonString := string(jsonDoc)

	d.Set("json", jsonString)
	d.SetId(strconv.Itoa(1))
	return diags
}

type PermissionSet struct {
	DefaultGroupPermission           string                      `json:"default_group_permission,omitempty"`
	DefaultCloudPermission           string                      `json:"default_cloud_permission,omitempty"`
	DefaultInstanceTypePermission    string                      `json:"default_instance_type_permission,omitempty"`
	DefaultBlueprintPermission       string                      `json:"default_blueprint_permission,omitempty"`
	DefaultReportTypePermission      string                      `json:"default_report_type_permission,omitempty"`
	DefaultPersonaPermission         string                      `json:"default_persona_permission,omitempty"`
	DefaultCatalogItemTypePermission string                      `json:"default_catalog_item_type_permission,omitempty"`
	DefaultVdiPoolPermission         string                      `json:"default_vdi_pool_permission,omitempty"`
	DefaultWorkflowPermission        string                      `json:"default_workflow_permission,omitempty"`
	DefaultTaskPermission            string                      `json:"default_task_permission,omitempty"`
	DefaultPersona                   string                      `json:"default_persona,omitempty"`
	FeaturePermissions               []featurePermission         `json:"feature_permissions,omitempty"`
	GroupPermissions                 []groupPermission           `json:"group_permissions,omitempty"`
	InstanceTypePermissions          []instanceTypePermission    `json:"instance_type_permissions,omitempty"`
	BlueprintPermissions             []blueprintPermission       `json:"blueprint_permissions,omitempty"`
	ReportTypePermissions            []reportTypePermission      `json:"report_type_permissions,omitempty"`
	PersonaPermissions               []personaPermission         `json:"persona_permissions,omitempty"`
	CatalogItemTypePermissions       []catalogItemTypePermission `json:"catalog_item_type_permissions,omitempty"`
	VdiPoolPermissions               []vdiPoolPermission         `json:"vdi_pool_permissions,omitempty"`
	TaskPermissions                  []taskPermission            `json:"task_permissions,omitempty"`
	WorkflowPermissions              []workflowPermission        `json:"workflow_permissions,omitempty"`
}

type featurePermission struct {
	Code   string `json:"code"`
	Access string `json:"access"`
}

type groupPermission struct {
	Id     int    `json:"id"`
	Access string `json:"access"`
}

type instanceTypePermission struct {
	Code   string `json:"code"`
	Access string `json:"access"`
}

type blueprintPermission struct {
	Code   string `json:"code"`
	Access string `json:"access"`
}

type reportTypePermission struct {
	Code   string `json:"code"`
	Access string `json:"access"`
}

type personaPermission struct {
	Code   string `json:"code"`
	Access string `json:"access"`
}

type catalogItemTypePermission struct {
	Code   string `json:"code"`
	Access string `json:"access"`
}

type vdiPoolPermission struct {
	Id     int    `json:"id"`
	Access string `json:"access"`
}

type taskPermission struct {
	Id     int    `json:"id"`
	Access string `json:"access"`
}

type workflowPermission struct {
	Id     int    `json:"id"`
	Access string `json:"access"`
}
