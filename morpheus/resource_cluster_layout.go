package morpheus

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceClusterLayout() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus cluster layout resource",
		CreateContext: resourceClusterLayoutCreate,
		ReadContext:   resourceClusterLayoutRead,
		UpdateContext: resourceClusterLayoutUpdate,
		DeleteContext: resourceClusterLayoutDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the cluster layout",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the cluster layout",
				Required:    true,
			},
			"version": {
				Type:        schema.TypeString,
				Description: "The version of the cluster layout",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the cluster layout",
				Optional:    true,
				Computed:    true,
			},
			"creatable": {
				Type:        schema.TypeBool,
				Description: "Whether the cluster layout can be used to create clusters or not",
				Optional:    true,
				Computed:    true,
			},
			"minimum_memory": {
				Type:        schema.TypeInt,
				Description: "The minimum amount of memory in bytes",
				Optional:    true,
				Computed:    true,
			},
			"cluster_type_id": {
				Type:        schema.TypeInt,
				Description: "The cluster type ID of the cluster layout",
				Required:    true,
			},
			"provision_type_id": {
				Type:        schema.TypeInt,
				Description: "The provision type ID of the cluster layout",
				Required:    true,
			},
			"enable_scaling": {
				Type:        schema.TypeBool,
				Description: "Whether to enable or disable horizontal scaling",
				Optional:    true,
				Computed:    true,
			},
			"install_docker": {
				Type:        schema.TypeBool,
				Description: "Whether to automatically install Docker or not",
				Optional:    true,
				Computed:    true,
			},
			"evar": {
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
						"value": {
							Type:        schema.TypeString,
							Description: "The environment variable value when the value can be in plaintext",
							Optional:    true,
							Computed:    true,
						},
						"masked_value": {
							Type:        schema.TypeString,
							Description: "The environment variable value when the value needs to be masked",
							Optional:    true,
							Sensitive:   true,
							Computed:    true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								if old == "" {
									return false
								}
								h := sha256.New()
								h.Write([]byte(new))
								sha256_hash := hex.EncodeToString(h.Sum(nil))
								return strings.EqualFold(strings.ToLower(old), strings.ToLower(sha256_hash))
							},
							DiffSuppressOnRefresh: true,
						},
						"export": {
							Type:        schema.TypeBool,
							Description: "Whether the environment variable is exported as an instance tag",
							Optional:    true,
						},
					},
				},
			},
			"option_type_ids": {
				Type:        schema.TypeList,
				Description: "A list of option type ids associated with the cluster layout",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == old
				},
				Computed: true,
			},
			/* WAITING ON API SUPPORT
			"spec_template_ids": {
				Type:        schema.TypeList,
				Description: "A list of spec templates ids associated with the cluster layout",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == old
				},
				Computed: true,
			},*/
			"workflow_id": {
				Type:        schema.TypeInt,
				Description: "Workflow ID to associate with the cluster layout",
				Optional:    true,
				Computed:    true,
			},
			"master_node_pool": {
				Type:        schema.TypeList,
				Description: "Master node configuration",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:        schema.TypeInt,
							Description: "The number of nodes",
							Required:    true,
						},
						"node_type_id": {
							Type:        schema.TypeInt,
							Description: "The id of the node type",
							Required:    true,
						},
						"priority_order": {
							Type:        schema.TypeInt,
							Description: "The priority order of the node type",
							Required:    true,
						},
					},
				},
			},
			"worker_node_pool": {
				Type:        schema.TypeList,
				Description: "Worker node configuration",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:        schema.TypeInt,
							Description: "The number of nodes",
							Required:    true,
						},
						"node_type_id": {
							Type:        schema.TypeInt,
							Description: "The id of the node type",
							Required:    true,
						},
						"priority_order": {
							Type:        schema.TypeInt,
							Description: "The priority order of the node type",
							Required:    true,
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

func resourceClusterLayoutCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	clusterLayout := make(map[string]interface{})
	clusterLayout["name"] = d.Get("name").(string)
	clusterLayout["computeVersion"] = d.Get("version").(string)
	clusterLayout["description"] = d.Get("description").(string)
	clusterLayout["creatable"] = d.Get("creatable").(bool)

	provisionType := make(map[string]interface{})
	provisionType["id"] = d.Get("provision_type_id").(int)
	clusterLayout["provisionType"] = provisionType

	groupType := make(map[string]interface{})
	groupType["id"] = d.Get("cluster_type_id").(int)
	clusterLayout["groupType"] = groupType

	clusterLayout["memoryRequirement"] = d.Get("minimum_memory").(int)

	if d.Get("workflow_id").(int) > 0 {
		taskSet := make(map[string]interface{})
		taskSet["id"] = d.Get("workflow_id").(int)
		var taskSets [1]map[string]interface{}
		taskSets[0] = taskSet
		clusterLayout["taskSets"] = taskSets
	}

	clusterLayout["hasAutoScale"] = d.Get("enable_scaling").(bool)
	clusterLayout["installContainerRuntime"] = d.Get("install_docker").(bool)

	// input types
	var optionTypes []map[string]interface{}
	if d.Get("option_type_ids") != nil {
		optionTypeList := d.Get("option_type_ids").([]interface{})
		// iterate over the array of tasks
		for i := 0; i < len(optionTypeList); i++ {
			row := make(map[string]interface{})
			row["id"] = optionTypeList[i]
			optionTypes = append(optionTypes, row)
		}
	}

	clusterLayout["optionTypes"] = optionTypes

	/* WAITING ON API SUPPORT
	// spec templates
	var specTemplates []map[string]interface{}
	if d.Get("spec_template_ids") != nil {
		specTemplateList := d.Get("spec_template_ids").([]interface{})
		// iterate over the array of spec templates
		for i := 0; i < len(specTemplateList); i++ {
			row := make(map[string]interface{})
			row["id"] = specTemplateList[i]
			specTemplates = append(specTemplates, row)
		}
	}

	clusterLayout["specTemplates"] = specTemplates
	*/

	clusterLayout["environmentVariables"] = parseClusterLayoutEnvironmentVariables(d.Get("evar").([]interface{}))

	clusterLayout["masters"] = parseClusterLayoutNodePools(d.Get("master_node_pool").([]interface{}))
	clusterLayout["workers"] = parseClusterLayoutNodePools(d.Get("worker_node_pool").([]interface{}))

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"layout": clusterLayout,
		},
	}

	resp, err := client.CreateClusterLayout(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		log.Println(err)
	}

	clusterLayoutID := fmt.Sprintf("%v", result["id"])

	// Successfully created resource, now set id
	d.SetId(clusterLayoutID)

	resourceClusterLayoutRead(ctx, d, meta)
	return diags
}

func resourceClusterLayoutRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindClusterLayoutByName(name)
	} else if id != "" {
		resp, err = client.GetClusterLayout(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Cluster layout cannot be read without name or id")
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
	var clusterLayout ClusterLayoutPayload
	json.Unmarshal(resp.Body, &clusterLayout)

	d.SetId(int64ToString(clusterLayout.ClusterLayout.ID))
	d.Set("name", clusterLayout.ClusterLayout.Name)
	d.Set("version", clusterLayout.ClusterLayout.ComputeVersion)
	d.Set("description", clusterLayout.ClusterLayout.Description)
	d.Set("creatable", clusterLayout.ClusterLayout.Creatable)
	d.Set("minimum_memory", clusterLayout.ClusterLayout.MemoryRequirement)
	d.Set("provision_type_id", clusterLayout.ClusterLayout.ProvisionType.ID)
	d.Set("cluster_type_id", clusterLayout.ClusterLayout.GroupType.ID)
	d.Set("install_docker", clusterLayout.ClusterLayout.InstallContainerRuntime)
	d.Set("enable_scaling", clusterLayout.ClusterLayout.HasAutoScale)

	if len(clusterLayout.ClusterLayout.TaskSets) > 0 {
		d.Set("workflow_id", clusterLayout.ClusterLayout.TaskSets[0].ID)
	}

	var evars []map[string]interface{}
	if clusterLayout.ClusterLayout.EnvironmentVariables != nil {
		// iterate over the array of environment variables
		for i := 0; i < len(clusterLayout.ClusterLayout.EnvironmentVariables); i++ {
			environmentVariable := clusterLayout.ClusterLayout.EnvironmentVariables[i]
			envPayload := make(map[string]interface{})
			envPayload["name"] = environmentVariable.Name
			if environmentVariable.Masked {
				envPayload["masked_value"] = environmentVariable.DefaultValueHash
			} else {
				envPayload["value"] = environmentVariable.DefaultValue
			}
			envPayload["export"] = environmentVariable.Export
			evars = append(evars, envPayload)
		}
	}
	d.Set("evar", evars)
	// inputs
	var inputs []int64
	if clusterLayout.ClusterLayout.OptionTypes != nil {
		// iterate over the array of option types
		for i := 0; i < len(clusterLayout.ClusterLayout.OptionTypes); i++ {
			input := clusterLayout.ClusterLayout.OptionTypes[i]
			inputs = append(inputs, int64(input.ID))
		}
	}
	stateInputs := matchTemplatesWithSchema(inputs, d.Get("option_type_ids").([]interface{}))
	d.Set("option_type_ids", stateInputs)

	/* WAITING ON API SUPPORT
	// spec tempaltes
	var specTemplates []int64
	if clusterLayout.ClusterLayout.SpecTemplates != nil {
		// iterate over the array of spec templates
		for i := 0; i < len(clusterLayout.ClusterLayout.SpecTemplates); i++ {
			specTemplate := clusterLayout.ClusterLayout.SpecTemplates[i]
			specTemplates = append(specTemplates, specTemplate.ID)
		}
	}

	stateSpecTemplates := matchTemplatesWithSchema(specTemplates, d.Get("spec_template_ids").([]interface{}))
	d.Set("spec_template_ids", stateSpecTemplates)
	*/

	return diags
}

func resourceClusterLayoutUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	clusterLayout := make(map[string]interface{})
	clusterLayout["name"] = d.Get("name").(string)
	clusterLayout["computeVersion"] = d.Get("version").(string)
	clusterLayout["description"] = d.Get("description").(string)
	clusterLayout["creatable"] = d.Get("creatable").(bool)

	provisionType := make(map[string]interface{})
	provisionType["id"] = d.Get("provision_type_id").(int)
	clusterLayout["provisionType"] = provisionType

	groupType := make(map[string]interface{})
	groupType["id"] = d.Get("cluster_type_id").(int)
	clusterLayout["groupType"] = groupType

	clusterLayout["memoryRequirement"] = d.Get("minimum_memory").(int)

	if d.Get("workflow_id").(int) > 0 {
		taskSet := make(map[string]interface{})
		taskSet["id"] = d.Get("workflow_id").(int)
		var taskSets [1]map[string]interface{}
		taskSets[0] = taskSet
		clusterLayout["taskSets"] = taskSets
	}

	// option types
	var optionTypes []map[string]interface{}
	if d.Get("option_type_ids") != nil {
		optionTypeList := d.Get("option_type_ids").([]interface{})
		// iterate over the array of option types
		for i := 0; i < len(optionTypeList); i++ {
			row := make(map[string]interface{})
			row["id"] = optionTypeList[i]
			optionTypes = append(optionTypes, row)
		}
	}

	clusterLayout["optionTypes"] = optionTypes

	/* WAITING ON API SUPPORT
	// spec templates
	var specTemplates []map[string]interface{}
	if d.Get("spec_template_ids") != nil {
		specTemplateList := d.Get("spec_template_ids").([]interface{})
		// iterate over the array of spec templates
		for i := 0; i < len(specTemplateList); i++ {
			row := make(map[string]interface{})
			row["id"] = specTemplateList[i]
			specTemplates = append(specTemplates, row)
		}
	}

	clusterLayout["specTemplates"] = specTemplates
	*/
	clusterLayout["hasAutoScale"] = d.Get("enable_scaling").(bool)
	clusterLayout["environmentVariables"] = parseClusterLayoutEnvironmentVariables(d.Get("evar").([]interface{}))
	clusterLayout["masters"] = parseClusterLayoutNodePools(d.Get("master_node_pool").([]interface{}))
	clusterLayout["workers"] = parseClusterLayoutNodePools(d.Get("worker_node_pool").([]interface{}))

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"layout": clusterLayout,
		},
	}

	resp, err := client.UpdateClusterLayout(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	//log.Printf("API RESPONSE: %s", resp)

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		log.Println(err)
	}

	clusterLayoutID := fmt.Sprintf("%v", result["id"])

	// Successfully updated resource, now set id
	d.SetId(clusterLayoutID)

	return resourceClusterLayoutRead(ctx, d, meta)
}

func resourceClusterLayoutDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteClusterLayout(toInt64(id), req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
			return nil
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)
			return diag.FromErr(err)
		}
	}
	//log.Printf("API RESPONSE: %s", resp)
	d.SetId("")
	return diags
}

func parseClusterLayoutNodePools(variables []interface{}) []map[string]interface{} {
	var nodepools []map[string]interface{}
	// iterate over the array of nodepools
	for i := 0; i < len(variables); i++ {
		row := make(map[string]interface{})
		nodepoolconfig := variables[i].(map[string]interface{})
		for k, v := range nodepoolconfig {
			switch k {
			case "count":
				row["nodeCount"] = v.(int)
			case "node_type_id":
				node_type := make(map[string]interface{})
				node_type["id"] = v.(int)
				row["containerType"] = node_type
			case "priority_order":
				row["priorityOrder"] = v.(int)
			}
		}
		nodepools = append(nodepools, row)
	}
	return nodepools
}

func parseClusterLayoutEnvironmentVariables(variables []interface{}) []map[string]interface{} {
	var evars []map[string]interface{}
	// iterate over the array of evars
	for i := 0; i < len(variables); i++ {
		row := make(map[string]interface{})
		evarconfig := variables[i].(map[string]interface{})
		for k, v := range evarconfig {
			switch k {
			case "name":
				row["name"] = v.(string)
				//row["evarName"] = v.(string)
				//row["valueType"] = "fixed"
			case "value":
				if v.(string) != "" {
					row["value"] = v.(string)
					row["masked"] = false
				}
			case "masked_value":
				if v.(string) != "" {
					row["value"] = v.(string)
					row["masked"] = true
				}
			case "export":
				row["export"] = v.(bool)
			}
		}
		evars = append(evars, row)
	}
	return evars
}

