package morpheus

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceServicePlan() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a service plan resource",
		CreateContext: resourceServicePlanCreate,
		ReadContext:   resourceServicePlanRead,
		UpdateContext: resourceServicePlanUpdate,
		DeleteContext: resourceServicePlanDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the service plan",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the service plan",
				Required:    true,
			},
			"active": {
				Type:        schema.TypeBool,
				Description: "Whether the service plan is active or not",
				Optional:    true,
				Default:     true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code for the service plan",
				Required:    true,
			},
			"display_order": {
				Type:        schema.TypeInt,
				Description: "The display or sort order of the service plan",
				Optional:    true,
			},
			"provision_type": {
				Type:         schema.TypeString,
				Description:  "The provision type of the service plan",
				ValidateFunc: validation.StringInSlice([]string{"aks", "alibaba", "amazon", "arm", "azure", "azureSqlServerDatabase", "maas", "cloudFoundryApp", "cloudFoundryDocker", "cloudFoundryService", "cloudFormation", "digitalocean", "docker", "swarm", "eks", "esxi", "external", "fusion", "gke", "google", "helm", "oneview", "huawei", "hyperv", "hypervisor", "bluemix", "kubernetes", "kvm", "lxc", "manual", "nutanix", "opentelekom", "openstack", "oraclecloud", "oraclevm", "rds", "scvmm", "scvmm-hypervisor", "softlayer", "terraform", "ucs", "unmanaged", "upcloud", "vagrant", "vcd.vapp", "vcloudair", "vcd", "virtualbox", "vmware", "workflow", "xen"}, false),
				Required:     true,
			},
			"region_code": {
				Type:        schema.TypeString,
				Description: "The region code for the service plan",
				Optional:    true,
			},
			"max_storage": {
				Type:        schema.TypeInt,
				Description: "The maximum amount of storage in bytes",
				Optional:    true,
			},
			"storage_size_type": {
				Type:         schema.TypeString,
				Description:  "The unit of measure used for the service plan storage (gb, mb)",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"gb", "mb"}, false),
			},
			"max_memory": {
				Type:        schema.TypeInt,
				Description: "The maximum amount of memory in bytes",
				Optional:    true,
			},
			"memory_size_type": {
				Type:         schema.TypeString,
				Description:  "The unit of measure used for the service plan memory (gb, mb)",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"gb", "mb"}, false),
			},
			"custom_memory": {
				Type:        schema.TypeBool,
				Description: "Whether customizable memory is an option",
				Optional:    true,
			},
			"max_cores": {
				Type:        schema.TypeInt,
				Description: "The maximum amount of processor cores",
				Optional:    true,
			},
			"custom_cores": {
				Type:        schema.TypeBool,
				Description: "Whether the option to customize the number of processor cores is avaiable",
				Optional:    true,
			},
			"cores_per_socket": {
				Type:        schema.TypeInt,
				Description: "The number of cores per socket",
				Optional:    true,
			},
			"customize_root_volume": {
				Type:        schema.TypeBool,
				Description: "Whether the root volume is customized",
				Optional:    true,
			},
			"customize_extra_volumes": {
				Type:        schema.TypeBool,
				Description: "Whether the additional volumes are customized",
				Optional:    true,
			},
			"add_volumes": {
				Type:        schema.TypeBool,
				Description: "Whether additional volumes",
				Optional:    true,
			},
			"max_disks_allowed": {
				Type:        schema.TypeInt,
				Description: "The maximum number of disks that are allowed to be added",
				Optional:    true,
			},
			"custom_storage_range": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"minimum": {
							Type:     schema.TypeString,
							Required: true,
						},
						"maximum": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"custom_memory_range": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"minimum": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"maximum": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"custom_cores_range": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"minimum": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"maximum": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"price_set_ids": {
				Type:        schema.TypeList,
				Description: "The list of price set ids associated with the service plan",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceServicePlanCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	servicePlan := make(map[string]interface{})

	servicePlan["name"] = d.Get("name").(string)
	servicePlan["code"] = d.Get("code").(string)
	servicePlan["active"] = d.Get("active").(bool)
	servicePlan["regionCode"] = d.Get("region_code").(string)
	servicePlan["sortOrder"] = d.Get("display_order").(int)

	servicePlan["maxMemory"] = d.Get("max_memory").(int)
	if d.Get("custom_memory") != nil {
		servicePlan["customMaxMemory"] = d.Get("custom_memory").(bool)
	}
	servicePlan["maxStorage"] = d.Get("max_storage").(int)
	if d.Get("max_cores") != nil {
		servicePlan["maxCores"] = d.Get("max_cores").(int)
	}
	if d.Get("custom_cores") != nil {
		servicePlan["customCores"] = d.Get("custom_cores").(bool)
	}

	ranges := make(map[string]interface{})
	if v, ok := d.GetOk("custom_cores_range"); ok {
		core_ranges := v.([]interface{})[0].(map[string]interface{})
		ranges["minCores"] = core_ranges["minimum"].(string)
		ranges["maxCores"] = core_ranges["maximum"].(string)
	}

	if v, ok := d.GetOk("custom_memory_range"); ok {
		memory_ranges := v.([]interface{})[0].(map[string]interface{})
		ranges["minMemory"] = memory_ranges["minimum"].(int)
		ranges["maxMemory"] = memory_ranges["maximum"].(int)
	}

	if v, ok := d.GetOk("custom_storage_range"); ok {
		storage_ranges := v.([]interface{})[0].(map[string]interface{})
		ranges["minStorage"] = storage_ranges["minimum"].(string)
		ranges["maxStorage"] = storage_ranges["maximum"].(string)
	}

	config := make(map[string]interface{})
	config["ranges"] = ranges
	servicePlan["config"] = config
	if d.Get("memory_size_type") != nil {
		config["memorySizeType"] = d.Get("memory_size_type").(string)
	}
	if d.Get("storage_size_type") != nil {
		config["storageSizeType"] = d.Get("storage_size_type").(string)
	}

	outResult, err := FindProvisionTypeByCode(client, d.Get("provision_type").(string))
	if err != nil {
		log.Printf("Provision Type Lookup: %s", err)
		return diag.FromErr(err)
	}
	respresult := outResult.Result.(*morpheus.GetProvisionTypeResult)
	servicePlan["provisionType"] = map[string]interface{}{
		"id": respresult.ProvisionType.ID,
	}

	switch d.Get("provision_type") {
	// All configuration settings
	case "vmware", "nutanix", "vcd":
		servicePlan["customMaxStorage"] = d.Get("customize_root_volume").(bool)
		servicePlan["customMaxDataStorage"] = d.Get("customize_extra_volumes").(bool)
		servicePlan["addVolumes"] = d.Get("add_volumes").(bool)
		servicePlan["maxDisks"] = d.Get("max_disks_allowed").(int)
		servicePlan["coresPerSocket"] = d.Get("cores_per_socket").(int)
	// No cores per socket
	case "xen", "vcloudair", "scvmm", "oraclevm", "oraclecloud", "kvm", "hyperv", "fusion", "esxi":
		servicePlan["customMaxStorage"] = d.Get("customize_root_volume").(bool)
		servicePlan["customMaxDataStorage"] = d.Get("customize_extra_volumes").(bool)
		servicePlan["addVolumes"] = d.Get("add_volumes").(bool)
		servicePlan["maxDisks"] = d.Get("max_disks_allowed").(int)
	// Customize root volume
	case "docker", "swarm":
		servicePlan["customMaxStorage"] = d.Get("customize_root_volume").(bool)
	}

	var priceSetIds []map[string]interface{}
	if d.Get("price_set_ids") != nil {
		priceSetIdList := d.Get("price_set_ids").([]interface{})
		var intPayload []int
		for i := 0; i < len(priceSetIdList); i++ {
			intPayload = append(intPayload, priceSetIdList[i].(int))
		}
		sort.Ints(intPayload)
		// iterate over the array of price set ids
		for _, indata := range intPayload {
			row := make(map[string]interface{})
			row["id"] = indata
			priceSetIds = append(priceSetIds, row)
		}
	}
	servicePlan["priceSets"] = priceSetIds

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"servicePlan": servicePlan,
		},
	}
	resp, err := client.CreatePlan(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		log.Fatal(err)
	}

	servicePlanID := fmt.Sprintf("%v", result["id"])
	// Successfully created resource, now set id
	d.SetId(servicePlanID)
	resourceServicePlanRead(ctx, d, meta)
	return diags
}

func resourceServicePlanRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindPlanByName(name)
	} else if id != "" {
		resp, err = client.GetPlan(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Plan cannot be read without name or id")
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
	var servicePlan MorpheusPlan
	json.Unmarshal(resp.Body, &servicePlan)

	d.SetId(intToString(int(servicePlan.ServicePlan.ID)))
	d.Set("name", servicePlan.ServicePlan.Name)
	d.Set("active", servicePlan.ServicePlan.Active)
	d.Set("code", servicePlan.ServicePlan.Code)
	d.Set("display_order", servicePlan.ServicePlan.Sortorder)
	d.Set("provision_type", servicePlan.ServicePlan.Provisiontype.Code)
	d.Set("region_code", servicePlan.ServicePlan.Regioncode)
	d.Set("max_memory", servicePlan.ServicePlan.Maxmemory)
	d.Set("max_storage", servicePlan.ServicePlan.Maxstorage)
	d.Set("storage_size_type", servicePlan.ServicePlan.Config.Storagesizetype)
	d.Set("memory_size_type", servicePlan.ServicePlan.Config.Memorysizetype)
	d.Set("custom_memory", servicePlan.ServicePlan.Custommaxmemory)
	d.Set("max_cores", servicePlan.ServicePlan.Maxcores)
	d.Set("custom_cores", servicePlan.ServicePlan.Customcores)
	if _, ok := d.GetOk("cores_per_socket"); ok {
		d.Set("cores_per_socket", servicePlan.ServicePlan.Corespersocket)
	}
	if d.Get("customize_root_volume") != nil {
		d.Set("customize_root_volume", servicePlan.ServicePlan.Custommaxstorage)
	}
	if d.Get("customize_extra_volumes") != nil {
		d.Set("customize_extra_volumes", servicePlan.ServicePlan.Custommaxdatastorage)
	}
	if d.Get("add_volumes") != nil {
		d.Set("add_volumes", servicePlan.ServicePlan.Addvolumes)
	}
	if _, ok := d.GetOk("max_disks_allowed"); ok {
		d.Set("max_disks_allowed", servicePlan.ServicePlan.Maxdisks)
	}
	if _, ok := d.GetOk("custom_storage_range"); ok {
		d.Set("custom_storage_range", flattenRanges("storage", &servicePlan))
	}

	if _, ok := d.GetOk("custom_memory_range"); ok {
		d.Set("custom_memory_range", flattenRanges("memory", &servicePlan))
	}

	if _, ok := d.GetOk("custom_cores_range"); ok {
		d.Set("custom_cores_range", flattenRanges("cores", &servicePlan))
	}

	// Create an array from the price set ids
	var priceSetIds []int
	if len(servicePlan.ServicePlan.Pricesets) > 0 {
		for _, v := range servicePlan.ServicePlan.Pricesets {
			priceSetIds = append(priceSetIds, int(v.ID))
		}
	}

	// Adjust the price set order returned from the API (numerical low to high)
	// to match the order defined in the Terraform code.
	statePriceSetPayload := matchPriceSetsWithSchema(priceSetIds, d.Get("price_set_ids").([]interface{}))
	d.Set("price_set_ids", statePriceSetPayload)
	return diags
}

func resourceServicePlanUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	servicePlan := make(map[string]interface{})

	servicePlan["name"] = d.Get("name").(string)
	servicePlan["code"] = d.Get("code").(string)
	servicePlan["active"] = d.Get("active").(bool)
	servicePlan["regionCode"] = d.Get("region_code").(string)
	servicePlan["sortOrder"] = d.Get("display_order").(int)

	servicePlan["maxMemory"] = d.Get("max_memory").(int)
	if d.Get("custom_memory") != nil {
		servicePlan["customMaxMemory"] = d.Get("custom_memory").(bool)
	}
	servicePlan["maxStorage"] = d.Get("max_storage").(int)
	if d.Get("max_cores") != nil {
		servicePlan["maxCores"] = d.Get("max_cores").(int)
	}
	if d.Get("custom_cores") != nil {
		servicePlan["customCores"] = d.Get("custom_cores").(bool)
	}

	ranges := make(map[string]interface{})
	if v, ok := d.GetOk("custom_cores_range"); ok {
		core_ranges := v.([]interface{})[0].(map[string]interface{})
		ranges["minCores"] = core_ranges["minimum"].(string)
		ranges["maxCores"] = core_ranges["maximum"].(string)
	}

	if v, ok := d.GetOk("custom_memory_range"); ok {
		memory_ranges := v.([]interface{})[0].(map[string]interface{})
		ranges["minMemory"] = memory_ranges["minimum"].(int)
		ranges["maxMemory"] = memory_ranges["maximum"].(int)
	}

	if v, ok := d.GetOk("custom_storage_range"); ok {
		storage_ranges := v.([]interface{})[0].(map[string]interface{})
		ranges["minStorage"] = storage_ranges["minimum"].(string)
		ranges["maxStorage"] = storage_ranges["maximum"].(string)
	}

	config := make(map[string]interface{})
	config["ranges"] = ranges
	servicePlan["config"] = config
	if d.Get("memory_size_type") != nil {
		config["memorySizeType"] = d.Get("memory_size_type").(string)
	}
	if d.Get("storage_size_type") != nil {
		config["storageSizeType"] = d.Get("storage_size_type").(string)
	}

	outResult, err := FindProvisionTypeByCode(client, d.Get("provision_type").(string))
	if err != nil {
		log.Printf("Provision Type Lookup: %s", err)
		return diag.FromErr(err)
	}
	respresult := outResult.Result.(*morpheus.GetProvisionTypeResult)
	servicePlan["provisionType"] = map[string]interface{}{
		"id": respresult.ProvisionType.ID,
	}

	switch d.Get("provision_type") {
	// All configuration settings
	case "vmware", "nutanix", "vcd":
		servicePlan["customMaxStorage"] = d.Get("customize_root_volume").(bool)
		servicePlan["customMaxDataStorage"] = d.Get("customize_extra_volumes").(bool)
		servicePlan["addVolumes"] = d.Get("add_volumes").(bool)
		servicePlan["maxDisks"] = d.Get("max_disks_allowed").(int)
		servicePlan["coresPerSocket"] = d.Get("cores_per_socket").(int)
	// No cores per socket
	case "xen", "vcloudair", "scvmm", "oraclevm", "oraclecloud", "kvm", "hyperv", "fusion", "esxi":
		servicePlan["customMaxStorage"] = d.Get("customize_root_volume").(bool)
		servicePlan["customMaxDataStorage"] = d.Get("customize_extra_volumes").(bool)
		servicePlan["addVolumes"] = d.Get("add_volumes").(bool)
		servicePlan["maxDisks"] = d.Get("max_disks_allowed").(int)
	// Customize root volume
	case "docker", "swarm":
		servicePlan["customMaxStorage"] = d.Get("customize_root_volume").(bool)
	}

	var priceSetIds []map[string]interface{}
	if d.Get("price_set_ids") != nil {
		priceSetIdList := d.Get("price_set_ids").([]interface{})
		var intPayload []int
		for i := 0; i < len(priceSetIdList); i++ {
			intPayload = append(intPayload, priceSetIdList[i].(int))
		}
		sort.Ints(intPayload)
		// iterate over the array of price set ids
		for _, indata := range intPayload {
			row := make(map[string]interface{})
			row["id"] = indata
			priceSetIds = append(priceSetIds, row)
		}
	}
	servicePlan["priceSets"] = priceSetIds

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"servicePlan": servicePlan,
		},
	}

	resp, err := client.UpdatePlan(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		log.Fatal(err)
	}

	d.SetId(id)
	return resourceServicePlanRead(ctx, d, meta)
}

func resourceServicePlanDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeletePlan(toInt64(id), req)
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

// FindPriceByName gets an existing price by name
func FindProvisionTypeByCode(client *morpheus.Client, code string) (*morpheus.Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListProvisionTypes(&morpheus.Request{
		QueryParams: map[string]string{
			"code": code,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*morpheus.ListProvisionTypesResult)
	provisionTypeCount := len(*listResult.ProvisionTypes)
	if provisionTypeCount != 1 {
		return resp, fmt.Errorf("found %d Provision Types for %v", provisionTypeCount, code)
	}
	firstRecord := (*listResult.ProvisionTypes)[0]
	provisionTypeID := firstRecord.ID
	return client.GetProvisionType(int64(provisionTypeID), &morpheus.Request{})
}

func flattenRanges(rangeType string, plan *MorpheusPlan) []interface{} {
	result := make([]interface{}, 0, 1)
	data := make(map[string]interface{})
	switch rangeType {
	case "storage":
		data["minimum"] = plan.ServicePlan.Config.Ranges.Minstorage
		data["maximum"] = plan.ServicePlan.Config.Ranges.Maxstorage
	case "memory":
		data["minimum"] = plan.ServicePlan.Config.Ranges.Minmemory
		data["maximum"] = plan.ServicePlan.Config.Ranges.Maxmemory
	case "cores":
		data["minimum"] = plan.ServicePlan.Config.Ranges.Mincores
		data["maximum"] = plan.ServicePlan.Config.Ranges.Maxcores
	}

	result = append(result, data)
	return result
}

// This cannot currently be handled efficiently by a DiffSuppressFunc.
// See: https://github.com/hashicorp/terraform-plugin-sdk/issues/477
func matchPriceSetsWithSchema(priceSets []int, declaredPriceSets []interface{}) []int {
	result := make([]int, len(declaredPriceSets))

	rMap := make(map[int]int, len(priceSets))
	for _, priceSet := range priceSets {
		rMap[priceSet] = priceSet
	}

	for i, declaredPriceSet := range declaredPriceSets {
		declaredPriceSet := declaredPriceSet.(int)

		if v, ok := rMap[declaredPriceSet]; ok {
			// matched recipient declared by ID
			result[i] = v
			delete(rMap, v)
		}
	}
	// append unmatched price sets to the result
	for _, rcpt := range rMap {
		result = append(result, rcpt)
	}
	return result
}

type MorpheusPlan struct {
	ServicePlan struct {
		ID                   int         `json:"id"`
		Name                 string      `json:"name"`
		Code                 string      `json:"code"`
		Active               bool        `json:"active"`
		Sortorder            int         `json:"sortOrder"`
		Description          string      `json:"description"`
		Maxstorage           int64       `json:"maxStorage"`
		Maxmemory            int         `json:"maxMemory"`
		Maxcpu               interface{} `json:"maxCpu"`
		Maxcores             int         `json:"maxCores"`
		Maxdisks             int         `json:"maxDisks"`
		Corespersocket       int         `json:"coresPerSocket"`
		Customcpu            bool        `json:"customCpu"`
		Customcores          bool        `json:"customCores"`
		Custommaxstorage     bool        `json:"customMaxStorage"`
		Custommaxdatastorage bool        `json:"customMaxDataStorage"`
		Custommaxmemory      bool        `json:"customMaxMemory"`
		Addvolumes           bool        `json:"addVolumes"`
		Memoryoptionsource   interface{} `json:"memoryOptionSource"`
		Cpuoptionsource      interface{} `json:"cpuOptionSource"`
		Datecreated          time.Time   `json:"dateCreated"`
		Lastupdated          time.Time   `json:"lastUpdated"`
		Regioncode           string      `json:"regionCode"`
		Visibility           string      `json:"visibility"`
		Editable             bool        `json:"editable"`
		Provisiontype        struct {
			ID                        int    `json:"id"`
			Name                      string `json:"name"`
			Code                      string `json:"code"`
			Rootdiskcustomizable      bool   `json:"rootDiskCustomizable"`
			Addvolumes                bool   `json:"addVolumes"`
			Customizevolume           bool   `json:"customizeVolume"`
			Hasconfigurablecpusockets bool   `json:"hasConfigurableCpuSockets"`
		} `json:"provisionType"`
		Tenants   string              `json:"tenants"`
		Pricesets []morpheus.PriceSet `json:"priceSets"`
		Config    struct {
			Storagesizetype string `json:"storageSizeType"`
			Memorysizetype  string `json:"memorySizeType"`
			Ranges          struct {
				Minstorage string `json:"minStorage"`
				Maxstorage string `json:"maxStorage"`
				Minmemory  int    `json:"minMemory"`
				Maxmemory  int    `json:"maxMemory"`
				Mincores   string `json:"minCores"`
				Maxcores   string `json:"maxCores"`
			} `json:"ranges"`
		} `json:"config"`
		Zones []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"zones"`
		Permissions struct {
			Resourcepermissions struct {
				Defaultstore  bool `json:"defaultStore"`
				Allplans      bool `json:"allPlans"`
				Defaulttarget bool `json:"defaultTarget"`
				Canmanage     bool `json:"canManage"`
				All           bool `json:"all"`
				Account       struct {
					ID int `json:"id"`
				} `json:"account"`
				Sites []struct {
					ID      int    `json:"id"`
					Name    string `json:"name"`
					Default bool   `json:"default"`
				} `json:"sites"`
				Plans []interface{} `json:"plans"`
			} `json:"resourcePermissions"`
		} `json:"permissions"`
	} `json:"servicePlan"`
}
