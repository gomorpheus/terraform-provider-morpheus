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
						"id": {
							Type:        schema.TypeInt,
							Description: "The id of the instance type",
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
						"code": {
							Type:        schema.TypeString,
							Description: "The report type code",
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
							Description: "The code of the persona (standard, vdi, serviceCatalog)",
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
				Description: "The catalog item type permissions associated with the user role",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Description: "The id of the catalog item type",
							Optional:    true,
						},
						"access": {
							Type:        schema.TypeString,
							Description: "The level of access granted to the catalog item type (default, full, none)",
							Optional:    true,
						},
					},
				},
			},
			"vdi_pool_permission": {
				Type:        schema.TypeList,
				Description: "The vdi pool permissions associated with the user role",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Description: "The id of the vdi pool",
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
			"override_permission_sets": {
				Type:        schema.TypeList,
				Description: "List of permission sets that are merged together into the exported json. In merging, the last permission applied in the list order is used. Non-overriding permissions will be added to the exported json.",
				Optional:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsJSON,
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
	permissionData.DefaultInstanceTypePermission = d.Get("default_instance_type_permission").(string)
	permissionData.DefaultBlueprintPermission = d.Get("default_blueprint_permission").(string)
	permissionData.DefaultReportTypePermission = d.Get("default_report_type_permission").(string)
	permissionData.DefaultPersonaPermission = d.Get("default_persona_permission").(string)
	permissionData.DefaultPersona = d.Get("default_persona").(string)
	permissionData.DefaultCatalogItemTypePermission = d.Get("default_catalog_item_type_permission").(string)
	permissionData.DefaultVdiPoolPermission = d.Get("default_vdi_pool_permission").(string)
	permissionData.DefaultWorkflowPermission = d.Get("default_workflow_permission").(string)
	permissionData.DefaultTaskPermission = d.Get("default_task_permission").(string)

	// Feature Permissions
	var features []featurePermission
	if d.Get("feature_permission") != nil {
		featureList := d.Get("feature_permission").([]interface{})
		// iterate over the array of features
		for i := 0; i < len(featureList); i++ {
			var row featurePermission
			featureConfig := featureList[i].(map[string]interface{})
			row.Code = featureConfig["code"].(string)
			row.Access = featureConfig["access"].(string)
			features = append(features, row)
		}
	}
	sort.Slice(features, func(i, j int) bool { return features[i].Code < features[j].Code })
	permissionData.FeaturePermissions = features

	// Group Permissions
	var groups []groupPermission
	if d.Get("group_permission") != nil {
		groupList := d.Get("group_permission").([]interface{})
		// iterate over the array of groups
		for i := 0; i < len(groupList); i++ {
			var row groupPermission
			groupConfig := groupList[i].(map[string]interface{})
			row.Id = groupConfig["id"].(int)
			row.Access = groupConfig["access"].(string)
			groups = append(groups, row)
		}
	}
	sort.Slice(groups, func(i, j int) bool { return groups[i].Id < groups[j].Id })
	permissionData.GroupPermissions = groups

	// Instance Type Permissions
	var instanceTypes []instanceTypePermission
	if d.Get("instance_type_permission") != nil {
		instanceTypeList := d.Get("instance_type_permission").([]interface{})
		// iterate over the array of instance types
		for i := 0; i < len(instanceTypeList); i++ {
			var row instanceTypePermission
			instanceTypeConfig := instanceTypeList[i].(map[string]interface{})
			row.Id = instanceTypeConfig["id"].(int)
			row.Access = instanceTypeConfig["access"].(string)
			instanceTypes = append(instanceTypes, row)
		}
	}
	sort.Slice(instanceTypes, func(i, j int) bool { return instanceTypes[i].Id < instanceTypes[j].Id })
	permissionData.InstanceTypePermissions = instanceTypes

	// Blueprints Permissions
	var blueprints []blueprintPermission
	if d.Get("blueprint_permission") != nil {
		blueprintList := d.Get("blueprint_permission").([]interface{})
		// iterate over the array of blueprints
		for i := 0; i < len(blueprintList); i++ {
			var row blueprintPermission
			blueprintConfig := blueprintList[i].(map[string]interface{})
			row.Id = blueprintConfig["id"].(int)
			row.Access = blueprintConfig["access"].(string)
			blueprints = append(blueprints, row)
		}
	}
	sort.Slice(blueprints, func(i, j int) bool { return blueprints[i].Id < blueprints[j].Id })
	permissionData.BlueprintPermissions = blueprints

	// Report Types Permissions
	var reportTypes []reportTypePermission
	if d.Get("report_type_permission") != nil {
		reportTypeList := d.Get("report_type_permission").([]interface{})
		// iterate over the array of report types
		for i := 0; i < len(reportTypeList); i++ {
			var row reportTypePermission
			reportTypeConfig := reportTypeList[i].(map[string]interface{})
			row.Code = reportTypeConfig["code"].(string)
			row.Access = reportTypeConfig["access"].(string)
			reportTypes = append(reportTypes, row)
		}
	}
	sort.Slice(reportTypes, func(i, j int) bool { return reportTypes[i].Code < reportTypes[j].Code })
	permissionData.ReportTypePermissions = reportTypes

	// Persona Permissions
	var personas []personaPermission
	if d.Get("persona_permission") != nil {
		personaList := d.Get("persona_permission").([]interface{})
		// iterate over the array of personas
		for i := 0; i < len(personaList); i++ {
			var row personaPermission
			personaConfig := personaList[i].(map[string]interface{})
			row.Code = personaConfig["code"].(string)
			row.Access = personaConfig["access"].(string)
			personas = append(personas, row)
		}
	}
	sort.Slice(personas, func(i, j int) bool { return personas[i].Code < personas[j].Code })
	permissionData.PersonaPermissions = personas

	// Catalog Item Types Permissions
	var catalogItemTypes []catalogItemTypePermission
	if d.Get("catalog_item_type_permission") != nil {
		catalogItemTypeList := d.Get("catalog_item_type_permission").([]interface{})
		// iterate over the array of catalog item types
		for i := 0; i < len(catalogItemTypeList); i++ {
			var row catalogItemTypePermission
			catalogItemTypeConfig := catalogItemTypeList[i].(map[string]interface{})
			row.Id = catalogItemTypeConfig["id"].(int)
			row.Access = catalogItemTypeConfig["access"].(string)
			catalogItemTypes = append(catalogItemTypes, row)
		}
	}
	sort.Slice(catalogItemTypes, func(i, j int) bool { return catalogItemTypes[i].Id < catalogItemTypes[j].Id })
	permissionData.CatalogItemTypePermissions = catalogItemTypes

	// VDI Pool Permissions
	var vdiPools []vdiPoolPermission
	if d.Get("vdi_pool_permission") != nil {
		vdiPoolList := d.Get("vdi_pool_permission").([]interface{})
		// iterate over the array of vdi pools
		for i := 0; i < len(vdiPoolList); i++ {
			var row vdiPoolPermission
			vdiPoolConfig := vdiPoolList[i].(map[string]interface{})
			row.Id = vdiPoolConfig["id"].(int)
			row.Access = vdiPoolConfig["access"].(string)
			vdiPools = append(vdiPools, row)
		}
	}
	sort.Slice(vdiPools, func(i, j int) bool { return vdiPools[i].Id < vdiPools[j].Id })
	permissionData.VdiPoolPermissions = vdiPools

	// Task Permissions
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

	// Workflow Permissions
	var workflows []workflowPermission
	if d.Get("workflow_permission") != nil {
		workflowList := d.Get("workflow_permission").([]interface{})
		// iterate over the array of workflows
		for i := 0; i < len(workflowList); i++ {
			var row workflowPermission
			workflowConfig := workflowList[i].(map[string]interface{})
			row.Id = workflowConfig["id"].(int)
			row.Access = workflowConfig["access"].(string)
			workflows = append(workflows, row)
		}
	}
	sort.Slice(workflows, func(i, j int) bool { return workflows[i].Id < workflows[j].Id })
	permissionData.WorkflowPermissions = workflows

	// Merge override permission sets with the base permission set in the order specified
	if v, ok := d.GetOk("override_permission_sets"); ok && len(v.([]interface{})) > 0 {
		for _, overrideJSON := range v.([]interface{}) {
			if overrideJSON == nil {
				continue
			}
			overridePayload := &PermissionSet{}
			json.Unmarshal([]byte(overrideJSON.(string)), overridePayload)
			if overridePayload.DefaultGroupPermission != "" {
				permissionData.DefaultGroupPermission = overridePayload.DefaultGroupPermission
			}
			if overridePayload.DefaultInstanceTypePermission != "" {
				permissionData.DefaultInstanceTypePermission = overridePayload.DefaultInstanceTypePermission
			}
			if overridePayload.DefaultBlueprintPermission != "" {
				permissionData.DefaultBlueprintPermission = overridePayload.DefaultBlueprintPermission
			}
			if overridePayload.DefaultReportTypePermission != "" {
				permissionData.DefaultReportTypePermission = overridePayload.DefaultReportTypePermission
			}
			if overridePayload.DefaultPersonaPermission != "" {
				permissionData.DefaultPersonaPermission = overridePayload.DefaultPersonaPermission
			}
			if overridePayload.DefaultPersona != "" {
				permissionData.DefaultPersona = overridePayload.DefaultPersona
			}
			if overridePayload.DefaultCatalogItemTypePermission != "" {
				permissionData.DefaultCatalogItemTypePermission = overridePayload.DefaultCatalogItemTypePermission
			}
			if overridePayload.DefaultVdiPoolPermission != "" {
				permissionData.DefaultVdiPoolPermission = overridePayload.DefaultVdiPoolPermission
			}
			if overridePayload.DefaultWorkflowPermission != "" {
				permissionData.DefaultWorkflowPermission = overridePayload.DefaultWorkflowPermission
			}
			if overridePayload.DefaultTaskPermission != "" {
				permissionData.DefaultTaskPermission = overridePayload.DefaultTaskPermission
			}

			// Feature Permissions
			if len(overridePayload.FeaturePermissions) > 0 {
				for indx, perm := range permissionData.FeaturePermissions {
					for _, overperm := range overridePayload.FeaturePermissions {
						if perm.Code == overperm.Code {
							permissionData.FeaturePermissions[indx].Access = overperm.Access
						}
					}
				}
			}
			permissionData.FeaturePermissions = append(permissionData.FeaturePermissions, overridePayload.FeaturePermissions...)
			permissionData.FeaturePermissions = removeDuplicate(permissionData.FeaturePermissions)
			sort.Slice(permissionData.FeaturePermissions, func(i, j int) bool {
				return permissionData.FeaturePermissions[i].Code < permissionData.FeaturePermissions[j].Code
			})

			// Group Permissions
			if len(overridePayload.GroupPermissions) > 0 {
				for indx, perm := range permissionData.GroupPermissions {
					for _, overperm := range overridePayload.GroupPermissions {
						if perm.Id == overperm.Id {
							permissionData.GroupPermissions[indx].Access = overperm.Access
						}
					}
				}
			}
			permissionData.GroupPermissions = append(permissionData.GroupPermissions, overridePayload.GroupPermissions...)
			permissionData.GroupPermissions = removeDuplicate(permissionData.GroupPermissions)
			sort.Slice(permissionData.GroupPermissions, func(i, j int) bool {
				return permissionData.GroupPermissions[i].Id < permissionData.GroupPermissions[j].Id
			})

			// Instance Type Permissions
			if len(overridePayload.InstanceTypePermissions) > 0 {
				for indx, perm := range permissionData.InstanceTypePermissions {
					for _, overperm := range overridePayload.InstanceTypePermissions {
						if perm.Id == overperm.Id {
							permissionData.InstanceTypePermissions[indx].Access = overperm.Access
						}
					}
				}
			}
			permissionData.InstanceTypePermissions = append(permissionData.InstanceTypePermissions, overridePayload.InstanceTypePermissions...)
			permissionData.InstanceTypePermissions = removeDuplicate(permissionData.InstanceTypePermissions)
			sort.Slice(permissionData.InstanceTypePermissions, func(i, j int) bool {
				return permissionData.InstanceTypePermissions[i].Id < permissionData.InstanceTypePermissions[j].Id
			})

			// Blueprint Permissions
			if len(overridePayload.BlueprintPermissions) > 0 {
				for indx, perm := range permissionData.BlueprintPermissions {
					for _, overperm := range overridePayload.BlueprintPermissions {
						if perm.Id == overperm.Id {
							permissionData.BlueprintPermissions[indx].Access = overperm.Access
						}
					}
				}
			}
			permissionData.BlueprintPermissions = append(permissionData.BlueprintPermissions, overridePayload.BlueprintPermissions...)
			permissionData.BlueprintPermissions = removeDuplicate(permissionData.BlueprintPermissions)
			sort.Slice(permissionData.BlueprintPermissions, func(i, j int) bool {
				return permissionData.BlueprintPermissions[i].Id < permissionData.BlueprintPermissions[j].Id
			})

			// Report Type Permissions
			if len(overridePayload.ReportTypePermissions) > 0 {
				for indx, perm := range permissionData.ReportTypePermissions {
					for _, overperm := range overridePayload.ReportTypePermissions {
						if perm.Code == overperm.Code {
							permissionData.ReportTypePermissions[indx].Access = overperm.Access
						}
					}
				}
			}
			permissionData.ReportTypePermissions = append(permissionData.ReportTypePermissions, overridePayload.ReportTypePermissions...)
			permissionData.ReportTypePermissions = removeDuplicate(permissionData.ReportTypePermissions)
			sort.Slice(permissionData.ReportTypePermissions, func(i, j int) bool {
				return permissionData.ReportTypePermissions[i].Code < permissionData.ReportTypePermissions[j].Code
			})

			// Persona Permissions
			if len(overridePayload.PersonaPermissions) > 0 {
				for indx, perm := range permissionData.PersonaPermissions {
					for _, overperm := range overridePayload.PersonaPermissions {
						if perm.Code == overperm.Code {
							permissionData.PersonaPermissions[indx].Access = overperm.Access
						}
					}
				}
			}
			permissionData.PersonaPermissions = append(permissionData.PersonaPermissions, overridePayload.PersonaPermissions...)
			permissionData.PersonaPermissions = removeDuplicate(permissionData.PersonaPermissions)
			sort.Slice(permissionData.PersonaPermissions, func(i, j int) bool {
				return permissionData.PersonaPermissions[i].Code < permissionData.PersonaPermissions[j].Code
			})

			// Catalog Item Type Permissions
			if len(overridePayload.CatalogItemTypePermissions) > 0 {
				for indx, perm := range permissionData.CatalogItemTypePermissions {
					for _, overperm := range overridePayload.CatalogItemTypePermissions {
						if perm.Id == overperm.Id {
							permissionData.CatalogItemTypePermissions[indx].Access = overperm.Access
						}
					}
				}
			}
			permissionData.CatalogItemTypePermissions = append(permissionData.CatalogItemTypePermissions, overridePayload.CatalogItemTypePermissions...)
			permissionData.CatalogItemTypePermissions = removeDuplicate(permissionData.CatalogItemTypePermissions)
			sort.Slice(permissionData.CatalogItemTypePermissions, func(i, j int) bool {
				return permissionData.CatalogItemTypePermissions[i].Id < permissionData.CatalogItemTypePermissions[j].Id
			})

			// VDI Pool Permissions
			if len(overridePayload.VdiPoolPermissions) > 0 {
				for indx, perm := range permissionData.VdiPoolPermissions {
					for _, overperm := range overridePayload.VdiPoolPermissions {
						if perm.Id == overperm.Id {
							permissionData.VdiPoolPermissions[indx].Access = overperm.Access
						}
					}
				}
			}
			permissionData.VdiPoolPermissions = append(permissionData.VdiPoolPermissions, overridePayload.VdiPoolPermissions...)
			permissionData.VdiPoolPermissions = removeDuplicate(permissionData.VdiPoolPermissions)
			sort.Slice(permissionData.VdiPoolPermissions, func(i, j int) bool {
				return permissionData.VdiPoolPermissions[i].Id < permissionData.VdiPoolPermissions[j].Id
			})

			// Workflow Permissions
			if len(overridePayload.WorkflowPermissions) > 0 {
				for indx, perm := range permissionData.WorkflowPermissions {
					for _, overperm := range overridePayload.WorkflowPermissions {
						if perm.Id == overperm.Id {
							permissionData.WorkflowPermissions[indx].Access = overperm.Access
						}
					}
				}
			}
			permissionData.WorkflowPermissions = append(permissionData.WorkflowPermissions, overridePayload.WorkflowPermissions...)
			permissionData.WorkflowPermissions = removeDuplicate(permissionData.WorkflowPermissions)
			sort.Slice(permissionData.WorkflowPermissions, func(i, j int) bool {
				return permissionData.WorkflowPermissions[i].Id < permissionData.WorkflowPermissions[j].Id
			})

			// Task Permissions
			if len(overridePayload.TaskPermissions) > 0 {
				for indx, perm := range permissionData.TaskPermissions {
					for _, overperm := range overridePayload.TaskPermissions {
						if perm.Id == overperm.Id {
							permissionData.TaskPermissions[indx].Access = overperm.Access
						}
					}
				}
			}
			permissionData.TaskPermissions = append(permissionData.TaskPermissions, overridePayload.TaskPermissions...)
			permissionData.TaskPermissions = removeDuplicate(permissionData.TaskPermissions)
			sort.Slice(permissionData.TaskPermissions, func(i, j int) bool {
				return permissionData.TaskPermissions[i].Id < permissionData.TaskPermissions[j].Id
			})
		}
	}

	jsonDoc, err := json.MarshalIndent(permissionData, "", "  ")
	log.Printf("API RESPONSE: %s", jsonDoc)

	if err != nil {
		return diag.Errorf("writing permission set: formatting JSON: %s", err)
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
	Id     int    `json:"id"`
	Access string `json:"access"`
}

type blueprintPermission struct {
	Id     int    `json:"id"`
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
	Id     int    `json:"id"`
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

func removeDuplicate[T comparable](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
