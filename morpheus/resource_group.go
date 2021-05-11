package morpheus

// this is for Groups/Sites.
// this resource has an extra Morpheus prefix in it
// to distinguish it from ResourceGroups.

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMorpheusGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus group resource.",
		CreateContext: resourceMorpheusGroupCreate,
		ReadContext:   resourceMorpheusGroupRead,
		UpdateContext: resourceMorpheusGroupUpdate,
		DeleteContext: resourceMorpheusGroupDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the group",
				Type:        schema.TypeString,
				Computed:    true,
			},
			// Required inputs
			"name": {
				Description: "A unique name scoped to your account for the group",
				Type:        schema.TypeString,
				Required:    true,
			},
			"code": {
				Description: "Optional code for use with policies",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"location": {
				Description: "Optional location argument for your group",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"cloud_ids": {
				Description: "An array of all the clouds assigned to this group",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func resourceMorpheusGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	code := d.Get("code").(string)
	location := d.Get("location").(string)
	// clouds := d.Get("clouds").([]interface{})

	// clouds is an array of string names, lookup each one via api.
	// then the api expects it an array of objects, but only looks for id right now
	// once api is better this should get simpler
	doUpdateClouds := false
	var clouds []map[string]interface{}
	cloudIds := d.Get("cloud_ids").(*schema.Set).List()
	if len(cloudIds) > 0 {
		doUpdateClouds = true
		for _, v := range cloudIds {
			cloudPayload := map[string]interface{}{
				"id": v,
			}
			clouds = append(clouds, cloudPayload)
		}
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"group": map[string]interface{}{
				"name":     name,
				"code":     code,
				"location": location,
			},
		},
	}
	jsonRequest, _ := json.Marshal(req.Body)
	log.Printf("API JSON REQUEST: %s", string(jsonRequest))
	log.Printf("API REQUEST: %s", req) // debug
	resp, err := client.CreateGroup(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.CreateGroupResult)
	group := result.Group

	// oh ya..update zones too.. should use Partial thingy
	// or, even better the api should do this all in 1 request
	// doUpdateClouds = false
	if doUpdateClouds {
		req2 := &morpheus.Request{
			Body: map[string]interface{}{
				"group": map[string]interface{}{
					"zones": clouds,
				},
			},
		}
		jsonRequest, _ := json.Marshal(req2.Body)
		log.Printf("API JSON REQUEST: %s", string(jsonRequest))
		log.Printf("API REQUEST: %s", req2) // debug
		resp2, err2 := client.UpdateGroupClouds(group.ID, req2)
		if err2 != nil {
			log.Printf("API FAILURE: %s - %s", resp2, err2)
			return diag.FromErr(err2)
		}
		log.Printf("API RESPONSE: %s", resp2)
	}

	// Successfully created resource, now set id
	d.SetId(int64ToString(group.ID))
	resourceMorpheusGroupRead(ctx, d, meta)
	return diags
}

func resourceMorpheusGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindGroupByName(name)
	} else if id != "" {
		resp, err = client.GetGroup(toInt64(id), &morpheus.Request{})
		// todo: ignore 404 errors...
	} else {
		return diag.Errorf("Group cannot be read without name or id")
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
	result := resp.Result.(*morpheus.GetGroupResult)
	group := result.Group
	if group != nil {
		d.SetId(int64ToString(group.ID))
		d.Set("name", group.Name)
		d.Set("code", group.Code)
		d.Set("location", group.Location)
		var clouds []int64
		if len(group.Clouds) > 0 {
			for _, v := range group.Clouds {
				clouds = append(clouds, v.ID)
			}
		}
		d.Set("cloud_ids", clouds)
	} else {
		return diag.Errorf("Group not found in response data.") // should not happen
	}

	return diags
}

func resourceMorpheusGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	name := d.Get("name").(string)
	code := d.Get("code").(string)
	location := d.Get("location").(string)
	// clouds := d.Get("clouds").([]interface{})

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"group": map[string]interface{}{
				"name":     name,
				"code":     code,
				"location": location,
			},
		},
	}
	jsonRequest, _ := json.Marshal(req.Body)
	log.Printf("API JSON REQUEST: %s", string(jsonRequest))
	log.Printf("API REQUEST: %s", req) // debug
	resp, err := client.UpdateGroup(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateGroupResult)
	group := result.Group

	// clouds is an array of string names, lookup each one via api.
	// then the api expects it an array of objects, but only looks for id right now
	// once api is better this should get simpler
	doUpdateClouds := false
	var clouds []map[string]interface{}
	cloudIds := d.Get("cloud_ids").(*schema.Set).List()
	if len(cloudIds) > 0 {
		doUpdateClouds = true
		for _, v := range cloudIds {
			cloudPayload := map[string]interface{}{
				"id": v,
			}
			clouds = append(clouds, cloudPayload)
		}
	}
	// oh ya..update zones too.. should use Partial thingy
	// or, even better the api should do this all in 1 request
	// doUpdateClouds = false
	if doUpdateClouds {
		req2 := &morpheus.Request{
			Body: map[string]interface{}{
				"group": map[string]interface{}{
					"zones": clouds,
				},
			},
		}
		jsonRequest, _ := json.Marshal(req2.Body)
		log.Printf("API JSON REQUEST: %s", string(jsonRequest))
		log.Printf("API REQUEST: %s", req2) // debug
		resp2, err2 := client.UpdateGroupClouds(group.ID, req2)
		if err2 != nil {
			log.Printf("API FAILURE: %s - %s", resp2, err2)
			return diag.FromErr(err2)
		}
		log.Printf("API RESPONSE: %s", resp2)
	}

	// Successfully updated resource, now set id
	d.SetId(int64ToString(group.ID))
	return resourceMorpheusGroupRead(ctx, d, meta)
}

func resourceMorpheusGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteGroup(toInt64(id), req)
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
	// result := resp.Result.(*morpheus.DeleteGroupResult)
	d.SetId("")
	return diags
}
