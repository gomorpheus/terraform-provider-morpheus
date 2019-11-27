package morpheus

// this is for Groups/Sites. 
// this resource has an extra Morpheus prefix in it
// to distinguish it from ResourceGroups. 

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	// "github.com/hashicorp/terraform/helper/schema"
	//_"github.com/hashicorp/terraform/helper/validation"

	"github.com/gomorpheus/morpheus-go-sdk"
	"log"
	"fmt"
	"errors"
)

func resourceMorpheusGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceMorpheusGroupCreate,
		Read:   resourceMorpheusGroupRead,
		Update: resourceMorpheusGroupUpdate,
		Delete: resourceMorpheusGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
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
			},
			"clouds": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}


func resourceMorpheusGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*morpheusapi.Client)
	name := d.Get("name").(string)
	code := d.Get("code").(string)
	location := d.Get("location").(string)
	// clouds := d.Get("clouds").([]interface{})
	
	// clouds is an array of string names, lookup each one via api.
	// then the api expects it an array of objects, but only looks for id right now
	// once api is better this should get simpler
	doUpdateClouds := false
	var clouds []map[string]interface{}
	//clouds := make([]map[string]interface{}, 0, len(cloudNames))
	if d.Get("clouds") != nil {
		doUpdateClouds = true
		cloudNames := d.Get("clouds").([]interface{})
		//clouds = make([]map[string]interface{}, 0, len(cloudNames))
		for i := 0; i < len(cloudNames); i++ {
			findResponse, findErr := client.FindCloudByName(cloudNames[i].(string))
			if findErr != nil {
				return findErr
			}
			cloud := findResponse.Result.(*morpheusapi.GetCloudResult).Cloud
			cloudPayload := map[string]interface{}{
				"id": cloud.ID,
				"name": cloud.Name,
			}
			clouds = append(clouds, cloudPayload)
		}
	}

	req := &morpheusapi.Request{
		Body: map[string]interface{}{
			"group": map[string]interface{}{
				"name": name,
				"code": code,
				"location": location,
				// "zones": clouds,
			},
		},
	}
	resp, err := client.CreateGroup(req)
	if err != nil {
		log.Printf("API FAILURE:", resp, err)
		return err
	}
	log.Printf("API RESPONSE: ", resp)
	result := resp.Result.(*morpheusapi.CreateGroupResult)
	group := result.Group
	
	
	// oh ya..update zones too.. should use Partial thingy
	// or, even better the api should do this all in 1 request
	// doUpdateClouds = false
	if doUpdateClouds {
		req2 := &morpheusapi.Request{
			Body: map[string]interface{}{
				"group": map[string]interface{}{
					"zones": clouds,
				},
			},
		}
		resp2, err2 := client.UpdateGroupClouds(group.ID, req2)
		if err2 != nil {
			log.Printf("API FAILURE:", resp2, err2)
			return err
		}
		log.Printf("API RESPONSE: ", resp2)
	}

	// Successfully created resource, now set id
	d.SetId(int64ToString(group.ID))
	return resourceMorpheusGroupRead(d, meta)
}

func resourceMorpheusGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*morpheusapi.Client)
	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheusapi.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindGroupByName(name)
	} else if id != "" {
		resp, err = client.GetGroup(toInt64(id), &morpheusapi.Request{})
		// todo: ignore 404 errors...
	} else {
		return errors.New("Group cannot be read without name or id")
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
	result := resp.Result.(*morpheusapi.GetGroupResult)
	group := result.Group
	if group != nil {
		d.SetId(int64ToString(group.ID))
		d.Set("name", group.Name)
		d.Set("code", group.Code)
		d.Set("location", group.Location)
		// d.Set("clouds", group.Clouds)
		// todo: more fields
	} else {
		return fmt.Errorf("Group not found in response data.")	 // should not happen
	}
	
	return nil
}

func resourceMorpheusGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*morpheusapi.Client)
	id := d.Id()
	name := d.Get("name").(string)
	code := d.Get("code").(string)
	location := d.Get("location").(string)
	// clouds := d.Get("clouds").([]interface{})

	req := &morpheusapi.Request{
		Body: map[string]interface{}{
			"group": map[string]interface{}{
				"name": name,
				"code": code,
				"location": location,
				// "clouds": clouds,
			},
		},
	}
	resp, err := client.UpdateGroup(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE:", resp, err)
		return err
	}
	log.Printf("API RESPONSE: ", resp)
	result := resp.Result.(*morpheusapi.UpdateGroupResult)
	group := result.Group
	// Successfully updated resource, now set id
	d.SetId(int64ToString(group.ID))
	return resourceMorpheusGroupRead(d, meta)
}

func resourceMorpheusGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*morpheusapi.Client)
	id := d.Id()
	req := &morpheusapi.Request{}
	resp, err := client.DeleteGroup(toInt64(id), req)
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
	// result := resp.Result.(*morpheusapi.DeleteGroupResult)
	//d.setId("") // implicit
	return nil
}
