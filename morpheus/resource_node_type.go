package morpheus

import (
	"context"
	"encoding/json"
	"regexp"
	"strconv"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	shortNameCharactersWarning = "Short names may not contain spaces or underscores."
)

var shortNameCharacters, _ = regexp.Compile("^[^ _]*$")

func resourceNodeType() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus node type resource",
		CreateContext: resourceNodeTypeCreate,
		ReadContext:   resourceNodeTypeRead,
		UpdateContext: resourceNodeTypeUpdate,
		DeleteContext: resourceNodeTypeDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the node type",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the node type",
				Required:    true,
			},
			"short_name": {
				Type:         schema.TypeString,
				Description:  "The short name of the node type",
				Required:     true,
				ValidateFunc: validation.StringMatch(shortNameCharacters, shortNameCharactersWarning),
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "The organization labels associated with the script template (Only supported on Morpheus 5.5.3 or higher)",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"technology": {
				Type:         schema.TypeString,
				Description:  "The technology of the node type (alibaba, amazon, azure, maas, esxi, fusion, google, huawei, hyperv, kvm, nutanix, opentelekom, openstack, oraclecloud, oraclevm, scvmm, upcloud, vcd.vapp, vcd, vmware, xen)",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"alibaba", "amazon", "azure", "maas", "esxi", "fusion", "google", "huawei", "hyperv", "kvm", "nutanix", "opentelekom", "openstack", "oraclecloud", "oraclevm", "scvmm", "upcloud", "vcd.vapp", "vcd", "vmware", "xen"}, false),
			},
			/* AWAITING API SUPPORT TO AVOID DUPLICATE ENTRIES
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
							Description: "The value of the environment variable",
							Optional:    true,
						},
						"export": {
							Type:        schema.TypeBool,
							Description: "Whether the environment variable is exported as an instance tag",
							Optional:    true,
						},
						"masked": {
							Type:        schema.TypeBool,
							Description: "Whether the environment variable is masked for security purposes",
							Optional:    true,
						},
					},
				},
			},*/
			"version": {
				Type:        schema.TypeString,
				Description: "The version of the node type",
				Required:    true,
			},
			"virtual_image_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the virtual image associated with the node type",
				Optional:    true,
				Computed:    true,
			},
			/* AWAITING API SUPPORT
			"logs_folder": {
				Type:        schema.TypeString,
				Description: "The log folder associated with the node type",
				Optional:    true,
			},
			"config_folder": {
				Type:        schema.TypeString,
				Description: "The config folder associated with the node type",
				Optional:    true,
			},
			"deploy_folder": {
				Type:        schema.TypeString,
				Description: "The deploy folder associated with the node type",
				Optional:    true,
			},
			*/
			/* Waiting to add support for kubernetes
			"kubernetes_manifest": {
				Type:        schema.TypeString,
				Description: "The kubernetes manifest associated with the node type",
				Optional:    true,
				Computed:    true,
			},
			*/
			"service_port": {
				Type:        schema.TypeList,
				Description: "Service ports associated with the node type",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:        schema.TypeString,
							Description: "The port number of the service",
							Optional:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "The name of the service port",
							Optional:    true,
						},
						"protocol": {
							Type:        schema.TypeString,
							Description: "The load balancer protocol (HTTP, HTTPS, TCP)",
							Optional:    true,
						},
					},
				},
			},
			"extra_options": {
				Type:        schema.TypeMap,
				Description: "VMware custom options associated with the node type",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"script_template_ids": {
				Type:        schema.TypeList,
				Description: "A list of script template ids associated with the node type",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == old
				},
				Computed: true,
			},
			"file_template_ids": {
				Type:        schema.TypeList,
				Description: "A list of file template ids associated with the node type",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == old
				},
				Computed: true,
			},
			"category": {
				Type:        schema.TypeString,
				Description: "The node type category",
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceNodeTypeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)

	// extra options
	config := make(map[string]interface{})
	if d.Get("extra_options") != nil {
		config["extraOptions"] = d.Get("extra_options")
	}

	containerType := make(map[string]interface{})
	containerType["name"] = name

	containerType["shortName"] = d.Get("short_name").(string)
	containerType["containerVersion"] = d.Get("version").(string)
	containerType["provisionTypeCode"] = d.Get("technology").(string)
	//"environmentVariables": parseNodeTypeEnvironmentVariables(d.Get("evar").([]interface{})),
	if d.Get("virtual_image_id") != 0 {
		containerType["virtualImageId"] = d.Get("virtual_image_id").(int)
	}
	containerType["config"] = config
	containerType["containerPorts"] = parseNodeTypeServicePorts(d.Get("service_port").([]interface{}))
	containerType["scripts"] = d.Get("script_template_ids")
	containerType["containerTemplates"] = d.Get("file_template_ids")
	containerType["category"] = d.Get("category").(string)
	containerType["serverType"] = "vm"

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		for _, s := range attr.(*schema.Set).List() {
			labelsPayload = append(labelsPayload, s.(string))
		}
	}
	containerType["labels"] = labelsPayload

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"containerType": containerType,
		},
	}

	resp, err := client.CreateNodeType(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateNodeTypeResult)
	nodeType := result.NodeType
	// Successfully created resource, now set id
	d.SetId(int64ToString(nodeType.ID))

	resourceNodeTypeRead(ctx, d, meta)
	return diags
}

func resourceNodeTypeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindNodeTypeByName(name)
	} else if id != "" {
		resp, err = client.GetNodeType(toInt64(id), &morpheus.Request{})
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
	log.Printf("API RESPONSE: %s", resp)

	// store resource data
	var nodeType NodeTypePayload
	json.Unmarshal(resp.Body, &nodeType)

	log.Printf("RESPONSE_PAYLOAD: %v", nodeType)
	d.SetId(int64ToString(nodeType.NodeType.ID))
	d.Set("name", nodeType.NodeType.Name)
	d.Set("short_name", nodeType.NodeType.ShortName)
	d.Set("labels", nodeType.Labels)
	d.Set("version", nodeType.NodeType.ContainerVersion)
	d.Set("technology", nodeType.NodeType.ProvisionType.Code)
	d.Set("virtual_image_id", nodeType.NodeType.VirtualImage.ID)
	d.Set("service_port", parseServicePortPayload(nodeType.NodeType.ContainerPorts))

	// script templates
	var scriptTemplates []int64
	if nodeType.NodeType.ContainerScripts != nil {
		// iterate over the array of script templates
		for i := 0; i < len(nodeType.NodeType.ContainerScripts); i++ {
			scriptTemplate := nodeType.NodeType.ContainerScripts[i]
			scriptTemplates = append(scriptTemplates, scriptTemplate.ID)
		}
	}

	// file templates
	var fileTemplates []int64
	if nodeType.NodeType.ContainerTemplates != nil {
		// iterate over the array of file templates
		for i := 0; i < len(nodeType.NodeType.ContainerTemplates); i++ {
			fileTemplate := nodeType.NodeType.ContainerTemplates[i]
			fileTemplates = append(fileTemplates, fileTemplate.ID)
		}
	}

	stateScriptTemplates := matchTemplatesWithSchema(scriptTemplates, d.Get("script_template_ids").([]interface{}))
	stateFileTemplates := matchTemplatesWithSchema(fileTemplates, d.Get("file_template_ids").([]interface{}))

	d.Set("script_template_ids", stateScriptTemplates)
	d.Set("file_template_ids", stateFileTemplates)
	if nodeType.NodeType.ProvisionType.Code == "vmware" {
		// Extra Options
		extraOptions := make(map[string]interface{})
		if nodeType.NodeType.Config.ExtraOptions != nil {
			log.Printf("FoundExtraOptions: %s", nodeType.NodeType.Config.ExtraOptions)
			for k, v := range nodeType.NodeType.Config.ExtraOptions {
				extraOptions[k] = v
			}
			d.Set("extra_options", extraOptions)
		}
	}
	d.Set("category", nodeType.NodeType.Category)
	return diags
}

func resourceNodeTypeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	name := d.Get("name").(string)

	// extra options
	config := make(map[string]interface{})
	if d.Get("extra_options") != nil {
		config["extraOptions"] = d.Get("extra_options")
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"containerType": map[string]interface{}{
				"name":              name,
				"shortName":         d.Get("short_name").(string),
				"containerVersion":  d.Get("version").(string),
				"provisionTypeCode": d.Get("technology").(string),
				//"environmentVariables": parseNodeTypeEnvironmentVariables(d.Get("evar").([]interface{})),
				"virtualImageId":     d.Get("virtual_image_id").(int),
				"config":             config,
				"containerPorts":     parseNodeTypeServicePorts(d.Get("service_port").([]interface{})),
				"containerScripts":   d.Get("script_template_ids"),
				"containerTemplates": d.Get("file_template_ids"),
				"category":           d.Get("category").(string),
				"serverType":         "vm",
			},
		},
	}

	resp, err := client.UpdateNodeType(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateNodeTypeResult)
	nodeType := result.NodeType
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(nodeType.ID))
	return resourceNodeTypeRead(ctx, d, meta)
}

func resourceNodeTypeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteNodeType(toInt64(id), req)
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

func parseNodeTypeServicePorts(variables []interface{}) []map[string]interface{} {
	var svcports []map[string]interface{}
	// iterate over the array of svcports
	for i := 0; i < len(variables); i++ {
		row := make(map[string]interface{})
		svcportconfig := variables[i].(map[string]interface{})
		for k, v := range svcportconfig {
			switch k {
			case "name":
				row["name"] = v.(string)
			case "port":
				row["port"] = v.(string)
			case "protocol":
				row["loadBalanceProtocol"] = v.(string)
			}
		}
		svcports = append(svcports, row)
	}
	return svcports
}

/* AWAITING API SUPPORT TO AVOID DUPLICATE ENTRIES
func parseNodeTypeEnvironmentVariables(variables []interface{}) []map[string]interface{} {
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
			case "port":
				row["value"] = v.(string)
			case "export":
				row["export"] = v.(bool)
			case "masked":
				row["masked"] = v
			}
		}
		evars = append(evars, row)
	}
	return evars
}
*/

// This cannot currently be handled efficiently by a DiffSuppressFunc.
// See: https://github.com/hashicorp/terraform-plugin-sdk/issues/477
func matchTemplatesWithSchema(templates []int64, declaredTemplates []interface{}) []int64 {
	result := make([]int64, len(declaredTemplates))

	rMap := make(map[int64]int64, len(templates))
	for _, template := range templates {
		rMap[template] = template
	}

	for i, definedTemplate := range declaredTemplates {
		definedTemplate := int64(definedTemplate.(int))

		if v, ok := rMap[definedTemplate]; ok {
			// matched node type declared by ID
			result[i] = v
			delete(rMap, v)
		}
	}
	// append unmatched node type to the result
	for _, rcpt := range rMap {
		result = append(result, rcpt)
	}
	return result
}

func parseServicePortPayload(variables []morpheus.ContainerPort) []map[string]interface{} {
	var svcports []map[string]interface{}
	// iterate over the array of svcports
	for i := 0; i < len(variables); i++ {
		row := make(map[string]interface{})
		row["name"] = variables[i].Name
		row["port"] = strconv.Itoa(int(variables[i].Port))
		row["protocol"] = variables[i].LoadBalanceProtocol
		svcports = append(svcports, row)
	}
	return svcports
}

type NodeTypePayload struct {
	morpheus.NodeType `json:"containerType"`
	/*NodeType struct {
		ID      int64 `json:"id"`
		Account struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"account"`
		Name             string `json:"name"`
		ShortName        string `json:"shortName"`
		Code             string `json:"code"`
		ContainerVersion string `json:"containerVersion"`
		ProvisionType    struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"provisionType"`
		VirtualImage struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"virtualImage"`
		Category string `json:"category"`
		Config   struct {
			ExtraOptions map[string]interface{} `json:"extraOptions"`
		} `json:"config"`
		ContainerPorts   []ContainerPort `json:"containerPorts"`
		ContainerScripts []struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"containerScripts"`
		ContainerTemplates []struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"containerTemplates"`
		EnvironmentVariables []struct {
			Evarname     string `json:"evarName"`
			Name         string `json:"name"`
			Defaultvalue string `json:"defaultValue"`
			Valuetype    string `json:"valueType"`
			Export       bool   `json:"export"`
			Masked       bool   `json:"masked"`
		} `json:"environmentVariables"`
	} `json:"containerType"`

	*/
}

/*
type ContainerPort struct {
	ID                  int64  `json:"id"`
	Name                string `json:"name"`
	Port                int64  `json:"port"`
	Loadbalanceprotocol string `json:"loadBalanceProtocol"`
	Exportname          string `json:"exportName"`
}
*/
