package morpheus

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	//"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/gomorpheus/morpheusapi"
	"log"
	"fmt"
	"errors"
	"strings"
)

func resourceInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceInstanceCreate,
		Read:   resourceInstanceRead,
		Update: resourceInstanceUpdate,
		Delete: resourceInstanceDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"cloud": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"group": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"layout": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"plan": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_pool": {
				Type:         schema.TypeString,
				Optional:     true,
			},
			"environment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"config": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"create_user": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			//todo: lookup user_group			
			"user_group": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"metadata": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
						},
						"value": {
							Type:         schema.TypeString,
							Required:     true,
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
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
						},
						"value": {
							Type:         schema.TypeString,
							Required:     true,
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
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// "id": {
						// 	Type:         schema.TypeInt,
						// 	Optional:     true,
						// 	// Default:     -1,
						// },
						"root": {
							Type:         schema.TypeBool,
							Optional:     true,
						},
						"name": {
							Type:         schema.TypeString,
							Optional:     true,
						},
						"size": {
							Type:         schema.TypeInt,
							Optional:     true,
						},
						"size_id": {
							Type:         schema.TypeInt,
							Optional:     true,
						},
						"storage_type": {
							Type:         schema.TypeInt,
							Optional:     true,
						},
						"datastore": {
							Type:         schema.TypeString,
							Optional:     true,
						},
						"datastore_id": {
							Type:         schema.TypeInt,
							Optional:     true,
						},
					},
				},
			},
			"interfaces": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// allows any value in the format network-*
						// otherwise looks up network by name or id
						"network": {
							Type:         schema.TypeString,
							Optional:     true,
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
							Type:         schema.TypeString,
							Optional:     true,
						},
						"ip_mode": {
							Type:         schema.TypeString,
							Optional:     true,
						},
						"network_interface_type_id": {
							Type:         schema.TypeInt,
							Optional:     true,
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
    "name": "jd-azureapache1",
    "cloud": "qa-azure2",
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
func resourceInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*morpheusapi.Client)
	name := d.Get("name").(string)

	// todo: needs to use /api/options/groups
	// group needs to be converted to instance.site.id
	var group *morpheusapi.Group
	groupName := d.Get("group").(string)
	if groupName == "" {
		return errors.New("instance configuration requires 'group'")
	} else {
		groupResponse, groupErr := client.FindGroupByName(groupName)
		if groupErr != nil {
			return groupErr
		}
		group = groupResponse.Result.(*morpheusapi.GetGroupResult).Group
	}

	// todo: needs to use /api/options/groups, heh
	// cloud needs to be converted to zoneId
	var cloud *morpheusapi.Cloud
	cloudName := d.Get("cloud").(string)
	if cloudName == "" {
		return errors.New("instance configuration requires 'cloud'")
	} else {
		findResponse, findErr := client.FindCloudByName(cloudName)
		if findErr != nil {
			return findErr
		}
		cloud = findResponse.Result.(*morpheusapi.GetCloudResult).Cloud
	}

	var optionResp *morpheusapi.Response
	var optionErr error

	// type
	var instanceType *morpheusapi.OptionSourceOption
	instanceTypeCode := d.Get("type").(string)
	if instanceTypeCode == "" {
		return errors.New("instance configuration requires 'type'")
	} else {
		optionResp, optionErr = client.GetOptionSource("instanceTypes", &morpheusapi.Request{
			QueryParams:map[string]string{
		        "groupId": int64ToString(group.ID),
		        "cloudId": int64ToString(cloud.ID),
		    },
		})
		if optionErr != nil {
			return optionErr
		}
		optionSourceData := optionResp.Result.(*morpheusapi.GetOptionSourceResult).Data
		var matchingOptions []*morpheusapi.OptionSourceOption
		for i := 0; i < len(*optionSourceData); i++ {
			item := (*optionSourceData)[i] // .(optionSourceOption)
			// if item.Value.(string) == instanceTypeCode {
			if item.Name == instanceTypeCode ||  item.Code == instanceTypeCode || int64ToString(item.ID) == instanceTypeCode {
				matchingOptions = append(matchingOptions, &item)	
			}
		}
		matchingOptionCount := len(matchingOptions)
		if matchingOptionCount != 1 {
			return fmt.Errorf("Found %d Instance Types for '%v'", matchingOptionCount, instanceTypeCode)
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
	var layout *morpheusapi.LayoutOption // OptionSourceOption
	layoutCode := d.Get("layout").(string)
	if layoutCode == "" {
		return errors.New("instance configuration requires 'layout'")
	} else {
		//optionResp, optionErr = client.GetOptionSource("layoutsForCloud", &morpheusapi.Request{
		optionResp, optionErr = client.GetOptionSourceLayouts(&morpheusapi.Request{
			QueryParams:map[string]string{
		        "groupId": int64ToString(group.ID),
		        "cloudId": int64ToString(cloud.ID),
		        "instanceTypeId": int64ToString(instanceType.ID),
		        // "version": instanceVersion,
		    },
		})
		if optionErr != nil {
			return optionErr
		}
		//optionSourceData := optionResp.Result.(*morpheusapi.GetOptionSourceResult).Data
		optionSourceData := optionResp.Result.(*morpheusapi.GetOptionSourceLayoutsResult).Data
		var matchingOptions []*morpheusapi.LayoutOption
		for i := 0; i < len(*optionSourceData); i++ {
			item := (*optionSourceData)[i]
			if item.Name == layoutCode ||  item.Code == layoutCode || int64ToString(item.ID) == layoutCode {
				if (instanceVersion != "") {
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
			return fmt.Errorf("Found %d Layouts for '%v'. You may need to specify version.", matchingOptionCount, layoutCode)
		}
		layout = matchingOptions[0]
	}

	// Plan
	// plan needs to be converted to instance.plan.id
	var plan *morpheusapi.InstancePlan
	planCode := d.Get("plan").(string)
	if planCode == "" {
		return errors.New("instance configuration requires 'plan'")
	} else {
		planResp, planErr := client.FindInstancePlanByCode(planCode, &morpheusapi.Request{
			QueryParams:map[string]string{
		        "groupId": int64ToString(group.ID),
		        "zoneId": int64ToString(cloud.ID),
		        "layoutId": int64ToString(layout.ID),
		    },
		})
		if planErr != nil {
			return planErr
		}
		plan = planResp.Result.(*morpheusapi.GetInstancePlanResult).Plan
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
			"id": plan.ID,
			"code": plan.Code,
			"name": plan.Name,
		},
		"layout": map[string]interface{}{
			"id": layout.ID,
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
	var resourcePool *morpheusapi.OptionSourceOption
	var resourcePoolIdStr string
	if d.Get("resource_pool") != nil {
		// look it up by name/externalId and pass config.resourcePoolId, ugh
		// /api/options/zonePools?groupId=456&siteId=456&zoneId=53&cloudId=53&instanceTypeId=1&planId=75&layoutId=3
		resourcePoolName := d.Get("resource_pool").(string)
		if resourcePoolName == "" {
			return errors.New("instance configuration requires 'resource_pool'")
		} else {
			//optionResp, optionErr = client.GetOptionSource("resourcePoolsForCloud", &morpheusapi.Request{
			optionResp, optionErr = client.GetOptionSource("zonePools", &morpheusapi.Request{
				QueryParams:map[string]string{
			        "groupId": int64ToString(group.ID),
			        "siteId": int64ToString(group.ID),
			        "cloudId": int64ToString(cloud.ID),
			        "instanceTypeId": int64ToString(instanceType.ID),
			        "layoutId": int64ToString(layout.ID),
			        "planId": int64ToString(plan.ID),
			    },
			})
			if optionErr != nil {
				return optionErr
			}
			// note: resourcePool.Value is a float64... why!?
			optionSourceData := optionResp.Result.(*morpheusapi.GetOptionSourceResult).Data
			var matchingOptions []*morpheusapi.OptionSourceOption
			for i := 0; i < len(*optionSourceData); i++ {
				item := (*optionSourceData)[i]
				if item.Name == resourcePoolName || item.ExternalId == resourcePoolName || int64ToString(int64(item.Value.(float64))) == resourcePoolName || int64ToString(item.ID) == resourcePoolName {
					matchingOptions = append(matchingOptions, &item)
				}
			}
			matchingOptionCount := len(matchingOptions)
			if matchingOptionCount != 1 {
				return fmt.Errorf("Found %d Resource Pools for '%v'.", matchingOptionCount, resourcePoolName)
			}
			resourcePool = matchingOptions[0]
			resourcePoolIdStr = int64ToString(int64(resourcePool.Value.(float64)))
			//vmware and aws want different properties, one wants externalId too still? err
			config["resourcePoolId"] = resourcePoolIdStr // resourcePool.Value.(string)
			config["resourcePool"] = resourcePoolIdStr // resourcePool.Value.(string)
		}
	}

	payload := map[string]interface{}{
		"zoneId": cloud.ID,
		"instance": instancePayload,
		"config": config,
	}

	// volumes
	if d.Get("volumes") != nil {
		
		// load datastore options
		//optionResp, optionErr = client.GetOptionSource("resourcePoolsForCloud", &morpheusapi.Request{
		datastoresResp, datastoresErr := client.GetOptionSource("datastores", &morpheusapi.Request{
			QueryParams:map[string]string{
		        "groupId": int64ToString(group.ID),
		        "siteId": int64ToString(group.ID),
		        "cloudId": int64ToString(cloud.ID),
		        "zoneId": int64ToString(cloud.ID),
		        "instanceTypeId": int64ToString(instanceType.ID),
		        "layoutId": int64ToString(layout.ID),
		        "planId": int64ToString(plan.ID),
		    },
		})
		if datastoresErr != nil {
			return datastoresErr
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
				var datastore *morpheusapi.OptionSourceOption
				datastoreList := datastoresResp.Result.(*morpheusapi.GetOptionSourceResult).Data
				var matchingDatastores []*morpheusapi.OptionSourceOption
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
					return fmt.Errorf("Found %d Datastores for '%v'.", matchingDatastoreCount, datastoreName)
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
		// networkOptionsResp, networkOptionsErr := client.GetOptionSource("zoneNetworkOptions", &morpheusapi.Request{
		networkOptionsResp, networkOptionsErr := client.GetOptionSourceZoneNetworkOptions(&morpheusapi.Request{
			QueryParams:map[string]string{
		        "groupId": int64ToString(group.ID),
		        "siteId": int64ToString(group.ID),
		        "cloudId": int64ToString(cloud.ID),
		        "zoneId": int64ToString(cloud.ID),
		        "instanceTypeId": int64ToString(instanceType.ID),
		        "layoutId": int64ToString(layout.ID),
		        "poolId": resourcePoolIdStr,
		    },
		})
		if networkOptionsErr != nil {
			return networkOptionsErr
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
					var network *morpheusapi.NetworkOption
					networkOptionsResult := networkOptionsResp.Result.(*morpheusapi.GetOptionSourceZoneNetworkOptionsResult)
					networkList := networkOptionsResult.Data.Networks
					var matchingNetworks []*morpheusapi.NetworkOption
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
						return fmt.Errorf("Found %d Networks for '%v'.", matchingNetworksCount, networkName)
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

	req := &morpheusapi.Request{Body: payload}
	resp, err := client.CreateInstance(req)
	log.Printf("API REQUEST:", req) // debug
	if err != nil {
		log.Printf("API FAILURE:", resp, err)
		return err
	}
	log.Printf("API RESPONSE: ", resp)
	result := resp.Result.(*morpheusapi.CreateInstanceResult)
	instance := result.Instance
	// Successfully created resource, now set id
	d.SetId(int64ToString(instance.ID))
	return resourceInstanceRead(d, meta)
}

func resourceInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*morpheusapi.Client)
	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheusapi.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindInstanceByName(name)
	} else if id != "" {
		resp, err = client.GetInstance(toInt64(id), &morpheusapi.Request{})
		// todo: ignore 404 errors...
	} else {
		return errors.New("Instance cannot be read without name or id")
	}
	if err != nil {
		// 404 is ok?
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404:", resp, err)
			return nil
		} else {
			log.Printf("API FAILURE:", resp, err)
			return err
		}
	}
	log.Printf("API RESPONSE:", resp)

	// store resource data	
	result := resp.Result.(*morpheusapi.GetInstanceResult)
	instance := result.Instance
	if instance == nil {
		return fmt.Errorf("Instance not found in response data.") // should not happen
	}
	
	d.SetId(int64ToString(instance.ID))
	d.Set("name", instance.Name)
	d.Set("description", instance.Description)
	// d.Set("environment", instance.Environment)
	// d.Set("tags", instance.Tags)
	d.Set("config", instance.Config)
	// todo: more fields

	return nil
}

func resourceInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*morpheusapi.Client)
	id := d.Id()
	name := d.Get("name").(string)
	code := d.Get("code").(string)
	location := d.Get("location").(string)
	// instances := d.Get("instances").([]interface{})

	req := &morpheusapi.Request{
		Body: map[string]interface{}{
			"zone": map[string]interface{}{
				"name": name,
				"code": code,
				"location": location,
				// "instances": instances,
			},
		},
	}
	resp, err := client.UpdateInstance(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE:", resp, err)
		return err
	}
	log.Printf("API RESPONSE: ", resp)
	result := resp.Result.(*morpheusapi.UpdateInstanceResult)
	instance := result.Instance
	// Successfully updated resource, now set id
	d.SetId(int64ToString(instance.ID))
	return resourceInstanceRead(d, meta)
}

func resourceInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*morpheusapi.Client)
	id := d.Id()
	req := &morpheusapi.Request{}
	// req := &morpheusapi.Request{
	// 	QueryParams:map[string]string{
	// 		"force": string(USE_FORCE),
	// 	},
	// }
	//return errors.New("oh no...")
	resp, err := client.DeleteInstance(toInt64(id), req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404:", resp, err)
			return nil
		} else {
			log.Printf("API FAILURE:", resp, err)
			return err
		}
	}
	log.Printf("API RESPONSE:", resp)
	// result := resp.Result.(*morpheusapi.DeleteInstanceResult)
	//d.setId("") // implicit
	return nil
}
