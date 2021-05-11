package morpheus

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVsphereInstance() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus instance resource.",
		CreateContext: resourceVsphereInstanceCreate,
		ReadContext:   resourceVsphereInstanceRead,
		UpdateContext: resourceVsphereInstanceUpdate,
		DeleteContext: resourceVsphereInstanceDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(45 * time.Minute),
			Delete: schema.DefaultTimeout(45 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the instance",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The name of the instance",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"description": {
				Description: "The user friendly description of the instance",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"cloud_id": {
				Description: "The ID of the cloud associated with the instance",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"group_id": {
				Description: "The ID of the group associated with the instance",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"instance_type_id": {
				Description: "The type of instance to provision",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"version": {
				Description: "",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"instance_layout_id": {
				Description: "The layout to provision the instance from",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"plan_id": {
				Description: "The service plan associated with the instance",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"resource_pool_id": {
				Description: "",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"environment": {
				Description: "The environment to assign the instance to",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"labels": {
				Type:        schema.TypeList,
				Description: "",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tags": {
				Description: "Tags to assign to the instance",
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"config": {
				Description: "The instance configuration options",
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"create_user": {
				Description: "",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"user_group": {
				Description: "",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"metadata": {
				Description: "Metadata assigned to the instance",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "The name of the metadata",
							Type:        schema.TypeString,
							Required:    true,
						},
						"value": {
							Description: "The value of the metadata",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
			"evar": {
				Type:        schema.TypeList,
				Description: "",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "",
							Optional:    true,
						},
						"value": {
							Type:        schema.TypeString,
							Description: "",
							Optional:    true,
						},
						"export": {
							Type:        schema.TypeBool,
							Description: "",
							Optional:    true,
						},
						"masked": {
							Type:        schema.TypeBool,
							Description: "",
							Optional:    true,
						},
					},
				},
			},
			"volumes": {
				Description: "The instance volumes to create",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"root": {
							Description: "",
							Type:        schema.TypeBool,
							Optional:    true,
						},
						"name": {
							Description: "The name/type of the LV being created",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"size": {
							Description: "The size of the LV being created",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"size_id": {
							Description: "The ID of an existing LV to assign to the instance",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"storage_type": {
							Description: "The ID of the LV type",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"datastore_id": {
							Description: "The ID of the datastore",
							Type:        schema.TypeInt,
							Optional:    true,
						},
					},
				},
			},
			"interfaces": {
				Description: "The instance network interfaces to create",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_id": {
							Description: "The network to assign the network interface to",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"ip_address": {
							Description: "",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"ip_mode": {
							Description: "",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"network_interface_type_id": {
							Description: "The network interface type",
							Type:        schema.TypeInt,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func resourceVsphereInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	group := d.Get("group_id").(int)
	cloud := d.Get("cloud_id").(int)

	// Plan
	planResp, err := client.GetPlan(int64(d.Get("plan_id").(int)), &morpheus.Request{})
	if err != nil {
		diag.FromErr(err)
	}
	planResult := planResp.Result.(*morpheus.GetPlanResult)
	plan := planResult.Plan

	// Instance Type
	instanceTypeResp, err := client.GetInstanceType(int64(d.Get("instance_type_id").(int)), &morpheus.Request{})
	if err != nil {
		diag.FromErr(err)
	}
	instanceTypeResult := instanceTypeResp.Result.(*morpheus.GetInstanceTypeResult)
	instanceTypeCode := instanceTypeResult.InstanceType.Code

	// Instance Layout
	instanceLayoutResp, err := client.GetInstanceLayout(int64(d.Get("instance_layout_id").(int)), &morpheus.Request{})
	if err != nil {
		diag.FromErr(err)
	}
	instanceLayoutResult := instanceLayoutResp.Result.(*morpheus.GetInstanceLayoutResult)
	instanceLayout := instanceLayoutResult.InstanceLayout
	log.Printf("Instance Layout: %v", instanceLayout)

	// config is a big map of who knows what
	var config map[string]interface{}
	if d.Get("config") != nil {
		config = d.Get("config").(map[string]interface{})
	}

	instancePayload := map[string]interface{}{
		"name": name,
		"type": instanceTypeCode,
		"site": map[string]interface{}{
			"id": group,
		},
		"plan": map[string]interface{}{
			"id":   plan.ID,
			"code": plan.Code,
			"name": plan.Name,
		},
		"layout": map[string]interface{}{
			"id":   instanceLayout.ID,
			"code": instanceLayout.Code,
			"name": instanceLayout.Name,
		},
	}

	if d.Get("description") != nil {
		instancePayload["description"] = d.Get("description").(string)
	}

	if d.Get("create_user") != nil {
		config["createUser"] = d.Get("create_user") //.(bool)
	}

	if d.Get("user_group") != nil {
		//todo: lookup user_group by name please
		//userGroupId := d.Get("user_group").(int64)
		userGroupPayload := map[string]interface{}{
			"id": d.Get("user_group"),
		}
		instancePayload["userGroup"] = userGroupPayload
	}

	// Resource Pool
	resourcePoolResp, err := client.GetResourcePool(int64(cloud), int64(d.Get("resource_pool_id").(int)), &morpheus.Request{})
	if err != nil {
		diag.FromErr(err)
	}
	resourcePoolResult := resourcePoolResp.Result.(*morpheus.GetResourcePoolResult)
	resourcePool := resourcePoolResult.ResourcePool
	log.Printf("Instance Layout: %v", resourcePool)
	config["resourcePoolId"] = resourcePool.ID

	payload := map[string]interface{}{
		"zoneId":   cloud,
		"instance": instancePayload,
		"config":   config,
	}

	// tags
	if d.Get("tags") != nil {
		tagsInput := d.Get("tags").(map[string]interface{})
		log.Println(tagsInput)
		var tags []map[string]interface{}
		for key, value := range tagsInput {
			tag := make(map[string]interface{})
			tag["name"] = key
			tag["value"] = value.(string)
			tags = append(tags, tag)
		}
		log.Printf("Tag JSON Input: %s", tags)
		payload["tags"] = tags
	}

	// labels
	if d.Get("labels") != nil {
		payload["labels"] = d.Get("labels")
	}

	// evars
	if d.Get("evar") != nil {
		evarList := d.Get("evar").([]interface{})
		var evars []map[string]interface{}
		// iterate over the array of evars
		for i := 0; i < len(evarList); i++ {
			row := make(map[string]interface{})
			evarconfig := evarList[i].(map[string]interface{})
			for k, v := range evarconfig {
				switch k {
				case "name":
					row["name"] = v.(string)
					log.Printf("evar name: %s", v.(string))
				case "value":
					row["value"] = v.(string)
					log.Printf("evar value: %s", v.(string))
				case "export":
					row["export"] = v.(bool)
					log.Printf("evar export: %t", v)
				case "masked":
					row["masked"] = v
					log.Printf("evar masked: %t", v)
				}
				log.Printf("evar string: %s", row)
			}
			evars = append(evars, row)
			log.Printf("evars payload: %s", evars)
		}
		payload["evars"] = evars
	}

	// volumes
	if d.Get("volumes") != nil {
		volumeList := d.Get("volumes").([]interface{})
		var volumes []map[string]interface{}
		for i := 0; i < len(volumeList); i++ {
			row := make(map[string]interface{})
			item := (volumeList)[i].(map[string]interface{})
			if item["root"] != nil {
				row["rootVolume"] = item["root"]
			}
			if item["name"] != nil {
				row["name"] = item["name"] // .(string)
			}
			if item["size"] != nil {
				row["size"] = item["size"] // .(int)
			}
			if item["size_id"] != nil {
				row["sizeId"] = item["size_id"] // .(int)
			}
			if item["storage_type"] != nil {
				row["storageType"] = item["storage_type"] // .(int)
			}
			if item["datastore_id"] != nil {
				row["datastoreId"] = item["datastore_id"] // .(int)
			}
			volumes = append(volumes, row)
		}
		payload["volumes"] = volumes // .([]map[string]interface{})
	}

	// networkInterfaces
	if d.Get("interfaces") != nil {
		interfaceList := d.Get("interfaces").([]interface{})
		var networkInterfaces []map[string]interface{}
		for i := 0; i < len(interfaceList); i++ {
			row := make(map[string]interface{})
			item := (interfaceList)[i].(map[string]interface{})
			if item["network_id"] != nil {
				row["network"] = map[string]interface{}{
					"id": fmt.Sprintf("network-%d", item["network_id"].(int)),
				}
			}
			if item["ip_address"] != nil {
				row["ipAddress"] = item["ip_address"] //.(string)
			}
			if item["ip_mode"] != nil {
				row["ipMode"] = item["ip_mode"] // .(string)
			}
			if item["network_interface_type_id"] != nil {
				row["networkInterfaceTypeId"] = item["network_interface_type_id"] //.(int)
			}
			networkInterfaces = append(networkInterfaces, row)
		}
		payload["networkInterfaces"] = networkInterfaces // .([]map[string]interface{})
	}

	req := &morpheus.Request{Body: payload}
	slcB, _ := json.Marshal(req.Body)
	log.Printf("API JSON REQUEST: %s", string(slcB))
	resp, err := client.CreateInstance(req)
	log.Printf("API REQUEST: %s", req) // debug
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.CreateInstanceResult)
	instance := result.Instance

	stateConf := &resource.StateChangeConf{
		Pending: []string{"provisioning", "starting", "stopping"},
		Target:  []string{"running", "failed"},
		Refresh: func() (interface{}, string, error) {
			instanceDetails, err := client.GetInstance(instance.ID, &morpheus.Request{})
			if err != nil {
				return "", "", err
			}
			result := instanceDetails.Result.(*morpheus.GetInstanceResult)
			instance := result.Instance
			return result, instance.Status, nil
		},
		Timeout:      3 * time.Hour,
		MinTimeout:   1 * time.Minute,
		Delay:        3 * time.Minute,
		PollInterval: 1 * time.Minute,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error creating instance: %s", err)
	}

	// Successfully created resource, now set id
	d.SetId(int64ToString(instance.ID))
	resourceVsphereInstanceRead(ctx, d, meta)
	return diags
}

func resourceVsphereInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindInstanceByName(name)
	} else if id != "" {
		resp, err = client.GetInstance(toInt64(id), &morpheus.Request{})
		// todo: ignore 404 errors...
	} else {
		return diag.Errorf("Instance cannot be read without name or id")
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
	result := resp.Result.(*morpheus.GetInstanceResult)
	instance := result.Instance
	if instance == nil {
		return diag.Errorf("Instance not found in response data.") // should not happen
	}

	d.SetId(int64ToString(instance.ID))
	d.Set("name", instance.Name)
	d.Set("description", instance.Description)
	d.Set("config", instance.Config)

	return diags
}

func resourceVsphereInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	name := d.Get("name").(string)
	code := d.Get("code").(string)
	location := d.Get("location").(string)

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"zone": map[string]interface{}{
				"name":     name,
				"code":     code,
				"location": location,
			},
		},
	}
	resp, err := client.UpdateInstance(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateInstanceResult)
	instance := result.Instance
	// Successfully updated resource, now set id
	d.SetId(int64ToString(instance.ID))
	return resourceVsphereInstanceRead(ctx, d, meta)
}

func resourceVsphereInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{
		QueryParams: map[string]string{},
	}
	if USE_FORCE {
		req.QueryParams["force"] = "true"
	}
	resp, err := client.DeleteInstance(toInt64(id), req)
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
