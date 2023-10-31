package morpheus

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"strings"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceInstanceType() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus instance type resource",
		CreateContext: resourceInstanceTypeCreate,
		ReadContext:   resourceInstanceTypeRead,
		UpdateContext: resourceInstanceTypeUpdate,
		DeleteContext: resourceInstanceTypeDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the instance type",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the instance type",
				Required:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The instance type code",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the instance type",
				Optional:    true,
				Computed:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "The organization labels associated with the script template (Only supported on Morpheus 5.5.3 or higher)",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"category": {
				Type:         schema.TypeString,
				Description:  "The instance type category (web, sql, nosql, apps, network, messaging, cache, os, cloud, utility)",
				ValidateFunc: validation.StringInSlice([]string{"web", "sql", "nosql", "apps", "network", "messaging", "cache", "os", "cloud", "utility"}, false),
				Required:     true,
			},
			"image_name": {
				Type:        schema.TypeString,
				Description: "The file name of the instance type logo image",
				Optional:    true,
			},
			"image_path": {
				Type:        schema.TypeString,
				Description: "The file path of the instance type logo image including the file name",
				Optional:    true,
			},
			"environment_prefix": {
				Type:        schema.TypeString,
				Description: "The prefix used for instance environment variables",
				Optional:    true,
				Computed:    true,
			},
			"enable_settings": {
				Type:        schema.TypeBool,
				Description: "Whether to enable settings for the instance type",
				Optional:    true,
				Computed:    true,
			},
			"enable_scaling": {
				Type:        schema.TypeBool,
				Description: "Whether to enable scaling for the instance type",
				Optional:    true,
				Computed:    true,
			},
			"enable_deployments": {
				Type:        schema.TypeBool,
				Description: "Whether to enable deployments for the instance type",
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
								log.Println(sha256_hash)
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
				Description: "The IDs of the inputs to associate with the instance type",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == old
				},
				Computed: true,
			},
			"visibility": {
				Type:         schema.TypeString,
				Description:  "The visibility of the instance type (public or private)",
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"public", "private"}, false),
			},
			"featured": {
				Type:        schema.TypeBool,
				Description: "Whether the instance type is marked as featured",
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceInstanceTypeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		for _, s := range attr.(*schema.Set).List() {
			labelsPayload = append(labelsPayload, s.(string))
		}
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"instanceType": map[string]interface{}{
				"name":                 name,
				"code":                 d.Get("code").(string),
				"description":          d.Get("description").(string),
				"labels":               labelsPayload,
				"category":             d.Get("category").(string),
				"visibility":           d.Get("visibility").(string),
				"optionTypes":          d.Get("option_type_ids"),
				"environmentPrefix":    d.Get("environment_prefix").(string),
				"environmentVariables": parseInstanceTypeEnvironmentVariables(d.Get("evar").([]interface{}), d),
				"hasSettings":          d.Get("enable_settings").(bool),
				"hasAutoScale":         d.Get("enable_scaling").(bool),
				"hasDeployment":        d.Get("enable_deployments").(bool),
				"featured":             d.Get("featured").(bool),
			},
		},
	}

	resp, err := client.CreateInstanceType(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateInstanceTypeResult)
	instanceType := result.InstanceType

	if d.Get("image_path") != "" && d.Get("image_name") != "" {
		data, err := os.ReadFile(d.Get("image_path").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		var filePayloads []*morpheus.FilePayload
		filePayload := &morpheus.FilePayload{
			ParameterName: "logo",
			FileName:      d.Get("image_name").(string),
			FileContent:   data,
		}
		filePayloads = append(filePayloads, filePayload)
		response, err := client.UpdateInstanceTypeLogo(instanceType.ID, filePayloads, &morpheus.Request{})
		if err != nil {
			log.Printf("API FAILURE: %s - %s", response, err)
			return diag.FromErr(err)
		}
		log.Printf("API RESPONSE: %s", response)
	}

	// Successfully created resource, now set id
	d.SetId(int64ToString(instanceType.ID))

	resourceInstanceTypeRead(ctx, d, meta)
	return diags
}

func resourceInstanceTypeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindInstanceTypeByName(name)
	} else if id != "" {
		resp, err = client.GetInstanceType(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Instance type cannot be read without name or id")
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
	var instanceTypePayload InstanceTypePayload
	json.Unmarshal(resp.Body, &instanceTypePayload)

	d.SetId(int64ToString(int64(instanceTypePayload.InstanceType.ID)))
	d.Set("name", instanceTypePayload.InstanceType.Name)
	d.Set("code", instanceTypePayload.InstanceType.Code)
	d.Set("description", instanceTypePayload.InstanceType.Description)
	d.Set("labels", instanceTypePayload.Labels)
	d.Set("category", instanceTypePayload.InstanceType.Category)
	d.Set("visibility", instanceTypePayload.InstanceType.Visibility)
	d.Set("environment_prefix", instanceTypePayload.InstanceType.EnvironmentPrefix)
	d.Set("enable_settings", instanceTypePayload.InstanceType.HasSettings)
	d.Set("enable_scaling", instanceTypePayload.InstanceType.HasAutoscale)
	d.Set("enable_deployments", instanceTypePayload.InstanceType.HasDeployment)
	d.Set("featured", instanceTypePayload.InstanceType.Featured)
	// inputs
	var inputs []int64
	if instanceTypePayload.InstanceType.OptionTypes != nil {
		// iterate over the array of option types
		for i := 0; i < len(instanceTypePayload.InstanceType.OptionTypes); i++ {
			input := instanceTypePayload.InstanceType.OptionTypes[i]
			inputs = append(inputs, int64(input.ID))
		}
	}

	stateInputs := matchTemplatesWithSchema(inputs, d.Get("option_type_ids").([]interface{}))
	d.Set("option_type_ids", stateInputs)

	var evars []map[string]interface{}
	if instanceTypePayload.InstanceType.EnvironmentVariables != nil {
		// iterate over the array of environment variables
		for i := 0; i < len(instanceTypePayload.InstanceType.EnvironmentVariables); i++ {
			environmentVariable := instanceTypePayload.InstanceType.EnvironmentVariables[i]
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
	return diags
}

func resourceInstanceTypeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	name := d.Get("name").(string)

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		for _, s := range attr.(*schema.Set).List() {
			labelsPayload = append(labelsPayload, s.(string))
		}
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"instanceType": map[string]interface{}{
				"name":                 name,
				"code":                 d.Get("code").(string),
				"description":          d.Get("description").(string),
				"labels":               labelsPayload,
				"category":             d.Get("category").(string),
				"visibility":           d.Get("visibility").(string),
				"optionTypes":          d.Get("option_type_ids"),
				"environmentPrefix":    d.Get("environment_prefix").(string),
				"environmentVariables": parseInstanceTypeEnvironmentVariables(d.Get("evar").([]interface{}), d),
				"hasSettings":          d.Get("enable_settings").(bool),
				"hasAutoScale":         d.Get("enable_scaling").(bool),
				"hasDeployment":        d.Get("enable_deployments").(bool),
				"featured":             d.Get("featured").(bool),
			},
		},
	}

	resp, err := client.UpdateInstanceType(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	//log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateInstanceTypeResult)
	instanceType := result.InstanceType

	if d.HasChange("image_name") || d.HasChange("image_path") {
		data, err := os.ReadFile(d.Get("image_path").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		var filePayloads []*morpheus.FilePayload
		filePayload := &morpheus.FilePayload{
			ParameterName: "logo",
			FileName:      d.Get("image_name").(string),
			FileContent:   data,
		}
		filePayloads = append(filePayloads, filePayload)
		client.UpdateInstanceTypeLogo(instanceType.ID, filePayloads, &morpheus.Request{})
	}

	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(instanceType.ID))
	return resourceInstanceTypeRead(ctx, d, meta)
}

func resourceInstanceTypeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteInstanceType(toInt64(id), req)
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

func parseInstanceTypeEnvironmentVariables(variables []interface{}, d *schema.ResourceData) []map[string]interface{} {
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

type InstanceTypePayload struct {
	morpheus.InstanceType `json:"instanceType"`
}
