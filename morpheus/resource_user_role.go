package morpheus

import (
	"context"
	"encoding/json"
	"sort"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUserRole() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus user role resource.",
		CreateContext: resourceUserRoleCreate,
		ReadContext:   resourceUserRoleRead,
		UpdateContext: resourceUserRoleUpdate,
		DeleteContext: resourceUserRoleDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the user role",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the user role",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the user role",
				Optional:    true,
				Computed:    true,
			},
			"multitenant_role": {
				Type:        schema.TypeBool,
				Description: "Whether the user role is automatically copied into all existing subtenants as well as placed into a subtenant when created",
				Optional:    true,
				Computed:    true,
			},
			"multitenant_locked": {
				Type:        schema.TypeBool,
				Description: "Whether subtenants are allowed to branch off or modify this role.",
				Optional:    true,
				Computed:    true,
			},
			"permission_set": {
				Type:        schema.TypeString,
				Description: "",
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceUserRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	data := PermissionSet{}
	json.Unmarshal([]byte(d.Get("permission_set").(string)), &data)
	log.Println("PERMISSIONS: ", data)

	var roleDefinition RolePermissionPayload
	roleDefinition.Name = d.Get("name").(string)
	roleDefinition.Description = d.Get("description").(string)
	roleDefinition.RoleType = "user"
	roleDefinition.Multitenant = d.Get("multitenant_role").(bool)
	roleDefinition.MultitenantLocked = d.Get("multitenant_locked").(bool)
	roleDefinition.DefaultPersona.Code = data.DefaultPersona
	roleDefinition.GlobalGroupAccess = data.DefaultGroupPermission
	roleDefinition.GlobalInstanceTypeAccess = data.DefaultInstanceTypePermission
	roleDefinition.GlobalBlueprintAccess = data.DefaultBlueprintPermission
	roleDefinition.GlobalReportTypeAccess = data.DefaultReportTypePermission
	roleDefinition.GlobalCatalogItemTypeAccess = data.DefaultCatalogItemTypePermission
	roleDefinition.GlobalVDIPoolAccess = data.DefaultVdiPoolPermission
	roleDefinition.GlobalWorkflowAccess = data.DefaultWorkflowPermission
	roleDefinition.GlobalTaskAccess = data.DefaultTaskPermission
	roleDefinition.FeaturePermissions = data.FeaturePermissions
	roleDefinition.GroupPermissions = data.GroupPermissions
	roleDefinition.PersonaPermissions = data.PersonaPermissions
	roleDefinition.InstanceTypePermissions = data.InstanceTypePermissions
	roleDefinition.Tasks = data.TaskPermissions
	roleDefinition.Workflows = data.WorkflowPermissions

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"role": roleDefinition,
		},
	}

	resp, err := client.CreateRole(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var role CreateRoleResult
	json.Unmarshal(resp.Body, &role)

	//result := resp.Result.(*CreateRoleResult)
	//role := result.Role
	log.Printf("USER ROLE RESPONSE: %v", role.Role)

	// Successfully created resource, now set id
	d.SetId(int64ToString(role.Role.ID))

	// Set Group Access
	for _, it := range data.GroupPermissions {
		req := &morpheus.Request{
			Body: map[string]interface{}{
				"groupId": it.Id,
				"access":  it.Access,
			},
		}

		resp, err := client.UpdateRoleGroupAccess(role.Role.ID, req)
		if err != nil {
			log.Printf("API FAILURE: %s - %s", resp, err)
			return diag.FromErr(err)
		}
		log.Printf("API RESPONSE: %s", resp)
	}

	resourceUserRoleRead(ctx, d, meta)
	return diags
}

func resourceUserRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindRoleByName(name)
	} else if id != "" {
		resp, err = client.GetRole(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Role cannot be read without name or id")
	}

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

	// store resource data
	result := resp.Result.(*morpheus.GetRoleResult)
	role := result

	d.SetId(int64ToString(role.Role.ID))
	d.Set("name", role.Role.Authority)
	d.Set("description", role.Role.Description)
	d.Set("multitenant_role", role.Role.MultiTenant)
	d.Set("multitenant_locked", role.Role.MultiTenantLocked)

	data := PermissionSet{}
	json.Unmarshal([]byte(d.Get("permission_set").(string)), &data)
	log.Println("PERMISSIONS: ", data)

	var featureList []string
	for _, feature := range data.FeaturePermissions {
		featureList = append(featureList, feature.Code)
	}

	var groupList []int
	for _, group := range data.GroupPermissions {
		groupList = append(groupList, group.Id)
	}

	var instanceTypeList []string
	for _, instanceType := range data.InstanceTypePermissions {
		instanceTypeList = append(instanceTypeList, instanceType.Code)
	}

	var personaList []string
	for _, persona := range data.PersonaPermissions {
		personaList = append(personaList, persona.Code)
	}

	var workflowList []int
	for _, workflow := range data.WorkflowPermissions {
		workflowList = append(workflowList, workflow.Id)
	}

	var taskList []int
	for _, task := range data.TaskPermissions {
		taskList = append(taskList, task.Id)
	}

	var permissionSet PermissionSet
	permissionSet.DefaultGroupPermission = role.GlobalSiteAccess
	permissionSet.DefaultInstanceTypePermission = role.GlobalInstanceTypeAccess
	permissionSet.DefaultBlueprintPermission = role.GlobalAppTemplateAccess
	permissionSet.DefaultReportTypePermission = role.GlobalReportTypeAccess
	permissionSet.DefaultPersona = role.Role.DefaultPersona.Code
	permissionSet.DefaultCatalogItemTypePermission = role.GlobalCatalogItemTypeAccess
	permissionSet.DefaultVdiPoolPermission = role.GlobalVDIPoolAccess
	permissionSet.DefaultWorkflowPermission = role.GlobalTaskSetAccess
	permissionSet.DefaultTaskPermission = role.GlobalTaskAccess
	// Feature Permissions
	var featurePermissions []featurePermission
	for _, feature := range role.FeaturePermissions {
		if containsString(featureList, feature.Code) {
			var featurePerm featurePermission
			featurePerm.Access = feature.Access
			featurePerm.Code = feature.Code
			featurePermissions = append(featurePermissions, featurePerm)
		}
	}
	sort.Slice(featurePermissions, func(i, j int) bool { return featurePermissions[i].Code < featurePermissions[j].Code })

	permissionSet.FeaturePermissions = featurePermissions

	// Group Permissions
	var groupPermissions []groupPermission
	log.Println("GROUP DATA: ", role.Sites)
	for _, group := range role.Sites {
		if containsInt(groupList, int(group.ID)) {
			var groupPerm groupPermission
			groupPerm.Access = group.Access
			groupPerm.Id = int(group.ID)
			groupPermissions = append(groupPermissions, groupPerm)
		}
	}

	sort.Slice(groupPermissions, func(i, j int) bool { return groupPermissions[i].Id < groupPermissions[j].Id })

	log.Println("GROUP PERMS: ", groupPermissions)
	permissionSet.GroupPermissions = groupPermissions

	// Instance Type Permissions
	var instanceTypePermissions []instanceTypePermission
	log.Println("Instance Type DATA: ", role.InstanceTypePermissions)
	for _, instanceType := range role.InstanceTypePermissions {
		if containsString(instanceTypeList, instanceType.Code) {
			var instanceTypePerm instanceTypePermission
			instanceTypePerm.Access = instanceType.Access
			instanceTypePerm.Code = instanceType.Code
			instanceTypePermissions = append(instanceTypePermissions, instanceTypePerm)
		}
	}

	sort.Slice(instanceTypePermissions, func(i, j int) bool { return instanceTypePermissions[i].Code < instanceTypePermissions[j].Code })

	log.Println("Instance Type PERMS: ", instanceTypePermissions)
	permissionSet.InstanceTypePermissions = instanceTypePermissions

	// Persona Permissions
	var personaPermissions []personaPermission
	log.Println("Persona DATA: ", role.PersonaPermissions)
	for _, persona := range role.PersonaPermissions {
		if containsString(personaList, persona.Code) {
			var personaPerm personaPermission
			personaPerm.Access = persona.Access
			personaPerm.Code = persona.Code
			personaPermissions = append(personaPermissions, personaPerm)
		}
	}

	sort.Slice(personaPermissions, func(i, j int) bool { return personaPermissions[i].Code < personaPermissions[j].Code })

	log.Println("Persona PERMS: ", personaPermissions)
	permissionSet.PersonaPermissions = personaPermissions

	// Workflow Permissions
	var workflowPermissions []workflowPermission
	for _, workflow := range role.TaskPermissions {
		if containsInt(workflowList, int(workflow.ID)) {
			var workflowPerm workflowPermission
			workflowPerm.Access = workflow.Access
			workflowPerm.Id = int(workflow.ID)
			workflowPermissions = append(workflowPermissions, workflowPerm)
		}
	}

	sort.Slice(workflowPermissions, func(i, j int) bool { return workflowPermissions[i].Id < workflowPermissions[j].Id })

	permissionSet.WorkflowPermissions = workflowPermissions

	// Task Permissions
	var taskPermissions []taskPermission
	for _, task := range role.TaskPermissions {
		if containsInt(taskList, int(task.ID)) {
			var taskPerm taskPermission
			taskPerm.Access = task.Access
			taskPerm.Id = int(task.ID)
			taskPermissions = append(taskPermissions, taskPerm)
		}
	}

	sort.Slice(taskPermissions, func(i, j int) bool { return taskPermissions[i].Id < taskPermissions[j].Id })

	permissionSet.TaskPermissions = taskPermissions

	jsonDoc, err := json.MarshalIndent(permissionSet, "", "  ")
	log.Printf("API RESPONSE: %s", jsonDoc)

	if err != nil {
		log.Println("error")
		// should never happen if the above code is correct
		//		return diags.AppendErrorf(diags, "writing IAM Policy Document: formatting JSON: %s", err)
	}
	jsonString := string(jsonDoc)
	d.Set("permission_set", jsonString)

	return diags
}

func resourceUserRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	log.Printf("USER ROLE ID: %d", toInt64(id))

	data := PermissionSet{}
	json.Unmarshal([]byte(d.Get("permission_set").(string)), &data)
	log.Println("PERMISSIONS: ", data)

	var roleDefinition RolePermissionPayload
	roleDefinition.Name = d.Get("name").(string)
	roleDefinition.Description = d.Get("description").(string)
	roleDefinition.RoleType = "user"
	roleDefinition.Multitenant = d.Get("multitenant_role").(bool)
	roleDefinition.MultitenantLocked = d.Get("multitenant_locked").(bool)
	roleDefinition.DefaultPersona.Code = data.DefaultPersona
	roleDefinition.GlobalGroupAccess = data.DefaultGroupPermission
	roleDefinition.GlobalWorkflowAccess = data.DefaultWorkflowPermission
	roleDefinition.GlobalTaskAccess = data.DefaultTaskPermission
	roleDefinition.FeaturePermissions = data.FeaturePermissions
	roleDefinition.GroupPermissions = data.GroupPermissions
	roleDefinition.InstanceTypePermissions = data.InstanceTypePermissions
	roleDefinition.PersonaPermissions = data.PersonaPermissions
	roleDefinition.Tasks = data.TaskPermissions
	roleDefinition.Workflows = data.WorkflowPermissions

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"role": roleDefinition,
		},
	}

	log.Printf("API REQUEST: %s", req)
	resp, err := client.UpdateRole(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	var role CreateRoleResult
	json.Unmarshal(resp.Body, &role)

	//	result := resp.Result.(*morpheus.UpdateRoleResult)
	//	role := result.Role

	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(role.Role.ID))
	return resourceUserRoleRead(ctx, d, meta)
}

func resourceUserRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteRole(toInt64(id), req)
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

func containsString(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func containsInt(n []int, num int) bool {
	for _, v := range n {
		if v == num {
			return true
		}
	}

	return false
}

type CreateRoleResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	Role    *morpheus.Role    `json:"role"`
}

type RolePermissionPayload struct {
	Name              string `json:"authority"`
	Description       string `json:"description"`
	Owner             int64  `json:"owner"`
	RoleType          string `json:"roleType"`
	Multitenant       bool   `json:"multitenant"`
	MultitenantLocked bool   `json:"multitenantLocked"`
	DefaultPersona    struct {
		Code string `json:"code"`
	} `json:"defaultPersona"`
	GlobalGroupAccess           string                   `json:"globalSiteAccess"`
	GlobalInstanceTypeAccess    string                   `json:"globalInstanceTypeAccess"`
	GlobalBlueprintAccess       string                   `json:"globalAppTemplateAccess"`
	GlobalReportTypeAccess      string                   `json:"globalReportTypeAccess"`
	GlobalCatalogItemTypeAccess string                   `json:"globalCatalogItemTypeAccess"`
	GlobalVDIPoolAccess         string                   `json:"globalVdiPoolAccess"`
	GlobalTaskAccess            string                   `json:"globalTaskAccess"`
	GlobalWorkflowAccess        string                   `json:"globalTaskSetAccess"`
	GroupPermissions            []groupPermission        `json:"sites"`
	FeaturePermissions          []featurePermission      `json:"permissions"`
	InstanceTypePermissions     []instanceTypePermission `json:"instanceTypes"`
	PersonaPermissions          []personaPermission      `json:"personas"`
	VdiPoolPermissions          []vdiPoolPermission      `json:"vdipools"`
	Tasks                       []taskPermission         `json:"tasks"`
	Workflows                   []workflowPermission     `json:"taskSets"`
}