type ClusterLayoutPayload struct {
	ClusterLayout struct {
		ID      int64 `json:"id"`
		Account struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"account"`
		Name                    string `json:"name"`
		Description             string `json:"description"`
		Code                    string `json:"code"`
		ComputeVersion          string `json:"computeVersion"`
		HasAutoScale            bool   `json:"hasAutoScale"`
		Creatable               bool   `json:"creatable"`
		MemoryRequirement       int64  `json:"memoryRequirement"`
		InstallContainerRuntime bool   `json:"installContainerRuntime"`
		ProvisionType           struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"provisionType"`
		GroupType struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"groupType"`
		TaskSets []struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"taskSets"`
		SpecTemplates []struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"specTemplates"`
		OptionTypes []struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"optionTypes"`
		EnvironmentVariables []struct {
			EvarName         string `json:"evarName"`
			Name             string `json:"name"`
			DefaultValue     string `json:"defaultValue"`
			DefaultValueHash string `json:"defaultValueHash"`
			ValueType        string `json:"valueType"`
			Export           bool   `json:"export"`
			Masked           bool   `json:"masked"`
		} `json:"environmentVariables"`
		ComputeServers []struct {
			ID                      int64       `json:"id"`
			PriorityOrder           int64       `json:"priorityOrder"`
			NodeCount               int64       `json:"nodeCount"`
			NodeType                string      `json:"nodeType"`
			MinNodeCount            int64       `json:"minNodeCount"`
			MaxNodeCount            interface{} `json:"maxNodeCount"`
			DynamicCount            bool        `json:"dynamicCount"`
			InstallContainerRuntime bool        `json:"installContainerRuntime"`
			InstallStorageRuntime   bool        `json:"installStorageRuntime"`
			Name                    string      `json:"name"`
			Code                    string      `json:"code"`
			Category                interface{} `json:"category"`
			Config                  interface{} `json:"config"`
			ContainertType          struct {
				ID               int64       `json:"id"`
				Account          interface{} `json:"account"`
				Name             string      `json:"name"`
				ShortName        string      `json:"shortName"`
				Code             string      `json:"code"`
				ContainerVersion string      `json:"containerVersion"`
				ProvisionType    struct {
					ID   int64  `json:"id"`
					Name string `json:"name"`
					Code string `json:"code"`
				} `json:"provisionType"`
				VirtualImage interface{} `json:"virtualImage"`
				Category     interface{} `json:"category"`
				Config       struct {
				} `json:"config"`
				ContainerPorts []struct {
					ID                  int64       `json:"id"`
					Name                string      `json:"name"`
					Port                int64       `json:"port"`
					LoadBalanceProtocol interface{} `json:"loadBalanceProtocol"`
					ExportName          string      `json:"exportName"`
				} `json:"containerPorts"`
				ContainerScripts   []interface{} `json:"containerScripts"`
				ContainerTemplates []struct {
					ID   int64  `json:"id"`
					Name string `json:"name"`
				} `json:"containerTemplates"`
				EnvironmentVariables []interface{} `json:"environmentVariables"`
			} `json:"containerType"`
			ComputeServerType struct {
				ID             interface{} `json:"id"`
				Code           interface{} `json:"code"`
				Name           interface{} `json:"name"`
				Managed        interface{} `json:"managed"`
				ExternalDelete interface{} `json:"externalDelete"`
			} `json:"computeServerType"`
			ProvisionService interface{} `json:"provisionService"`
			PlanCategory     interface{} `json:"planCategory"`
			NamePrefix       interface{} `json:"namePrefix"`
			NameSuffix       string      `json:"nameSuffix"`
			ForceNameIndex   bool        `json:"forceNameIndex"`
			LoadBalance      bool        `json:"loadBalance"`
		}
	} `json:"layout"`
}
