package morpheus

import (
	"context"
	"fmt"
	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEnvironment() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus environment resource",
		CreateContext: resourceEnvironmentCreate,
		ReadContext:   resourceEnvironmentRead,
		UpdateContext: resourceEnvironmentUpdate,
		DeleteContext: resourceEnvironmentDelete,

		Schema: map[string]*schema.Schema{
			"active": {
				Type:        schema.TypeBool,
				Description: "Whether the environment is enabled or not",
				Optional:    true,
				Default:     true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the environment",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the environment",
				Optional:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the environment",
				Optional:    true,
			},
			"visibility": {
				Type:        schema.TypeString,
				Description: "Whether the environment is visible in sub-tenants or not",
				Optional:    true,
				Default:     "private",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceEnvironmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	code := d.Get("code").(string)
	req := &morpheus.Request{
		Body: map[string]interface{}{
			"environment": map[string]interface{}{
				"active":      d.Get("active").(bool),
				"name":        name,
				"description": description,
				"code":        code,
				"visibility":  d.Get("visibility").(string),
			},
		},
	}

	resp, err := client.CreateEnvironment(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}

	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateEnvironmentResult)
	environment := result.Environment
	// Successfully created resource, now set id
	d.SetId(int64ToString(environment.ID))

	return diags
}

func resourceEnvironmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindEnvironmentByName(name)
	} else if id != "" {
		resp, err = client.GetEnvironment(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Environment cannot be read without name or id")
	}

	if err != nil {
		// 404 is ok?
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
			return nil
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)
			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)

	// store resource data
	result := resp.Result.(*morpheus.GetEnvironmentResult)
	environment := result.Environment
	if environment != nil {
		d.SetId(int64ToString(environment.ID))
		d.Set("active", environment.Active)
		d.Set("name", environment.Name)
		d.Set("description", environment.Description)
		d.Set("visibility", environment.Visibility)
		d.Set("code", environment.Code)
	} else {
		log.Println(environment)
		err := fmt.Errorf("read operation: environment not found in response data") // should not happen
		return diag.FromErr(err)
	}

	return diags
}

func resourceEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	id := d.Id()
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	code := d.Get("code").(string)

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"environment": map[string]interface{}{
				"active":      d.Get("active").(bool),
				"name":        name,
				"description": description,
				"code":        code,
				"visibility":  d.Get("visibility").(string),
			},
		},
	}
	resp, err := client.UpdateEnvironment(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateEnvironmentResult)
	account := result.Environment
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(account.ID))
	return resourceEnvironmentRead(ctx, d, meta)
}

func resourceEnvironmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteEnvironment(toInt64(id), req)
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
