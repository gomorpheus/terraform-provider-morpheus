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

func resourceTenantRole() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus tenant role resource (This resource requires Morpheus 6.0.4 or later).",
		CreateContext: resourceTenantRoleCreate,
		ReadContext:   resourceTenantRoleRead,
		UpdateContext: resourceTenantRoleUpdate,
		DeleteContext: resourceTenantRoleDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the tenant role",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the tenant role",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the tenant role",
				Optional:    true,
				Computed:    true,
			},
			"permission_set": {
				Type:             schema.TypeString,
				Description:      "The permission set JSON document",
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressEquivalentJsonDiffs,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTenantRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	data := PermissionSet{}
	if err := json.Unmarshal([]byte(d.Get("permission_set").(string)), &data); err != nil {
		return diag.FromErr(err)
	}

	var roleDefinition TenantRolePermissionPayload
	roleDefinition.Name = d.Get("name").(string)
	roleDefinition.Description = d.Get("description").(string)
	roleDefinition.RoleType = "account"
	roleDefinition.DefaultPersona.Code = data.DefaultPersona
	roleDefinition.GlobalCloudAccess = data.DefaultCloudPermission
	roleDefinition.GlobalInstanceTypeAccess = data.DefaultInstanceTypePermission
	roleDefinition.GlobalBlueprintAccess = data.DefaultBlueprintPermission
	roleDefinition.GlobalReportTypeAccess = data.DefaultReportTypePermission
	roleDefinition.GlobalCatalogItemTypeAccess = data.DefaultCatalogItemTypePermission
	roleDefinition.GlobalVDIPoolAccess = data.DefaultVdiPoolPermission
	roleDefinition.GlobalWorkflowAccess = data.DefaultWorkflowPermission
	roleDefinition.GlobalTaskAccess = data.DefaultTaskPermission
	roleDefinition.FeaturePermissions = data.FeaturePermissions
	roleDefinition.CloudPermissions = data.CloudPermissions
	roleDefinition.InstanceTypePermissions = data.InstanceTypePermissions
	roleDefinition.BlueprintPermissions = data.BlueprintPermissions
	roleDefinition.ReportTypePermissions = data.ReportTypePermissions
	roleDefinition.PersonaPermissions = data.PersonaPermissions
	roleDefinition.CatalogItemTypePermissions = data.CatalogItemTypePermissions
	roleDefinition.VdiPoolPermissions = data.VdiPoolPermissions
	roleDefinition.Workflows = data.WorkflowPermissions
	roleDefinition.Tasks = data.TaskPermissions

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"role": roleDefinition,
		},
	}

	jsonRequest, _ := json.Marshal(req.Body)
	log.Printf("API JSON REQUEST: %s", string(jsonRequest))

	resp, err := client.CreateRole(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var role CreateRoleResult
	if err := json.Unmarshal(resp.Body, &role); err != nil {
		return diag.FromErr(err)
	}

	// Successfully created resource, now set id
	d.SetId(int64ToString(role.Role.ID))

	resourceTenantRoleRead(ctx, d, meta)
	return diags
}

func resourceTenantRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	result := resp.Result.(*morpheus.GetRoleResult)
	role := result

	d.SetId(int64ToString(role.Role.ID))
	d.Set("name", role.Role.Authority)
	d.Set("description", role.Role.Description)

	// Convert the Morpheus API response into the permission set JSON format for comparison
	data := PermissionSet{}
	if err := json.Unmarshal([]byte(d.Get("permission_set").(string)), &data); err != nil {
		return diag.FromErr(err)
	}

	var featureList []string
	for _, feature := range data.FeaturePermissions {
		featureList = append(featureList, feature.Code)
	}

	var cloudList []int
	for _, cloud := range data.CloudPermissions {
		cloudList = append(cloudList, cloud.Id)
	}

	var instanceTypeList []int
	for _, instanceType := range data.InstanceTypePermissions {
		instanceTypeList = append(instanceTypeList, instanceType.Id)
	}

	var blueprintList []int
	for _, blueprint := range data.BlueprintPermissions {
		blueprintList = append(blueprintList, blueprint.Id)
	}

	var reportTypeList []string
	for _, reportType := range data.ReportTypePermissions {
		reportTypeList = append(reportTypeList, reportType.Code)
	}

	var personaList []string
	for _, persona := range data.PersonaPermissions {
		personaList = append(personaList, persona.Code)
	}

	var catalogItemTypeList []int
	for _, catalogItemType := range data.CatalogItemTypePermissions {
		catalogItemTypeList = append(catalogItemTypeList, catalogItemType.Id)
	}

	var vdiPoolList []int
	for _, vdiPool := range data.VdiPoolPermissions {
		vdiPoolList = append(vdiPoolList, vdiPool.Id)
	}

	var workflowList []int
	for _, workflow := range data.WorkflowPermissions {
		workflowList = append(workflowList, workflow.Id)
	}

	var taskList []int
	for _, task := range data.TaskPermissions {
		taskList = append(taskList, task.Id)
	}

	// Set the default permissions
	var permissionSet PermissionSet
	permissionSet.DefaultCloudPermission = role.GlobalZoneAccess
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

	// Cloud Permissions
	var cloudPermissions []cloudPermission
	for _, cloud := range role.Zones {
		if containsInt(cloudList, int(cloud.ID)) {
			var cloudPerm cloudPermission
			cloudPerm.Access = cloud.Access
			cloudPerm.Id = int(cloud.ID)
			cloudPermissions = append(cloudPermissions, cloudPerm)
		}
	}

	sort.Slice(cloudPermissions, func(i, j int) bool { return cloudPermissions[i].Id < cloudPermissions[j].Id })
	permissionSet.CloudPermissions = cloudPermissions

	// Instance Type Permissions
	var instanceTypePermissions []instanceTypePermission
	for _, instanceType := range role.InstanceTypePermissions {
		if containsInt(instanceTypeList, int(instanceType.ID)) {
			var instanceTypePerm instanceTypePermission
			instanceTypePerm.Access = instanceType.Access
			instanceTypePerm.Id = int(instanceType.ID)
			instanceTypePermissions = append(instanceTypePermissions, instanceTypePerm)
		}
	}

	sort.Slice(instanceTypePermissions, func(i, j int) bool { return instanceTypePermissions[i].Id < instanceTypePermissions[j].Id })
	permissionSet.InstanceTypePermissions = instanceTypePermissions

	// Blueprint Permissions
	var blueprintPermissions []blueprintPermission
	for _, blueprint := range role.AppTemplatePermissions {
		if containsInt(blueprintList, int(blueprint.ID)) {
			var blueprintPerm blueprintPermission
			blueprintPerm.Access = blueprint.Access
			blueprintPerm.Id = int(blueprint.ID)
			blueprintPermissions = append(blueprintPermissions, blueprintPerm)
		}
	}

	sort.Slice(blueprintPermissions, func(i, j int) bool { return blueprintPermissions[i].Id < blueprintPermissions[j].Id })
	permissionSet.BlueprintPermissions = blueprintPermissions

	// Report Type Permissions
	var reportTypePermissions []reportTypePermission
	for _, reportType := range role.ReportTypePermissions {
		if containsString(reportTypeList, reportType.Code) {
			var reportTypePerm reportTypePermission
			reportTypePerm.Access = reportType.Access
			reportTypePerm.Code = reportType.Code
			reportTypePermissions = append(reportTypePermissions, reportTypePerm)
		}
	}

	sort.Slice(reportTypePermissions, func(i, j int) bool { return reportTypePermissions[i].Code < reportTypePermissions[j].Code })
	permissionSet.ReportTypePermissions = reportTypePermissions

	// Persona Permissions
	var personaPermissions []personaPermission
	for _, persona := range role.PersonaPermissions {
		if containsString(personaList, persona.Code) {
			var personaPerm personaPermission
			personaPerm.Access = persona.Access
			personaPerm.Code = persona.Code
			personaPermissions = append(personaPermissions, personaPerm)
		}
	}

	sort.Slice(personaPermissions, func(i, j int) bool { return personaPermissions[i].Code < personaPermissions[j].Code })
	permissionSet.PersonaPermissions = personaPermissions

	// Catalog Item Type Permissions
	var catalogItemTypePermissions []catalogItemTypePermission
	for _, catalogItemType := range role.CatalogItemTypePermissions {
		if containsInt(catalogItemTypeList, int(catalogItemType.ID)) {
			var catalogItemTypePerm catalogItemTypePermission
			catalogItemTypePerm.Access = catalogItemType.Access
			catalogItemTypePerm.Id = int(catalogItemType.ID)
			catalogItemTypePermissions = append(catalogItemTypePermissions, catalogItemTypePerm)
		}
	}

	sort.Slice(catalogItemTypePermissions, func(i, j int) bool { return catalogItemTypePermissions[i].Id < catalogItemTypePermissions[j].Id })
	permissionSet.CatalogItemTypePermissions = catalogItemTypePermissions

	// VDI Pool Permissions
	var vdiPoolPermissions []vdiPoolPermission
	for _, vdiPool := range role.VDIPoolPermissions {
		if containsInt(vdiPoolList, int(vdiPool.ID)) {
			var vdiPoolPerm vdiPoolPermission
			vdiPoolPerm.Access = vdiPool.Access
			vdiPoolPerm.Id = int(vdiPool.ID)
			vdiPoolPermissions = append(vdiPoolPermissions, vdiPoolPerm)
		}
	}

	sort.Slice(vdiPoolPermissions, func(i, j int) bool { return vdiPoolPermissions[i].Id < vdiPoolPermissions[j].Id })
	permissionSet.VdiPoolPermissions = vdiPoolPermissions

	// Workflow Permissions
	var workflowPermissions []workflowPermission
	for _, workflow := range role.TaskSetPermissions {
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
		return diag.FromErr(err)
	}
	jsonString := string(jsonDoc)
	d.Set("permission_set", jsonString)

	return diags
}

func resourceTenantRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	data := PermissionSet{}
	if err := json.Unmarshal([]byte(d.Get("permission_set").(string)), &data); err != nil {
		return diag.FromErr(err)
	}

	var roleDefinition TenantRolePermissionPayload
	roleDefinition.Name = d.Get("name").(string)
	roleDefinition.Description = d.Get("description").(string)
	roleDefinition.RoleType = "account"
	roleDefinition.DefaultPersona.Code = data.DefaultPersona
	roleDefinition.GlobalCloudAccess = data.DefaultCloudPermission
	roleDefinition.GlobalInstanceTypeAccess = data.DefaultInstanceTypePermission
	roleDefinition.GlobalBlueprintAccess = data.DefaultBlueprintPermission
	roleDefinition.GlobalReportTypeAccess = data.DefaultReportTypePermission
	roleDefinition.GlobalCatalogItemTypeAccess = data.DefaultCatalogItemTypePermission
	roleDefinition.GlobalVDIPoolAccess = data.DefaultVdiPoolPermission
	roleDefinition.GlobalWorkflowAccess = data.DefaultWorkflowPermission
	roleDefinition.GlobalTaskAccess = data.DefaultTaskPermission
	roleDefinition.FeaturePermissions = data.FeaturePermissions
	roleDefinition.CloudPermissions = data.CloudPermissions
	roleDefinition.InstanceTypePermissions = data.InstanceTypePermissions
	roleDefinition.BlueprintPermissions = data.BlueprintPermissions
	roleDefinition.ReportTypePermissions = data.ReportTypePermissions
	roleDefinition.PersonaPermissions = data.PersonaPermissions
	roleDefinition.CatalogItemTypePermissions = data.CatalogItemTypePermissions
	roleDefinition.VdiPoolPermissions = data.VdiPoolPermissions
	roleDefinition.Workflows = data.WorkflowPermissions
	roleDefinition.Tasks = data.TaskPermissions

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"role": roleDefinition,
		},
	}

	resp, err := client.UpdateRole(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	var role CreateRoleResult
	if err := json.Unmarshal(resp.Body, &role); err != nil {
		return diag.FromErr(err)
	}

	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(role.Role.ID))
	return resourceTenantRoleRead(ctx, d, meta)
}

func resourceTenantRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

type TenantRolePermissionPayload struct {
	Name           string `json:"authority"`
	Description    string `json:"description"`
	Owner          int64  `json:"owner"`
	RoleType       string `json:"roleType"`
	DefaultPersona struct {
		Code string `json:"code"`
	} `json:"defaultPersona"`
	GlobalCloudAccess           string                      `json:"globalZoneAccess"`
	GlobalInstanceTypeAccess    string                      `json:"globalInstanceTypeAccess"`
	GlobalBlueprintAccess       string                      `json:"globalAppTemplateAccess"`
	GlobalReportTypeAccess      string                      `json:"globalReportTypeAccess"`
	GlobalPersonaAccess         string                      `json:"globalPersonaAccess"`
	GlobalCatalogItemTypeAccess string                      `json:"globalCatalogItemTypeAccess"`
	GlobalVDIPoolAccess         string                      `json:"globalVdiPoolAccess"`
	GlobalTaskAccess            string                      `json:"globalTaskAccess"`
	GlobalWorkflowAccess        string                      `json:"globalTaskSetAccess"`
	FeaturePermissions          []featurePermission         `json:"permissions"`
	CloudPermissions            []cloudPermission           `json:"zones"`
	InstanceTypePermissions     []instanceTypePermission    `json:"instanceTypes"`
	BlueprintPermissions        []blueprintPermission       `json:"appTemplates"`
	ReportTypePermissions       []reportTypePermission      `json:"reportTypes"`
	PersonaPermissions          []personaPermission         `json:"personas"`
	CatalogItemTypePermissions  []catalogItemTypePermission `json:"catalogItemTypes"`
	VdiPoolPermissions          []vdiPoolPermission         `json:"vdiPools"`
	Tasks                       []taskPermission            `json:"tasks"`
	Workflows                   []workflowPermission        `json:"taskSets"`
}
