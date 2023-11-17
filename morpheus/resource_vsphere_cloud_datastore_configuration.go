package morpheus

import (
	"context"
	"fmt"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceVSphereCloudDatastoreConfiguration() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus vSphere cloud datastore resource",
		CreateContext: resourceVSphereCloudDatastoreConfigurationCreate,
		ReadContext:   resourceVSphereCloudDatastoreConfigurationRead,
		UpdateContext: resourceVSphereCloudDatastoreConfigurationUpdate,
		DeleteContext: resourceVSphereCloudDatastoreConfigurationDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The id of the vSphere cloud datastore",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the vSphere cloud datastore",
				Required:    true,
			},
			"cloud_id": {
				Type:        schema.TypeInt,
				Description: "The id of the vSphere cloud",
				Required:    true,
			},
			"active": {
				Type:        schema.TypeBool,
				Description: "Whether the cloud datastore is active",
				Optional:    true,
				Computed:    true,
			},
			"group_access_all": {
				Type:        schema.TypeBool,
				Description: "Whether to grant all groups access to the datastore",
				Optional:    true,
				Computed:    true,
			},
			"group_access_ids": {
				Type:        schema.TypeSet,
				Description: "A list of group ids to grant access to the datastore",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"visibility": {
				Description:  "Determines whether the cloud datastore is visible in sub-tenants or not",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"private", "public"}, false),
				Default:      "private",
			},
			"tenant_access": {
				Type:        schema.TypeList,
				Description: "The tenant datastore access",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Description: "The id of the tenant",
							Optional:    true,
						},
						"default_store": {
							Type:        schema.TypeBool,
							Description: "Whether to mark the cloud datastore as a default store for this tenant",
							Optional:    true,
						},
						"image_target": {
							Type:        schema.TypeBool,
							Description: "Whether to mark the cloud datastore as an image target for this tenant",
							Optional:    true,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceVSphereCloudDatastoreConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	cloudId := d.Get("cloud_id").(int)
	name := d.Get("name").(string)
	// Find by name, then get by ID
	resp, err := client.ListCloudDatastores(int64(cloudId), &morpheus.Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	datastoreResult := resp.Result.(*morpheus.ListCloudDatastoresResult)
	if len(*datastoreResult.Datastores) == 0 {
		return diag.Errorf("Unable to find a datastore named %s", name)
	}
	datastoreId := (*datastoreResult.Datastores)[0].ID
	datastore := make(map[string]interface{})
	datastore["active"] = d.Get("active").(bool)
	datastore["visibility"] = d.Get("visibility").(string)
	resourcePermissions := make(map[string]interface{})
	resourcePermissions["all"] = d.Get("group_access_all").(bool)

	var groupIds []map[string]interface{}
	groupIdList := d.Get("group_access_ids").(*schema.Set).List()
	if len(groupIdList) > 0 {
		for _, v := range groupIdList {
			accessPayload := map[string]interface{}{
				"id": v,
			}
			groupIds = append(groupIds, accessPayload)
		}
	}
	resourcePermissions["sites"] = groupIds
	datastore["resourcePermissions"] = resourcePermissions

	var accounts []int
	var defaultTargets []int
	var defaultStores []int

	tenantData := d.Get("tenant_access").([]interface{})

	for i := 0; i < len(tenantData); i++ {
		evarconfig := tenantData[i].(map[string]interface{})
		for k, v := range evarconfig {
			switch k {
			case "id":
				accounts = append(accounts, v.(int))
			case "default_store":
				if v.(bool) {
					defaultStores = append(defaultStores, evarconfig["id"].(int))
				}
			case "image_target":
				if v.(bool) {
					defaultTargets = append(defaultTargets, evarconfig["id"].(int))
				}
			}
		}
	}
	var tenantPerm TenantPermission
	tenantPerm.Accounts = accounts
	tenantPerm.Defaultstore = defaultStores
	tenantPerm.Defaulttarget = defaultTargets

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"datastore":         datastore,
			"tenantPermissions": tenantPerm,
		},
	}

	resp, err = client.Execute(&morpheus.Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("/api/zones/%d/data-stores/%d", int64(cloudId), datastoreId),
		QueryParams: map[string]string{},
		Body:        req.Body,
		Result:      &morpheus.UpdateCloudDatastoreResult{},
	})
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.UpdateCloudDatastoreResult)
	datastoreUpdateResult := result.Datastore
	// Successfully created resource, now set id
	d.SetId(int64ToString(datastoreUpdateResult.ID))

	resourceVSphereCloudDatastoreConfigurationRead(ctx, d, meta)
	return diags
}

func resourceVSphereCloudDatastoreConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error

	// Find by name, then get by ID
	resp, err = client.ListCloudDatastores(4, &morpheus.Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("READ API RESPONSE: %s", resp)
	datastoreResult := resp.Result.(*morpheus.ListCloudDatastoresResult)
	if len(*datastoreResult.Datastores) == 0 {
		return diag.Errorf("Unable to find a datastore named %s", name)
	}
	datastore := (*datastoreResult.Datastores)[0]

	d.SetId(int64ToString(datastore.ID))
	d.Set("name", datastore.Name)
	d.Set("active", datastore.Active)
	d.Set("visibility", datastore.Visibility)
	d.Set("group_access_all", datastore.ResourcePermission.All)
	var group_ids []int
	for _, site := range datastore.ResourcePermission.Sites {
		group_ids = append(group_ids, site.ID)
	}
	d.Set("group_access_ids", group_ids)
	d.Set("tenant_access", parseDatastoreTenant(datastore.Tenants))

	return diags
}

func resourceVSphereCloudDatastoreConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	id := d.Id()
	cloudId := d.Get("cloud_id").(int)
	datastore := make(map[string]interface{})
	datastore["active"] = d.Get("active").(bool)
	datastore["visibility"] = d.Get("visibility").(string)
	resourcePermissions := make(map[string]interface{})
	resourcePermissions["all"] = d.Get("group_access_all").(bool)

	var groupIds []map[string]interface{}
	groupIdList := d.Get("group_access_ids").(*schema.Set).List()
	if len(groupIdList) > 0 {
		for _, v := range groupIdList {
			accessPayload := map[string]interface{}{
				"id": v,
			}
			groupIds = append(groupIds, accessPayload)
		}
	}
	resourcePermissions["sites"] = groupIds
	datastore["resourcePermissions"] = resourcePermissions

	var accounts []int
	var defaultTargets []int
	var defaultStores []int

	tenantData := d.Get("tenant_access").([]interface{})

	for i := 0; i < len(tenantData); i++ {
		evarconfig := tenantData[i].(map[string]interface{})
		for k, v := range evarconfig {
			switch k {
			case "id":
				accounts = append(accounts, v.(int))
			case "default_store":
				if v.(bool) {
					defaultStores = append(defaultStores, evarconfig["id"].(int))
				}
			case "image_target":
				if v.(bool) {
					defaultTargets = append(defaultTargets, evarconfig["id"].(int))
				}
			}
		}
	}
	var tenantPerm TenantPermission
	tenantPerm.Accounts = accounts
	tenantPerm.Defaultstore = defaultStores
	tenantPerm.Defaulttarget = defaultTargets

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"datastore":         datastore,
			"tenantPermissions": tenantPerm,
		},
	}

	resp, err := client.Execute(&morpheus.Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("/api/zones/%d/data-stores/%d", int64(cloudId), toInt64(id)),
		QueryParams: map[string]string{},
		Body:        req.Body,
		Result:      &morpheus.UpdateCloudDatastoreResult{},
	})

	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("UPDATE API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.UpdateCloudDatastoreResult)
	datastoreUpdateResult := result.Datastore

	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(datastoreUpdateResult.ID))
	return resourceVSphereCloudDatastoreConfigurationRead(ctx, d, meta)
}

func resourceVSphereCloudDatastoreConfigurationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	d.SetId("")
	return diags
}

func parseDatastoreTenant(variables Tenant) []map[string]interface{} {
	var tenantConfigs []map[string]interface{}
	// iterate over the array of tenantConfigs
	for i := 0; i < len(variables); i++ {
		row := make(map[string]interface{})
		row["id"] = variables[i].ID
		row["default_store"] = variables[i].DefaultStore
		row["image_target"] = variables[i].DefaultTarget
		tenantConfigs = append(tenantConfigs, row)
	}
	return tenantConfigs
}

type Tenant []struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	DefaultStore  bool   `json:"defaultStore"`
	DefaultTarget bool   `json:"defaultTarget"`
}

type TenantPermission struct {
	Accounts      []int `json:"accounts"`
	Defaulttarget []int `json:"defaultTarget"`
	Defaultstore  []int `json:"defaultStore"`
}
