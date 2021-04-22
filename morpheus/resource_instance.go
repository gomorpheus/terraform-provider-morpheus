package morpheus

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"log"
	"strings"

	"github.com/gomorpheus/morpheus-go-sdk"
)

func resourceInstance() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus instance resource.",
		CreateContext: resourceInstanceCreate,
		ReadContext:   resourceInstanceRead,
		UpdateContext: resourceInstanceUpdate,
		DeleteContext: resourceInstanceDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the instance",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The name of the instance",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The user friendly description of the instance",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"cloud": {
				Description: "",
				Type:        schema.TypeString,
				Required:    true,
			},
			"group": {
				Description: "The group to provision the instance into",
				Type:        schema.TypeString,
				Required:    true,
			},
			"type": {
				Description: "The type of instance to provision",
				Type:        schema.TypeString,
				Required:    true,
			},
			"version": {
				Description: "",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"layout": {
				Description: "The layout to provision the instance from",
				Type:        schema.TypeString,
				Required:    true,
			},
			"plan": {
				Description: "The service plan associated with the instance",
				Type:        schema.TypeString,
				Required:    true,
			},
			"resource_pool": {
				Description: "",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"environment": {
				Description: "The environment to assign the instance to",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"tags": {
				Description: "Tags to assign to the instance",
				Type:        schema.TypeList,
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
			//todo: lookup user_group
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
						// "masked": {
						// 	Type:         schema.TypeBool,
						// 	Optional:     true,
						// 	Default:     false,
						// },
					},
				},
			},
			"evars": {
				Description: "The environment variables to assign to the instance",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "The name of the environment variable",
							Type:        schema.TypeString,
							Required:    true,
						},
						"value": {
							Description: "The value of the environment variable",
							Type:        schema.TypeString,
							Required:    true,
						},
						// "masked": {
						// 	Type:         schema.TypeBool,
						// 	Optional:     true,
						// 	Default:     false,
						// },
						// "export": {
						// 	Type:         schema.TypeBool,
						// 	Optional:     true,
						// 	Default:     true, // or false?
						// },
					},
				},
			},
			"volumes": {
				Description: "The instance volumes to create",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// "id": {
						// 	Type:         schema.TypeInt,
						// 	Optional:     true,
						// 	// Default:     -1,
						// },
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
						"datastore": {
							Description: "",
							Type:        schema.TypeString,
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
						// allows any value in the format network-*
						// otherwise looks up network by name or id
						"network": {
							Description: "The network to assign the network interface to",
							Type:        schema.TypeString,
							Optional:    true,
						},
						//todo: support look for these too,
						// until then, can just use network: "networkGroup-55"
						// "network_group": {
						// 	Type:         schema.TypeString,
						// 	Optional:     true,
						// },
						// "subnet": {
						// 	Type:         schema.TypeString,
						// 	Optional:     true,
						// },
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

/*

Here is the sequence of options and api requests for the select dropdowns.

Groups:
	http://localhost:8080/api/options/groups
Clouds:
	http://localhost:8080/api/options/clouds?groupId=1
Type:
	http://localhost:8080/api/options/instanceTypes?groupId=1
	http://localhost:8080/api/instance-types?code=activemq
	http://localhost:8080/api/instance-types/1
Name:
Description:
Environment:
	http://localhost:8080/api/options/environments
Tags (optional):
Version:
	http://localhost:8080/api/options/instanceVersions?groupId=1&cloudId=39&instanceTypeId=1
Layout:
	http://localhost:8080/api/options/layoutsForCloud?groupId=1&cloudId=39&instanceTypeId=1&version=5.11
Plan:
	http://localhost:8080/api/instances/service-plans?zoneId=39&layoutId=9&siteId=1

Volumes:

Interfaces:

Metadata:

Environment Variables:

the resulting payload looks like this:

{
  "zoneId": 40,
  "instance": {
    "name": "tftest",
    "cloud": "qa-azure",
    "site": {
      "id": 1
    },
    "type": "apache",
    "instanceType": {
      "code": "apache"
    },
    "layout": {
      "id": 1292
    },
    "plan": {
      "id": 187,
      "code": "azure.plan.westus.Basic_A0",
      "name": "Basic_A0 (1 Core, 0.75GB Memory) (westus)"
    }
  },
  "plan": {
    "id": 187,
    "code": "azure.plan.westus.Basic_A0",
    "name": "Basic_A0 (1 Core, 0.75GB Memory) (westus)"
  },
  "config": {
    "resourcePoolId": 205,
    "azuresecurityGroupId": null,
    "availabilitySet": null,
    "azurefloatingIp": "on",
    "createUser": true
  },
  "volumes": [
    {
      "id": -1,
      "rootVolume": true,
      "name": "root",
      "size": 32,
      "sizeId": null,
      "storageType": 14,
      "datastoreId": 1141
    }
  ],
  "networkInterfaces": [
    {
      "network": {
        "id": "network-413"
      }
    }
  ]
}
*/

// resourceInstanceCreate is a real doozy.
// We need to make the API easier to use. In the meantime,
// try to make this terraform resource easy to use, that means settings cloud and group
// with names easily like cloud: "My Cloud"
func resourceInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)

	// todo: needs to use /api/options/groups
	// group needs to be converted to instance.site.id
	var group *morpheus.Group
	groupName := d.Get("group").(string)
	if groupName == "" {
		return diag.Errorf("instance configuration requires 'group'")
	} else {
		groupResponse, groupErr := client.FindGroupByName(groupName)
		if groupErr != nil {
			return diag.FromErr(groupErr)
		}
		group = groupResponse.Result.(*morpheus.GetGroupResult).Group
	}

	// todo: needs to use /api/options/groups, heh
	// cloud needs to be converted to zoneId
	var cloud *morpheus.Cloud
	cloudName := d.Get("cloud").(string)
	if cloudName == "" {
		return diag.Errorf("instance configuration requires 'cloud'")
	} else {
		findResponse, findErr := client.FindCloudByName(cloudName)
		if findErr != nil {
			return diag.FromErr(findErr)
		}
		cloud = findResponse.Result.(*morpheus.GetCloudResult).Cloud
	}

	var optionResp *morpheus.Response
	var optionErr error

	// type
	var instanceType *morpheus.OptionSourceOption
	instanceTypeCode := d.Get("type").(string)
	if instanceTypeCode == "" {
		return diag.Errorf("instance configuration requires 'type'")
	} else {
		optionResp, optionErr = client.GetOptionSource("instanceTypes", &morpheus.Request{
			QueryParams: map[string]string{
				"groupId": int64ToString(group.ID),
				"cloudId": int64ToString(cloud.ID),
			},
		})
		if optionErr != nil {
			return diag.FromErr(optionErr)
		}
		optionSourceData := optionResp.Result.(*morpheus.GetOptionSourceResult).Data
		var matchingOptions []*morpheus.OptionSourceOption
		for i := 0; i < len(*optionSourceData); i++ {
			item := (*optionSourceData)[i] // .(optionSourceOption)
			// if item.Value.(string) == instanceTypeCode {
			if item.Name == instanceTypeCode || item.Code == instanceTypeCode || int64ToString(item.ID) == instanceTypeCode {
				matchingOptions = append(matchingOptions, &item)
			}
		}
		matchingOptionCount := len(matchingOptions)
		if matchingOptionCount != 1 {
			return diag.Errorf("Found %d Instance Types for '%v'", matchingOptionCount, instanceTypeCode)
		}
		instanceType = matchingOptions[0]
	}

	// version
	// the api never even looks for this, it just dictates the options for layout..
	var instanceVersion string
	if d.Get("version") != nil {
		instanceVersion = d.Get("version").(string)
	}

	// Layout
	var layout *morpheus.LayoutOption // OptionSourceOption
	layoutCode := d.Get("layout").(string)
	if layoutCode == "" {
		return diag.Errorf("instance configuration requires 'layout'")
	} else {
		//optionResp, optionErr = client.GetOptionSource("layoutsForCloud", &morpheus.Request{
		optionResp, optionErr = client.GetOptionSourceLayouts(&morpheus.Request{
			QueryParams: map[string]string{
				"groupId":        int64ToString(group.ID),
				"cloudId":        int64ToString(cloud.ID),
				"instanceTypeId": int64ToString(instanceType.ID),
				// "version": instanceVersion,
			},
		})
		if optionErr != nil {
			return diag.FromErr(optionErr)
		}
		//optionSourceData := optionResp.Result.(*morpheus.GetOptionSourceResult).Data
		optionSourceData := optionResp.Result.(*morpheus.GetOptionSourceLayoutsResult).Data
		var matchingOptions []*morpheus.LayoutOption
		for i := 0; i < len(*optionSourceData); i++ {
			item := (*optionSourceData)[i]
			if item.Name == layoutCode || item.Code == layoutCode || int64ToString(item.ID) == layoutCode {
				if instanceVersion != "" {
					if instanceVersion == item.Version {
						matchingOptions = append(matchingOptions, &item)
					}
				} else {
					matchingOptions = append(matchingOptions, &item)
				}
			}
		}
		matchingOptionCount := len(matchingOptions)
		if matchingOptionCount != 1 {
			return diag.Errorf("Found %d Layouts for '%v'. You may need to specify version.", matchingOptionCount, layoutCode)
		}
		layout = matchingOptions[0]
	}

	// Plan
	// plan needs to be converted to instance.plan.id
	var plan *morpheus.InstancePlan
	planCode := d.Get("plan").(string)
	if planCode == "" {
		return diag.Errorf("instance configuration requires 'plan'")
	} else {
		planResp, planErr := client.FindInstancePlanByCode(planCode, &morpheus.Request{
			QueryParams: map[string]string{
				"groupId":  int64ToString(group.ID),
				"zoneId":   int64ToString(cloud.ID),
				"layoutId": int64ToString(layout.ID),
			},
		})
		if planErr != nil {
			return diag.FromErr(planErr)
		}
		plan = planResp.Result.(*morpheus.GetInstancePlanResult).Plan
	}

	// config is a big map of who knows what
	var config map[string]interface{}
	if d.Get("config") != nil {
		config = d.Get("config").(map[string]interface{})
	}

	instancePayload := map[string]interface{}{
		"name": name,
		// "type": instanceTypeCode,
		"type": instanceType.Code,
		"site": map[string]interface{}{
			"id": group.ID,
		},
		"plan": map[string]interface{}{
			"id":   plan.ID,
			"code": plan.Code,
			"name": plan.Name,
		},
		"layout": map[string]interface{}{
			"id":   layout.ID,
			"code": layout.Code,
			"name": layout.Name,
		},
	}

	if d.Get("description") != nil {
		instancePayload["description"] = d.Get("description").(string)
	}

	// tags
	if d.Get("tags") != nil {
		// collect as a string, heh, you HAVE to use []interface{} for some reason..
		var tags []string
		tagList := d.Get("tags").([]interface{})
		for i := 0; i < len(tagList); i++ {
			tags = append(tags, tagList[i].(string))
		}
		// api should accept an array!
		// instancePayload["tags"] = tagList
		instancePayload["tags"] = strings.Join(tags, ", ")
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
	var resourcePool *morpheus.OptionSourceOption
	var resourcePoolIdStr string
	if d.Get("resource_pool") != nil {
		// look it up by name/externalId and pass config.resourcePoolId, ugh
		// /api/options/zonePools?groupId=456&siteId=456&zoneId=53&cloudId=53&instanceTypeId=1&planId=75&layoutId=3
		resourcePoolName := d.Get("resource_pool").(string)
		if resourcePoolName == "" {
			return diag.Errorf("instance configuration requires 'resource_pool'")
		} else {
			//optionResp, optionErr = client.GetOptionSource("resourcePoolsForCloud", &morpheus.Request{
			optionResp, optionErr = client.GetOptionSource("zonePools", &morpheus.Request{
				QueryParams: map[string]string{
					"groupId":        int64ToString(group.ID),
					"siteId":         int64ToString(group.ID),
					"cloudId":        int64ToString(cloud.ID),
					"instanceTypeId": int64ToString(instanceType.ID),
					"layoutId":       int64ToString(layout.ID),
					"planId":         int64ToString(plan.ID),
				},
			})
			if optionErr != nil {
				return diag.FromErr(optionErr)
			}
			// note: resourcePool.Value is a float64... why!?
			optionSourceData := optionResp.Result.(*morpheus.GetOptionSourceResult).Data
			var matchingOptions []*morpheus.OptionSourceOption
			for i := 0; i < len(*optionSourceData); i++ {
				item := (*optionSourceData)[i]
				if item.Name == resourcePoolName || item.ExternalId == resourcePoolName || int64ToString(int64(item.Value.(float64))) == resourcePoolName || int64ToString(item.ID) == resourcePoolName {
					matchingOptions = append(matchingOptions, &item)
				}
			}
			matchingOptionCount := len(matchingOptions)
			if matchingOptionCount != 1 {
				return diag.Errorf("Found %d Resource Pools for '%v'.", matchingOptionCount, resourcePoolName)
			}
			resourcePool = matchingOptions[0]
			resourcePoolIdStr = int64ToString(int64(resourcePool.Value.(float64)))
			//vmware and aws want different properties, one wants externalId too still? err
			config["resourcePoolId"] = resourcePoolIdStr // resourcePool.Value.(string)
			config["resourcePool"] = resourcePoolIdStr   // resourcePool.Value.(string)
		}
	}

	payload := map[string]interface{}{
		"zoneId":   cloud.ID,
		"instance": instancePayload,
		"config":   config,
	}

	// volumes
	if d.Get("volumes") != nil {

		// load datastore options
		//optionResp, optionErr = client.GetOptionSource("resourcePoolsForCloud", &morpheus.Request{
		datastoresResp, datastoresErr := client.GetOptionSource("datastores", &morpheus.Request{
			QueryParams: map[string]string{
				"groupId":        int64ToString(group.ID),
				"siteId":         int64ToString(group.ID),
				"cloudId":        int64ToString(cloud.ID),
				"zoneId":         int64ToString(cloud.ID),
				"instanceTypeId": int64ToString(instanceType.ID),
				"layoutId":       int64ToString(layout.ID),
				"planId":         int64ToString(plan.ID),
			},
		})
		if datastoresErr != nil {
			return diag.FromErr(datastoresErr)
		}

		// payload["volumes"] = d.Get("volumes").([]map[string]interface{})
		// need to make NetworkInterface an object...
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
			// todo: lookup storage_type by code
			if item["storage_type"] != nil {
				row["storageType"] = item["storage_type"] // .(int)
			}
			if item["datastore_id"] != nil {
				row["datastoreId"] = item["datastore_id"] // .(int)
			}
			// lookup datastore by name (partial name) and set datastoreId
			if item["datastore"] != nil {
				datastoreName := item["datastore"].(string)
				var datastore *morpheus.OptionSourceOption
				datastoreList := datastoresResp.Result.(*morpheus.GetOptionSourceResult).Data
				var matchingDatastores []*morpheus.OptionSourceOption
				for i := 0; i < len(*datastoreList); i++ {
					item := (*datastoreList)[i]
					if item.Name == datastoreName || int64ToString(item.ID) == datastoreName {
						matchingDatastores = append(matchingDatastores, &item)
					}
				}
				// match on partial name, because the name value returned here is silly
				// and includes a suffix like " - 2.2TB Free"
				if len(matchingDatastores) == 0 {
					for i := 0; i < len(*datastoreList); i++ {
						item := (*datastoreList)[i]
						if strings.HasPrefix(item.Name, datastoreName) {
							matchingDatastores = append(matchingDatastores, &item)
						}
					}
				}
				matchingDatastoreCount := len(matchingDatastores)
				if matchingDatastoreCount != 1 {
					return diag.Errorf("Found %d Datastores for '%v'.", matchingDatastoreCount, datastoreName)
				}
				datastore = matchingDatastores[0]

				row["datastoreId"] = datastore.ID // .(int)
			}
			volumes = append(volumes, row)
		}
		payload["volumes"] = volumes // .([]map[string]interface{})
	}

	// networkInterfaces
	if d.Get("interfaces") != nil {

		// load networking options
		// networkOptionsResp, networkOptionsErr := client.GetOptionSource("zoneNetworkOptions", &morpheus.Request{
		networkOptionsResp, networkOptionsErr := client.GetOptionSourceZoneNetworkOptions(&morpheus.Request{
			QueryParams: map[string]string{
				"groupId":        int64ToString(group.ID),
				"siteId":         int64ToString(group.ID),
				"cloudId":        int64ToString(cloud.ID),
				"zoneId":         int64ToString(cloud.ID),
				"instanceTypeId": int64ToString(instanceType.ID),
				"layoutId":       int64ToString(layout.ID),
				"poolId":         resourcePoolIdStr,
			},
		})
		if networkOptionsErr != nil {
			return diag.FromErr(networkOptionsErr)
		}

		// need to make NetworkInterface an object...
		interfaceList := d.Get("interfaces").([]interface{})
		var networkInterfaces []map[string]interface{}
		for i := 0; i < len(interfaceList); i++ {
			row := make(map[string]interface{})
			item := (interfaceList)[i].(map[string]interface{})
			// api expects "network":{"id":"network-45"}
			if item["network_id"] != nil {
				row["network"] = map[string]interface{}{
					"id": item["network_id"].(string),
				}
			}
			if item["network"] != nil {
				var networkId string
				networkName := item["network"].(string)
				if strings.HasPrefix(networkId, "network-") || strings.HasPrefix(networkId, "networkGroup-") || strings.HasPrefix(networkId, "subnet-") {
					//name passed as id like network- or networkGroup-
					networkId = networkName
				} else {
					// lookup network by name
					networkName := item["network"].(string)
					var network *morpheus.NetworkOption
					networkOptionsResult := networkOptionsResp.Result.(*morpheus.GetOptionSourceZoneNetworkOptionsResult)
					networkList := networkOptionsResult.Data.Networks
					var matchingNetworks []*morpheus.NetworkOption
					for i := 0; i < len(*networkList); i++ {
						item := (*networkList)[i]
						if item.Name == networkName || item.ID == networkName {
							matchingNetworks = append(matchingNetworks, &item)
						}
					}
					// match on partial name, because the name value returned here is silly
					// and includes a suffix like " - 2.2TB Free"
					if len(matchingNetworks) == 0 {
						for i := 0; i < len(*networkList); i++ {
							item := (*networkList)[i]
							if strings.HasPrefix(item.Name, networkName) {
								matchingNetworks = append(matchingNetworks, &item)
							}
						}
					}
					matchingNetworksCount := len(matchingNetworks)
					if matchingNetworksCount != 1 {
						return diag.Errorf("Found %d Networks for '%v'.", matchingNetworksCount, networkName)
					}
					network = matchingNetworks[0]

					networkId = network.ID

				}
				row["network"] = map[string]interface{}{
					"id": networkId,
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
	resp, err := client.CreateInstance(req)
	log.Printf("API REQUEST: %s", req) // debug
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.CreateInstanceResult)
	instance := result.Instance
	// Successfully created resource, now set id
	d.SetId(int64ToString(instance.ID))
	resourceInstanceRead(ctx, d, meta)
	return diags
}

func resourceInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	// d.Set("environment", instance.Environment)
	// d.Set("tags", instance.Tags)
	d.Set("config", instance.Config)
	// todo: more fields

	return diags
}

func resourceInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	name := d.Get("name").(string)
	code := d.Get("code").(string)
	location := d.Get("location").(string)
	// instances := d.Get("instances").([]interface{})

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"zone": map[string]interface{}{
				"name":     name,
				"code":     code,
				"location": location,
				// "instances": instances,
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
	return resourceInstanceRead(ctx, d, meta)
}

func resourceInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
			return nil
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)
			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)
	// result := resp.Result.(*morpheus.DeleteInstanceResult)
	d.SetId("") // implicit
	return diags
}
