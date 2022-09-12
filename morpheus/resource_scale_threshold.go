package morpheus

import (
	"context"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceScaleThreshold() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus scale threshold resource.",
		CreateContext: resourceScaleThresholdCreate,
		ReadContext:   resourceScaleThresholdRead,
		UpdateContext: resourceScaleThresholdUpdate,
		DeleteContext: resourceScaleThresholdDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the scale threshold",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the scale threshold",
				Required:    true,
			},
			"auto_upscale": {
				Type:        schema.TypeBool,
				Description: "Whether to scale up the number of instances",
				Required:    true,
			},
			"auto_downscale": {
				Type:        schema.TypeBool,
				Description: "Whether to scale down the number of instances",
				Required:    true,
			},
			"min_count": {
				Type:        schema.TypeInt,
				Description: "The minimum number of instances to scale down to",
				Required:    true,
			},
			"max_count": {
				Type:        schema.TypeInt,
				Description: "The maximum number of instances to scale up to",
				Required:    true,
			},
			"enable_cpu_threshold": {
				Type:        schema.TypeBool,
				Description: "Whether scaling operations based upon cpu usage is enabled or not",
				Optional:    true,
				Computed:    true,
			},
			"min_cpu_percentage": {
				Type:        schema.TypeFloat,
				Description: "The minimum cpu percentage for scaling",
				Optional:    true,
				Computed:    true,
			},
			"max_cpu_percentage": {
				Type:        schema.TypeFloat,
				Description: "The maximum memory percentage for scaling",
				Optional:    true,
				Computed:    true,
			},
			"enable_memory_threshold": {
				Type:        schema.TypeBool,
				Description: "Whether scaling operations based upon memory usage is enabled or not",
				Optional:    true,
				Computed:    true,
			},
			"min_memory_percentage": {
				Type:        schema.TypeFloat,
				Description: "The minimum memory percentage for scaling",
				Optional:    true,
				Computed:    true,
			},
			"max_memory_percentage": {
				Type:        schema.TypeFloat,
				Description: "The maximum memory percentage for scaling",
				Optional:    true,
				Computed:    true,
			},
			"enable_disk_threshold": {
				Type:        schema.TypeBool,
				Description: "Whether scaling operations based upon disk usage is enabled or not",
				Optional:    true,
				Computed:    true,
			},
			"min_disk_percentage": {
				Type:        schema.TypeFloat,
				Description: "The minimum disk percentage for scaling",
				Optional:    true,
				Computed:    true,
			},
			"max_disk_percentage": {
				Type:        schema.TypeFloat,
				Description: "The maximum disk percentage for scaling",
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceScaleThresholdCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"scaleThreshold": map[string]interface{}{
				"name":          d.Get("name").(string),
				"autoUp":        d.Get("auto_upscale").(bool),
				"autoDown":      d.Get("auto_downscale").(bool),
				"minCount":      d.Get("min_count").(int),
				"maxCount":      d.Get("max_count").(int),
				"cpuEnabled":    d.Get("enable_cpu_threshold").(bool),
				"minCpu":        d.Get("min_cpu_percentage").(float64),
				"maxCpu":        d.Get("max_cpu_percentage").(float64),
				"memoryEnabled": d.Get("enable_memory_threshold").(bool),
				"minMemory":     d.Get("min_memory_percentage").(float64),
				"maxMemory":     d.Get("max_memory_percentage").(float64),
				"diskEnabled":   d.Get("enable_disk_threshold").(bool),
				"minDisk":       d.Get("min_disk_percentage").(float64),
				"maxDisk":       d.Get("max_disk_percentage").(float64),
			},
		},
	}

	resp, err := client.CreateScaleThreshold(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateScaleThresholdResult)
	scaleThreshold := result.ScaleThreshold
	// Successfully created resource, now set id
	d.SetId(int64ToString(scaleThreshold.ID))

	resourceScaleThresholdRead(ctx, d, meta)
	return diags
}

func resourceScaleThresholdRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindScaleThresholdByName(name)
	} else if id != "" {
		resp, err = client.GetScaleThreshold(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("ScaleThreshold cannot be read without name or id")
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
	result := resp.Result.(*morpheus.GetScaleThresholdResult)
	scaleThreshold := result.ScaleThreshold
	if scaleThreshold != nil {
		d.SetId(int64ToString(scaleThreshold.ID))
		d.Set("name", scaleThreshold.Name)
		d.Set("auto_upscale", scaleThreshold.AutoUp)
		d.Set("auto_downscale", scaleThreshold.AutoDown)
		d.Set("min_count", scaleThreshold.MinCount)
		d.Set("max_count", scaleThreshold.MaxCount)
		d.Set("enable_cpu_threshold", scaleThreshold.CpuEnabled)
		d.Set("min_cpu_percentage", scaleThreshold.MinCpu)
		d.Set("max_cpu_percentage", scaleThreshold.MaxCpu)
		d.Set("enable_memory_threshold", scaleThreshold.MemoryEnabled)
		d.Set("min_memory_percentage", scaleThreshold.MinMemory)
		d.Set("max_memory_percentage", scaleThreshold.MaxMemory)
		d.Set("enable_disk_threshold", scaleThreshold.DiskEnabled)
		d.Set("min_disk_percentage", scaleThreshold.MinDisk)
		d.Set("max_disk_percentage", scaleThreshold.MaxDisk)
	} else {
		return diag.Errorf("read operation: scale threshold not found in response data") // should not happen
	}
	return diags
}

func resourceScaleThresholdUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"scaleThreshold": map[string]interface{}{
				"name":          d.Get("name").(string),
				"autoUp":        d.Get("auto_upscale").(bool),
				"autoDown":      d.Get("auto_downscale").(bool),
				"minCount":      d.Get("min_count").(int),
				"maxCount":      d.Get("max_count").(int),
				"cpuEnabled":    d.Get("enable_cpu_threshold").(bool),
				"minCpu":        d.Get("min_cpu_percentage").(float64),
				"maxCpu":        d.Get("max_cpu_percentage").(float64),
				"memoryEnabled": d.Get("enable_memory_threshold").(bool),
				"minMemory":     d.Get("min_memory_percentage").(float64),
				"maxMemory":     d.Get("max_memory_percentage").(float64),
				"diskEnabled":   d.Get("enable_disk_threshold").(bool),
				"minDisk":       d.Get("min_disk_percentage").(float64),
				"maxDisk":       d.Get("max_disk_percentage").(float64),
			},
		},
	}

	log.Printf("API REQUEST: %s", req)
	resp, err := client.UpdateScaleThreshold(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateScaleThresholdResult)
	account := result.ScaleThreshold
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(account.ID))
	return resourceScaleThresholdRead(ctx, d, meta)
}

func resourceScaleThresholdDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteScaleThreshold(toInt64(id), req)
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
