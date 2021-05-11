package morpheus

import (
	"context"
	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceVsphereCloud() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus cloud resource.",
		CreateContext: resourceVsphereCloudCreate,
		ReadContext:   resourceVsphereCloudRead,
		UpdateContext: resourceVsphereCloudUpdate,
		DeleteContext: resourceVsphereCloudDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the cloud",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "A unique name scoped to your account for the cloud",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The user friendly description of the cloud",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"code": {
				Description: "Optional code for use with policies",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"location": {
				Description: "Optional location for your cloud",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"visibility": {
				Description:  "Determines whether the resource is visible in sub-tenants or not",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"private", "public", ""}, false),
				Default:      "private",
			},
			"tenants": {
				Description: "",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"enabled": {
				Description: "Determines whether the cloud is active or not",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"api_url": {
				Type:        schema.TypeString,
				Description: "",
				Required:    true,
			},
			"username": {
				Type:        schema.TypeString,
				Description: "The username of the VMware vSphere account",
				Required:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "The password of the VMware vSphere account",
				Required:    true,
			},
			"datacenter": {
				Type:        schema.TypeString,
				Description: "",
				Required:    true,
			},
			"cluster": {
				Type:        schema.TypeString,
				Description: "",
				Required:    true,
			},
			"resource_pool": {
				Type:        schema.TypeString,
				Description: "",
				Optional:    true,
			},
			"rpc_mode": {
				Type:        schema.TypeString,
				Description: "",
				Optional:    true,
				Default:     "guestexec",
			},
			"hide_host_selection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"import_existing": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"enable_vnc": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"groups": {
				Description: "The group the cloud is assigned to",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceVsphereCloudCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	code := d.Get("code").(string)
	visibility := d.Get("visibility").(string)

	// config is a big map of who knows what
	config := make(map[string]interface{})
	config["certificateProvider"] = "internal"
	config["apiUrl"] = d.Get("api_url")
	config["username"] = d.Get("username")
	config["password"] = d.Get("password")
	config["datacenter"] = d.Get("datacenter")
	config["cluster"] = d.Get("cluster")
	config["rpcMode"] = d.Get("rpc_mode")
	config["hideHostSelection"] = d.Get("hide_host_selection")
	config["enableVnc"] = d.Get("enable_vnc")
	config["importExisting"] = d.Get("import_existing")

	payload := map[string]interface{}{
		"zone": map[string]interface{}{
			"name":        name,
			"code":        code,
			"location":    d.Get("location").(string),
			"description": d.Get("description").(string),
			"zoneType": map[string]interface{}{
				"code": "vmware",
			},
			"config":     config,
			"visibility": visibility,
		},
	}

	req := &morpheus.Request{Body: payload}

	resp, err := client.CreateCloud(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.CreateCloudResult)
	cloud := result.Cloud
	// Successfully created resource, now set id
	d.SetId(int64ToString(cloud.ID))
	resourceCloudRead(ctx, d, meta)
	return diags
}

func resourceVsphereCloudRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindCloudByName(name)
	} else if id != "" {
		resp, err = client.GetCloud(toInt64(id), &morpheus.Request{})
		// todo: ignore 404 errors...
	} else {
		return diag.Errorf("Cloud cannot be read without name or id")
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
	result := resp.Result.(*morpheus.GetCloudResult)
	cloud := result.Cloud
	if cloud == nil {
		return diag.Errorf("Cloud not found in response data.") // should not happen
	}

	d.SetId(int64ToString(cloud.ID))
	d.Set("name", cloud.Name)
	d.Set("code", cloud.Code)
	d.Set("location", cloud.Location)
	d.Set("visibility", cloud.Visibility)
	d.Set("enabled", cloud.Enabled)
	d.Set("api_url", cloud.Config["apiUrl"])

	return diags
}

func resourceVsphereCloudUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	name := d.Get("name").(string)
	code := d.Get("code").(string)
	location := d.Get("location").(string)
	// clouds := d.Get("clouds").([]interface{})

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"zone": map[string]interface{}{
				"name":     name,
				"code":     code,
				"location": location,
				// "clouds": clouds,
			},
		},
	}
	resp, err := client.UpdateCloud(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateCloudResult)
	cloud := result.Cloud
	// Successfully updated resource, now set id
	d.SetId(int64ToString(cloud.ID))
	return resourceCloudRead(ctx, d, meta)
}

func resourceVsphereCloudDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteCloud(toInt64(id), req)
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
