package morpheus

import (
	"context"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceResourcePoolGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus resource pool group resource",
		CreateContext: resourceResourcePoolGroupCreate,
		ReadContext:   resourceResourcePoolGroupRead,
		UpdateContext: resourceResourcePoolGroupUpdate,
		DeleteContext: resourceResourcePoolGroupDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the resource pool group",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the resource pool group",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the resource pool group",
				Optional:    true,
				Computed:    true,
			},
			"mode": {
				Type:         schema.TypeString,
				Description:  "The load balancing mode of the resource pool group (roundrobin, availablecapacity)",
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"roundrobin", "availablecapacity"}, false),
			},
			"resource_pool_ids": {
				Type:        schema.TypeSet,
				Description: "A list of resource pool ids associated with the resource pool group",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"all_group_access": {
				Type:        schema.TypeBool,
				Description: "Whether all groups will be granted access to the resource pool group",
				Optional:    true,
			},
			"group_access": {
				Type:        schema.TypeList,
				Description: "A list of Morpheus group configuration to enable group access to the resource pool group",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:        schema.TypeInt,
							Description: "The ID of the Morpheus group to grant access to the resource pool group",
							Required:    true,
						},
						"default": {
							Type:        schema.TypeBool,
							Description: "Whether the resource pool group will be a default for the associated group",
							Required:    true,
						},
					},
				},
			},
			"visibility": {
				Type:         schema.TypeString,
				Description:  "Whether the resource pool group is visible in sub-tenants or not",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"private", "public"}, false),
				Default:      "private",
			},
			"tenant_ids": {
				Type:        schema.TypeSet,
				Description: "A list of tenant ids associated with the resource pool group",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceResourcePoolGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	resourcePermissions := make(map[string]interface{})
	tenantPermissions := make(map[string]interface{})

	tenantsPayload := make([]int, 0)
	if attr, ok := d.GetOk("tenant_ids"); ok {
		for _, s := range attr.(*schema.Set).List() {
			tenantsPayload = append(tenantsPayload, s.(int))
		}
	}

	tenantPermissions["accounts"] = tenantsPayload

	resourcePermissions["all"] = d.Get("all_group_access").(bool)
	// Group Access
	if d.Get("group_access") != "" {
		resourcePermissions["sites"] = parseGroupAccess(d.Get("group_access").([]interface{}))
	}

	poolsPayload := make([]int, 0)
	if attr, ok := d.GetOk("resource_pool_ids"); ok {
		for _, s := range attr.(*schema.Set).List() {
			poolsPayload = append(poolsPayload, s.(int))
		}
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"resourcePoolGroup": map[string]interface{}{
				"name":        d.Get("name").(string),
				"description": d.Get("description").(string),
				"mode":        d.Get("mode").(string),
				"visibility":  d.Get("visibility").(string),
				"pools":       poolsPayload,
			},
			"resourcePermissions": resourcePermissions,
			"tenantPermissions":   tenantPermissions,
		},
	}
	resp, err := client.CreateResourcePoolGroup(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateResourcePoolGroupResult)
	resourcePoolGroup := result.ResourcePoolGroup
	// Successfully created resource, now set id
	d.SetId(int64ToString(resourcePoolGroup.ID))

	resourceResourcePoolGroupRead(ctx, d, meta)
	return diags
}

func resourceResourcePoolGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindResourcePoolGroupByName(name)
	} else if id != "" {
		resp, err = client.GetResourcePoolGroup(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Resource pool group cannot be read without name or id")
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
	result := resp.Result.(*morpheus.GetResourcePoolGroupResult)
	resourcePoolGroup := result.ResourcePoolGroup
	d.SetId(int64ToString(resourcePoolGroup.ID))
	d.Set("name", resourcePoolGroup.Name)
	d.Set("description", resourcePoolGroup.Description)
	d.Set("mode", resourcePoolGroup.Mode)
	var resourcePools []int64
	if len(resourcePoolGroup.Pools) > 0 {
		resourcePools = append(resourcePools, resourcePoolGroup.Pools...)
	}
	d.Set("resource_pool_ids", resourcePools)
	d.Set("all_group_access", resourcePoolGroup.ResourcePermission.All)
	// Group Access
	var groupAccess []map[string]interface{}
	if len(resourcePoolGroup.ResourcePermission.Sites) != 0 {
		for _, group := range resourcePoolGroup.ResourcePermission.Sites {
			groupData := make(map[string]interface{})
			groupData["group_id"] = group.ID
			groupData["default"] = group.Default
			groupAccess = append(groupAccess, groupData)
		}
	}
	d.Set("group_access", groupAccess)
	// tenant ids
	var tenantIds []int64
	// iterate over the array of tasks
	for _, tenant := range resourcePoolGroup.Tenants {
		tenantIds = append(tenantIds, tenant.ID)
	}
	d.Set("tenant_ids", tenantIds)
	d.Set("visibility", resourcePoolGroup.Visibility)
	return diags
}

func resourceResourcePoolGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	resourcePermissions := make(map[string]interface{})
	tenantPermissions := make(map[string]interface{})

	tenantsPayload := make([]int, 0)
	if attr, ok := d.GetOk("tenant_ids"); ok {
		for _, s := range attr.(*schema.Set).List() {
			tenantsPayload = append(tenantsPayload, s.(int))
		}
	}

	tenantPermissions["accounts"] = tenantsPayload

	resourcePermissions["all"] = d.Get("all_group_access").(bool)
	// Group Access
	if d.Get("group_access") != "" {
		resourcePermissions["sites"] = parseGroupAccess(d.Get("group_access").([]interface{}))
	}

	poolsPayload := make([]int, 0)
	if attr, ok := d.GetOk("resource_pool_ids"); ok {
		for _, s := range attr.(*schema.Set).List() {
			poolsPayload = append(poolsPayload, s.(int))
		}
	}
	req := &morpheus.Request{
		Body: map[string]interface{}{
			"resourcePoolGroup": map[string]interface{}{
				"name":        d.Get("name").(string),
				"description": d.Get("description").(string),
				"mode":        d.Get("mode").(string),
				"visibility":  d.Get("visibility").(string),
				"pools":       poolsPayload,
			},
			"resourcePermissions": resourcePermissions,
			"tenantPermissions":   tenantPermissions,
		},
	}
	resp, err := client.UpdateResourcePoolGroup(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateResourcePoolGroupResult)
	resourcePoolGroup := result.ResourcePoolGroup
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(resourcePoolGroup.ID))
	return resourceResourcePoolGroupRead(ctx, d, meta)
}

func resourceResourcePoolGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteResourcePoolGroup(toInt64(id), req)
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

func parseGroupAccess(variables []interface{}) []map[string]interface{} {
	var accessData []map[string]interface{}
	// iterate over the array of group access
	for i := 0; i < len(variables); i++ {
		row := make(map[string]interface{})
		groupconfig := variables[i].(map[string]interface{})
		for k, v := range groupconfig {
			switch k {
			case "group_id":
				row["id"] = v.(int)
			case "default":
				row["default"] = v.(bool)
			}
		}
		accessData = append(accessData, row)
	}
	return accessData
}
