package morpheus

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/gomorpheus/morpheus-go/morpheusapi"
	"log"
	"fmt"
	"errors"
)

func resourceCloud() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudCreate,
		Read:   resourceCloudRead,
		Update: resourceCloudUpdate,
		Delete: resourceCloudDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"code": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"location": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default: "", //eh?
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"visibility": &schema.Schema{
				Type:     schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"private", "public", ""}, false),
				Default:      "private",
			},
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default: true,
			},


			"config": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		

   // 			"api_url": &schema.Schema{
			// 	Type:     schema.TypeString,
			// 	Required: false,
			// },
			// "username": &schema.Schema{
			// 	Type:     schema.TypeString,
			// 	Required: false,
			// },
			// "password": &schema.Schema{
			// 	Type:     schema.TypeString,
			// 	Required: false,
			// },
			// "datacenter": &schema.Schema{
			// 	Type:     schema.TypeString,
			// 	Required: false,
			// },
			// "cluster": &schema.Schema{
			// 	Type:     schema.TypeString,
			// 	Required: false,
			// },
			// "resource_pool": &schema.Schema{
			// 	Type:     schema.TypeString,
			// 	Required: false,
			// },
			// "rpc_mode": &schema.Schema{
			// 	Type:     schema.TypeString,
			// 	Required: false,
			// 	// Default: "guestexec",
			// },
			// "hide_host_selection": &schema.Schema{
			// 	Type:     schema.TypeBool,
			// 	Required: false,
			// 	// Default: false,
			// },
			// "enable_vnc": &schema.Schema{
			// 	Type:     schema.TypeBool,
			// 	Required: false,
			// 	// Default: false,
			// },
			// "import_existing": &schema.Schema{
			// 	Type:     schema.TypeBool,
			// 	Required: false,
			// 	// Default: false,
			// },

			// a TON more to add...

			"tenants": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"groups": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

		},
	}
}


func resourceCloudCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*morpheusapi.Client)
	name := d.Get("name").(string)
	code := d.Get("code").(string)
	location := d.Get("location").(string)
	visibility := d.Get("visibility").(string)

	// api expects zoneType.code, silly
	cloudTypeCode := d.Get("type").(string)

	// config is a big map of who knows what
	var config map[string]interface{}
	if d.Get("config") != nil {
		config = d.Get("config").(map[string]interface{})
	}

	payload := map[string]interface{}{
		"zone": map[string]interface{}{
			"name": name,
			"code": code,
			"location": location,
			"zoneType": map[string]interface{}{
				"code": cloudTypeCode,
			},
			"config": config,
			"visibility": visibility,
			// "groups": groups,
		},
	}

	req := &morpheusapi.Request{Body: payload}

	resp, err := client.CreateCloud(req)
	if err != nil {
		log.Printf("API FAILURE:", resp, err)
		return err
	}
	log.Printf("API RESPONSE: ", resp)
	result := resp.Result.(*morpheusapi.CreateCloudResult)
	cloud := result.Cloud
	// Successfully created resource, now set id
	d.SetId(int64ToString(cloud.ID))
	return resourceCloudRead(d, meta)
}

func resourceCloudRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*morpheusapi.Client)
	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheusapi.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindCloudByName(name)
	} else if id != "" {
		resp, err = client.GetCloud(toInt64(id), &morpheusapi.Request{})
		// todo: ignore 404 errors...
	} else {
		return errors.New("Cloud cannot be read without name or id")
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
	result := resp.Result.(*morpheusapi.GetCloudResult)
	cloud := result.Cloud
	if cloud == nil {
		return fmt.Errorf("Cloud not found in response data.") // should not happen
	}
	
	d.SetId(int64ToString(cloud.ID))
	d.Set("name", cloud.Name)
	d.Set("code", cloud.Code)
	d.Set("location", cloud.Location)
	d.Set("visibility", cloud.Visibility)
	d.Set("enabled", cloud.Enabled)
	// d.Set("groups", cloud.Groups)
	// todo: more fields

	return nil
}

func resourceCloudUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*morpheusapi.Client)
	id := d.Id()
	name := d.Get("name").(string)
	code := d.Get("code").(string)
	location := d.Get("location").(string)
	// clouds := d.Get("clouds").([]interface{})

	req := &morpheusapi.Request{
		Body: map[string]interface{}{
			"zone": map[string]interface{}{
				"name": name,
				"code": code,
				"location": location,
				// "clouds": clouds,
			},
		},
	}
	resp, err := client.UpdateCloud(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE:", resp, err)
		return err
	}
	log.Printf("API RESPONSE: ", resp)
	result := resp.Result.(*morpheusapi.UpdateCloudResult)
	cloud := result.Cloud
	// Successfully updated resource, now set id
	d.SetId(int64ToString(cloud.ID))
	return resourceCloudRead(d, meta)
}

func resourceCloudDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*morpheusapi.Client)
	id := d.Id()
	req := &morpheusapi.Request{}
	resp, err := client.DeleteCloud(toInt64(id), req)
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
	// result := resp.Result.(*morpheusapi.DeleteCloudResult)
	//d.setId("") // implicit
	return nil
}
