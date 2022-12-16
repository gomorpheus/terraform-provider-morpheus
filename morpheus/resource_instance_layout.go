package morpheus

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceInstanceLayout() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus instance layout resource",
		CreateContext: resourceInstanceLayoutCreate,
		ReadContext:   resourceInstanceLayoutRead,
		UpdateContext: resourceInstanceLayoutUpdate,
		DeleteContext: resourceInstanceLayoutDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the instance layout",
				Computed:    true,
			},
			"instance_type_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the associated instance type",
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the instance layout",
				Required:    true,
			},
			"version": {
				Type:        schema.TypeString,
				Description: "The version of the instance layout",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The instance layout category",
				Optional:    true,
				Computed:    true,
			},
			"creatable": {
				Type:        schema.TypeBool,
				Description: "Whether the instance layout can be used to create an instance",
				Optional:    true,
				Computed:    true,
			},
			"technology": {
				Type:         schema.TypeString,
				Description:  "The technology of the instance layout (alibaba, amazon, arm, azure, maas, cloudFormation, docker, esxi, fusion, google, huawei, hyperv, kubernetes, kvm, nutanix, opentelekom, openstack, oraclecloud, oraclevm, scvmm, terraform, upcloud, vcd.vapp, vcd, vmware, workflow, xen)",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"alibaba", "amazon", "arm", "azure", "maas", "cloudFormation", "docker", "esxi", "fusion", "google", "huawei", "hyperv", "kubernetes", "kvm", "nutanix", "opentelekom", "openstack", "oraclecloud", "oraclevm", "scvmm", "terraform", "upcloud", "vcd.vapp", "vcd", "vmware", "workflow", "xen"}, false),
			},
			"minimum_memory": {
				Type:        schema.TypeInt,
				Description: "The memory requirement in megabytes",
				Optional:    true,
				Computed:    true,
			},
			"workflow_id": {
				Type:        schema.TypeInt,
				Description: "The id of the provisioning workflow associated with the instance layout",
				Optional:    true,
				Computed:    true,
			},
			"support_convert_to_managed": {
				Type:        schema.TypeBool,
				Description: "Whether the instance layout supports deployed instances to be converted to managed",
				Optional:    true,
				Computed:    true,
			},
			/* AWAITING API SUPPORT
			"enable_scaling": {
				Type:        schema.TypeBool,
				Description: "The instance layout category",
				Optional:    true,
			},
			*/
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
									return true
								}
								h := sha256.New()
								h.Write([]byte(new))
								sha256_hash := hex.EncodeToString(h.Sum(nil))
								return strings.EqualFold(strings.ToLower(old), strings.ToLower(sha256_hash))
							},
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
				Description: "A list of option type ids associated with the instance layout",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == old
				},
				Computed: true,
			},
			"node_type_ids": {
				Type:        schema.TypeList,
				Description: "A list of node type ids associated with the instance layout",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == old
				},
				Computed: true,
			},
			"spec_template_ids": {
				Type:        schema.TypeList,
				Description: "A list of spec template ids associated with the instance layout",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == old
				},
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceInstanceLayoutCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	instanceLayout := make(map[string]interface{})
	instanceLayout["name"] = d.Get("name").(string)
	instanceLayout["instanceVersion"] = d.Get("version").(string)
	instanceLayout["description"] = d.Get("description").(string)
	instanceLayout["creatable"] = d.Get("creatable").(bool)
	instanceLayout["provisionTypeCode"] = d.Get("technology").(string)
	instanceLayout["memoryRequirement"] = d.Get("minimum_memory").(int)
	instanceLayout["taskSetId"] = d.Get("workflow_id").(int)
	instanceLayout["supportsConvertToManaged"] = d.Get("support_convert_to_managed").(bool)
	//instanceLayout["hasAutoScale"] = d.Get("enable_scaling").(bool)
	instanceLayout["optionTypes"] = d.Get("option_type_ids")
	instanceLayout["environmentVariables"] = parseInstanceLayoutEnvironmentVariables(d.Get("evar").([]interface{}), d)

	switch d.Get("technology") {
	case "alibaba", "amazon", "azure", "maas", "docker", "esxi", "fusion", "google", "huawei", "hyperv", "kubernetes", "kvm", "nutanix", "opentelekom", "openstack", "oraclecloud", "oraclevm", "scvmm", "upcloud", "vcd.vapp", "vcd", "vmware", "xen":
		instanceLayout["containerTypes"] = d.Get("node_type_ids")
	case "arm", "cloudFormation", "terraform":
		instanceLayout["specTemplates"] = d.Get("spec_template_ids")
	case "workflow":
		break
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"instanceTypeLayout": instanceLayout,
		},
	}

	resp, err := client.CreateInstanceLayout(int64(d.Get("instance_type_id").(int)), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateInstanceLayoutResult)
	instanceLayoutResponse := result.InstanceLayout
	// Successfully created resource, now set id
	d.SetId(int64ToString(instanceLayoutResponse.ID))

	resourceInstanceLayoutRead(ctx, d, meta)
	return diags
}

func resourceInstanceLayoutRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindInstanceLayoutByName(name)
	} else if id != "" {
		resp, err = client.GetInstanceLayout(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Node type cannot be read without name or id")
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
	//log.Printf("API RESPONSE: %s", resp)

	// store resource data
	var instanceLayout InstanceLayoutPayload
	json.Unmarshal(resp.Body, &instanceLayout)

	d.SetId(int64ToString(instanceLayout.InstanceLayout.ID))
	d.Set("name", instanceLayout.InstanceLayout.Name)
	d.Set("version", instanceLayout.InstanceLayout.ContainerVersion)
	d.Set("description", instanceLayout.InstanceLayout.Description)
	d.Set("creatable", instanceLayout.InstanceLayout.Creatable)
	d.Set("minimum_memory", instanceLayout.InstanceLayout.MemoryRequirement)
	if len(instanceLayout.InstanceLayout.TaskSets) > 0 {
		d.Set("workflow_id", instanceLayout.InstanceLayout.TaskSets[0].ID)
	}
	d.Set("support_convert_to_managed", instanceLayout.InstanceLayout.SupportsConvertToManaged)

	var evars []map[string]interface{}
	if instanceLayout.InstanceLayout.EnvironmentVariables != nil {
		// iterate over the array of environment variables
		for i := 0; i < len(instanceLayout.InstanceLayout.EnvironmentVariables); i++ {
			environmentVariable := instanceLayout.InstanceLayout.EnvironmentVariables[i]
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
	if instanceLayout.InstanceLayout.OptionTypes != nil {
		// iterate over the array of option types
		for i := 0; i < len(instanceLayout.InstanceLayout.OptionTypes); i++ {
			input := instanceLayout.InstanceLayout.OptionTypes[i]
			inputs = append(inputs, int64(input.ID))
		}
	}
	stateInputs := matchTemplatesWithSchema(inputs, d.Get("option_type_ids").([]interface{}))
	d.Set("option_type_ids", stateInputs)

	// spec templates
	if d.Get("spec_template_ids") != nil {
		var specTemplates []int64
		if instanceLayout.InstanceLayout.SpecTemplates != nil {
			// iterate over the array of script templates
			for i := 0; i < len(instanceLayout.InstanceLayout.SpecTemplates); i++ {
				specTemplate := instanceLayout.InstanceLayout.SpecTemplates[i]
				specTemplates = append(specTemplates, specTemplate.ID)
			}
		}
		stateSpecTemplates := matchTemplatesWithSchema(specTemplates, d.Get("spec_template_ids").([]interface{}))
		d.Set("spec_template_ids", stateSpecTemplates)
	}

	// node types
	if d.Get("node_type_ids") != nil {
		var nodeTypes []int64
		if instanceLayout.InstanceLayout.ContainerTypes != nil {
			// iterate over the array of node types
			for i := 0; i < len(instanceLayout.InstanceLayout.ContainerTypes); i++ {
				nodeType := instanceLayout.InstanceLayout.ContainerTypes[i]
				nodeTypes = append(nodeTypes, nodeType.ID)
			}
		}
		stateNodeTypes := matchTemplatesWithSchema(nodeTypes, d.Get("node_type_ids").([]interface{}))
		d.Set("node_type_ids", stateNodeTypes)
	}

	return diags
}

func resourceInstanceLayoutUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	instanceLayout := make(map[string]interface{})
	instanceLayout["name"] = d.Get("name").(string)
	instanceLayout["instanceVersion"] = d.Get("version").(string)
	instanceLayout["description"] = d.Get("description").(string)
	instanceLayout["creatable"] = d.Get("creatable").(bool)
	instanceLayout["provisionTypeCode"] = d.Get("technology").(string)
	instanceLayout["memoryRequirement"] = d.Get("minimum_memory").(int)
	instanceLayout["taskSetId"] = d.Get("workflow_id").(int)
	instanceLayout["supportsConvertToManaged"] = d.Get("support_convert_to_managed").(bool)
	//instanceLayout["hasAutoScale"] = d.Get("enable_scaling").(bool)
	instanceLayout["optionTypes"] = d.Get("option_type_ids")
	instanceLayout["environmentVariables"] = parseInstanceLayoutEnvironmentVariables(d.Get("evar").([]interface{}), d)

	switch d.Get("technology") {
	case "alibaba", "amazon", "azure", "maas", "docker", "esxi", "fusion", "google", "huawei", "hyperv", "kubernetes", "kvm", "nutanix", "opentelekom", "openstack", "oraclecloud", "oraclevm", "scvmm", "upcloud", "vcd.vapp", "vcd", "vmware", "xen":
		instanceLayout["containerTypes"] = d.Get("node_type_ids")
	case "arm", "cloudFormation", "terraform":
		instanceLayout["specTemplates"] = d.Get("spec_template_ids")
	case "workflow":
		break
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"instanceTypeLayout": instanceLayout,
		},
	}

	resp, err := client.UpdateInstanceLayout(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateInstanceLayoutResult)
	instanceLayoutResponse := result.InstanceLayout
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(instanceLayoutResponse.ID))
	return resourceInstanceLayoutRead(ctx, d, meta)
}

func resourceInstanceLayoutDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteInstanceLayout(toInt64(id), req)
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

func parseInstanceLayoutEnvironmentVariables(variables []interface{}, d *schema.ResourceData) []map[string]interface{} {
	var evars []map[string]interface{}
	// iterate over the array of evars
	for i := 0; i < len(variables); i++ {
		row := make(map[string]interface{})
		evarconfig := variables[i].(map[string]interface{})
		for k, v := range evarconfig {
			switch k {
			case "name":
				row["name"] = v.(string)
				row["evarName"] = v.(string)
				row["valueType"] = "fixed"
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

type InstanceLayoutPayload struct {
	InstanceLayout struct {
		ID      int64 `json:"id"`
		Account struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"account"`
		Name                     string `json:"name"`
		Description              string `json:"description"`
		Code                     string `json:"code"`
		ContainerVersion         string `json:"instanceVersion"`
		Creatable                bool   `json:"creatable"`
		MemoryRequirement        int64  `json:"memoryRequirement"`
		SupportsConvertToManaged bool   `json:"supportsConvertToManaged"`
		ProvisionType            struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"provisionType"`
		TaskSets []struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"taskSets"`
		ContainerTypes []struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"containerTypes"`
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
	} `json:"instanceTypeLayout"`
}
