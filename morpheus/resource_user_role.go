package morpheus

import (
	"context"
	"encoding/json"

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
			"permissions": {
				Type:        schema.TypeSet,
				Description: "Role permissions",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"global_groups_permission": {
							Type:        schema.TypeString,
							Description: "The name of the Morpheus plan.",
							Optional:    true,
							Computed:    true,
						},
						"global_instance_types_permission": {
							Type:        schema.TypeString,
							Description: "The name of the Morpheus plan.",
							Optional:    true,
							Computed:    true,
						},
						"global_blueprints_permission": {
							Type:        schema.TypeString,
							Description: "The name of the Morpheus plan.",
							Optional:    true,
							Computed:    true,
						},
						"global_report_types_permission": {
							Type:        schema.TypeString,
							Description: "The name of the Morpheus plan.",
							Optional:    true,
							Computed:    true,
						},
						"global_catalog_item_types_permission": {
							Type:        schema.TypeString,
							Description: "The name of the Morpheus plan.",
							Optional:    true,
							Computed:    true,
						},
						"global_vdi_pools_permission": {
							Type:        schema.TypeString,
							Description: "The name of the Morpheus plan.",
							Optional:    true,
							Computed:    true,
						},
						"default_persona": {
							Type:        schema.TypeString,
							Description: "The name of the Morpheus plan.",
							Optional:    true,
							Computed:    true,
						},
						"feature_permission": {
							Type:        schema.TypeList,
							Description: "The environment variables to create",
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
										Description: "The value of the environment variable",
										Optional:    true,
									},
								},
							},
						},
						"group_permission": {
							Type:        schema.TypeList,
							Description: "The environment variables to create",
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
										Description: "The value of the environment variable",
										Optional:    true,
									},
								},
							},
						},
						"instance_type_permission": {
							Type:        schema.TypeList,
							Description: "The environment variables to create",
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
										Description: "The value of the environment variable",
										Optional:    true,
									},
								},
							},
						},
						"blueprint_permission": {
							Type:        schema.TypeList,
							Description: "The environment variables to create",
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
										Description: "The value of the environment variable",
										Optional:    true,
									},
								},
							},
						},
						"report_type_permission": {
							Type:        schema.TypeList,
							Description: "The environment variables to create",
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
										Description: "The value of the environment variable",
										Optional:    true,
									},
								},
							},
						},
						"persona_permission": {
							Type:        schema.TypeList,
							Description: "The environment variables to create",
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
										Description: "The value of the environment variable",
										Optional:    true,
									},
								},
							},
						},
						"catalog_item_type_permission": {
							Type:        schema.TypeList,
							Description: "The environment variables to create",
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
										Description: "The value of the environment variable",
										Optional:    true,
									},
								},
							},
						},
					},
				},
				Optional: true,
				MaxItems: 1,
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

	log.Printf("PERMISSIONS: %s", d.Get("permissions"))
	req := &morpheus.Request{
		Body: map[string]interface{}{
			"role": map[string]interface{}{
				"authority":         d.Get("name").(string),
				"description":       d.Get("description").(string),
				"roleType":          "user",
				"multitenant":       d.Get("multitenant_role").(bool),
				"multitenantLocked": d.Get("multitenant_locked").(bool),
			},
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

	perms := d.Get("permissions")

	featureRequest := &morpheus.Request{
		Body: map[string]interface{}{
			"permissionCode": "global_vdi_pools_permission",
			"access":         "full",
		},
	}

	permResp, permErr := client.UpdateRoleFeaturePermission(role.Role.ID, featureRequest)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", permResp, permErr)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", permResp)

	// Successfully created resource, now set id
	d.SetId(int64ToString(role.Role.ID))

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
	role := result.Role

	d.SetId(int64ToString(role.ID))
	d.Set("name", role.Authority)
	d.Set("description", role.Description)
	d.Set("multitenant_role", role.MultiTenant)
	d.Set("multitenant_locked", role.MultiTenantLocked)

	return diags
}

func resourceUserRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	log.Printf("USER ROLE ID: %d", toInt64(id))

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"role": map[string]interface{}{
				"authority":         d.Get("name").(string),
				"description":       d.Get("description").(string),
				"roleType":          "user",
				"multitenant":       d.Get("multitenant_role").(bool),
				"multitenantLocked": d.Get("multitenant_locked").(bool),
			},
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

type CreateRoleResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	Role    *morpheus.Role    `json:"role"`
}
