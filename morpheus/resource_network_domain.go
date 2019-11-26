package morpheus

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/gomorpheus/morpheusapi"
	"log"
	"fmt"
	"errors"
)

func resourceNetworkDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkDomainCreate,
		Read:   resourceNetworkDomainRead,
		Update: resourceNetworkDomainUpdate,
		Delete: resourceNetworkDomainDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_zone": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default: false,
			},
			"domain_controller": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default: false,
			},
			"active": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default: true,
			},
			"visibility": &schema.Schema{
				Type:     schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"private", "public", ""}, false),
				Default:      "private",
			},
			"tenant": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}


func resourceNetworkDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*morpheusapi.Client)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	// publicZone := d.Get("public_zone").(bool) // .(bool)
	// domainController := d.Get("domain_controller").(bool) // .(bool)
	//active := d.Get("active").(bool)

	req := &morpheusapi.Request{
		Body: map[string]interface{}{
			"networkDomain": map[string]interface{}{
				"name": name,
				"description": description,
				// "publicZone": publicZone,
				// "domainController": domainController,
				// "active":active,
			},
		},
	}
	resp, err := client.CreateNetworkDomain(req)
	if err != nil {
		log.Printf("API FAILURE:", resp, err)
		return err
	}
	log.Printf("API RESPONSE: ", resp)

	result := resp.Result.(*morpheusapi.CreateNetworkDomainResult)
	networkDomain := result.NetworkDomain
	// Successfully created resource, now set id
	d.SetId(int64ToString(networkDomain.ID))

	return resourceNetworkDomainRead(d, meta)
}

func resourceNetworkDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*morpheusapi.Client)
	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheusapi.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindNetworkDomainByName(name)
	} else if id != "" {
		resp, err = client.GetNetworkDomain(toInt64(id), &morpheusapi.Request{})
	} else {
		return errors.New("NetworkDomain cannot be read without name or id")
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
	result := resp.Result.(*morpheusapi.GetNetworkDomainResult)
	networkDomain := result.NetworkDomain
	if networkDomain != nil {
		d.SetId(int64ToString(networkDomain.ID))
		d.Set("name", networkDomain.Name)
		d.Set("description", networkDomain.Description)
		d.Set("active", networkDomain.Active)
		d.Set("public_zone", networkDomain.PublicZone)
		d.Set("domain_controller", networkDomain.DomainController)
		d.Set("visibility", networkDomain.Visibility)
		// d.Set("fqdn", networkDomain.Fqdn)
		// todo: more fields
	} else {
		return fmt.Errorf("NetworkDomain not found in response data.")	 // should not happen
	}
	
	return nil
}

func resourceNetworkDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*morpheusapi.Client)
	id := d.Id()
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	// publicZone := d.Get("public_zone").(bool) // .(bool)
	// domainController := d.Get("domain_controller").(bool) // .(bool)
	//active := d.Get("active").(bool)

	req := &morpheusapi.Request{
		Body: map[string]interface{}{
			"networkDomain": map[string]interface{}{
				"name": name,
				"description": description,
				// "publicZone": publicZone,
				// "domainController": domainController,
				//"active":active,
			},
		},
	}
	resp, err := client.UpdateNetworkDomain(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE:", resp, err)
		return err
	}
	log.Printf("API RESPONSE: ", resp)
	result := resp.Result.(*morpheusapi.UpdateNetworkDomainResult)
	networkDomain := result.NetworkDomain
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(networkDomain.ID))
	return resourceNetworkDomainRead(d, meta)
}

func resourceNetworkDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*morpheusapi.Client)
	id := d.Id()
	req := &morpheusapi.Request{}
	resp, err := client.DeleteNetworkDomain(toInt64(id), req)
	//result := resp.Result.(*morpheusapi.DeleteNetworkDomainResult)
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
	//d.setId("") // implicit
	return nil
}
