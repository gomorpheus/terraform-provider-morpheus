package morpheus

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceServiceNowIntegration() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a ServiceNow integration resource",
		CreateContext: resourceServiceNowIntegrationCreate,
		ReadContext:   resourceServiceNowIntegrationRead,
		UpdateContext: resourceServiceNowIntegrationUpdate,
		DeleteContext: resourceServiceNowIntegrationDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The id of the ServiceNow integration",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the ServiceNow integration",
				Required:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the SerivceNow integration is enabled",
				Optional:    true,
				Computed:    true,
			},
			"url": {
				Type:        schema.TypeString,
				Description: "The url of the ServiceNow instance",
				Required:    true,
			},
			"credential_id": {
				Description:   "The id of the credential store entry used for authentication",
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"username", "password"},
			},
			"username": {
				Type:        schema.TypeString,
				Description: "The username of the account used to connect to ServiceNow",
				Optional:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "The password of the account used to connect to ServiceNow",
				Optional:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256_hash := hex.EncodeToString(h.Sum(nil))
					return strings.EqualFold(old, sha256_hash)
				},
				DiffSuppressOnRefresh: true,
			},
			"cmdb_custom_mapping": {
				Type:        schema.TypeString,
				Description: "A JSON encoded payload to populate a specific field in the ServiceNow table and with a specific mapping",
				Optional:    true,
			},
			"cmdb_class_mapping": {
				Type:        schema.TypeMap,
				Description: "The mapping between Morpheus server types and ServiceNow CI classes",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"default_cmdb_business_class": {
				Type:        schema.TypeString,
				Description: "The default ServiceNow table that records are written to if they aren't explicitly defined",
				Optional:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceServiceNowIntegrationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	integration := make(map[string]interface{})

	integration["type"] = "serviceNow"
	integration["name"] = d.Get("name").(string)
	integration["enabled"] = d.Get("enabled").(bool)
	integration["serviceUrl"] = d.Get("url").(string)
	config := make(map[string]interface{})

	if d.Get("credential_id").(int) != 0 {
		credential := make(map[string]interface{})
		credential["type"] = "username-password"
		credential["id"] = d.Get("credential_id").(int)
		integration["credential"] = credential
	} else {
		integration["serviceUsername"] = d.Get("username").(string)
		integration["servicePassword"] = d.Get("password").(string)
	}

	if d.Get("cmdb_class_mapping") != nil {
		classMappingResponse, err := client.GetOptionSource("serviceNowServerMappings", &morpheus.Request{})
		if err != nil {
			diag.FromErr(err)
		}
		classMappingResult := classMappingResponse.Result.(*morpheus.GetOptionSourceResult)
		classMappingsInput := d.Get("cmdb_class_mapping").(map[string]interface{})
		var classMappings []Mapping
		for key, value := range classMappingsInput {
			matchStatus := false
			for _, mapping := range *classMappingResult.Data {
				if key == mapping.Name {
					var classMapping Mapping
					classMapping.Name = mapping.Name
					classMapping.ID = strconv.Itoa(int(mapping.Value.(float64)))
					classMapping.NowClass = value.(string)
					classMappings = append(classMappings, classMapping)
					matchStatus = true
				}
			}
			if !matchStatus {
				return diag.Errorf("The %s cmdb mapping class is not a supported class", key)
			}
		}
		config["serviceNowCmdbClassMapping"] = classMappings
	}
	config["serviceNowCMDBBusinessObject"] = d.Get("default_cmdb_business_class").(string)
	config["serviceNowCustomCmdbMapping"] = d.Get("cmdb_custom_mapping")

	integration["config"] = config

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"integration": integration,
		},
	}

	resp, err := client.CreateIntegration(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateIntegrationResult)
	integrationResult := result.Integration
	// Successfully created resource, now set id
	d.SetId(int64ToString(integrationResult.ID))

	resourceServiceNowIntegrationRead(ctx, d, meta)
	return diags
}

func resourceServiceNowIntegrationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindIntegrationByName(name)
	} else if id != "" {
		resp, err = client.GetIntegration(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Integration cannot be read without name or id")
	}

	if err != nil {
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
	result := resp.Result.(*morpheus.GetIntegrationResult)
	integration := result.Integration
	d.SetId(int64ToString(integration.ID))
	d.Set("name", integration.Name)
	d.Set("enabled", integration.Enabled)
	d.Set("url", integration.URL)
	if integration.Credential.ID == 0 {
		d.Set("username", integration.Username)
		d.Set("password", integration.PasswordHash)
	} else {
		d.Set("credential_id", integration.Credential.ID)
	}
	d.Set("cmdb_custom_mapping", integration.Config.ServiceNowCustomCmdbMapping)
	classMappings := make(map[string]interface{})
	// iterate over the array of classMappings
	for i := 0; i < len(integration.Config.ServiceNowCmdbClassMapping); i++ {
		classMap := integration.Config.ServiceNowCmdbClassMapping[i]
		classMapName := classMap.Name
		classMappings[classMapName] = classMap.NowClass
	}
	d.Set("cmdb_class_mapping", classMappings)
	d.Set("default_cmdb_business_class", integration.Config.ServiceNowCMDBBusinessObject)

	return diags
}

func resourceServiceNowIntegrationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	integration := make(map[string]interface{})

	integration["name"] = d.Get("name").(string)
	integration["enabled"] = d.Get("enabled").(bool)
	integration["type"] = "serviceNow"
	integration["serviceUrl"] = d.Get("url").(string)

	config := make(map[string]interface{})

	if d.Get("credential_id").(int) != 0 {
		credential := make(map[string]interface{})
		credential["type"] = "username-password"
		credential["id"] = d.Get("credential_id").(int)
		integration["credential"] = credential
	} else {
		if d.HasChange("username") {
			integration["serviceUsername"] = d.Get("username").(string)
		}
		if d.HasChange("password") {
			integration["servicePassword"] = d.Get("password").(string)
		}
	}

	if d.Get("cmdb_class_mapping") != nil {
		// Query the API to fetch the ID of the class map
		classMappingResponse, err := client.GetOptionSource("serviceNowServerMappings", &morpheus.Request{})
		if err != nil {
			diag.FromErr(err)
		}
		classMappingResult := classMappingResponse.Result.(*morpheus.GetOptionSourceResult)
		classMappingsInput := d.Get("cmdb_class_mapping").(map[string]interface{})
		var classMappings []Mapping
		for key, value := range classMappingsInput {
			matchStatus := false
			for _, mapping := range *classMappingResult.Data {
				if key == mapping.Name {
					var classMapping Mapping
					classMapping.Name = mapping.Name
					classMapping.ID = strconv.Itoa(int(mapping.Value.(float64)))
					classMapping.NowClass = value.(string)
					classMappings = append(classMappings, classMapping)
					matchStatus = true
				}
			}
			if !matchStatus {
				return diag.Errorf("The %s cmdb mapping class is not a supported class", key)
			}
		}
		config["serviceNowCmdbClassMapping"] = classMappings
	}
	config["serviceNowCMDBBusinessObject"] = d.Get("default_cmdb_business_class").(string)
	if d.HasChange("cmdb_custom_mapping") {
		config["serviceNowCustomCmdbMapping"] = d.Get("cmdb_custom_mapping")
	}
	integration["config"] = config

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"integration": integration,
		},
	}

	resp, err := client.UpdateIntegration(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateIntegrationResult)
	integrationResult := result.Integration

	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(integrationResult.ID))
	return resourceServiceNowIntegrationRead(ctx, d, meta)
}

func resourceServiceNowIntegrationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteIntegration(toInt64(id), req)
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

type Mapping struct {
	ID       string `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	NowClass string `json:"nowClass"`
}
