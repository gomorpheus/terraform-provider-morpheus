package morpheus

import (
	"context"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGuidanceSetting() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus guidance setting resource.",
		CreateContext: resourceGuidanceSettingCreate,
		ReadContext:   resourceGuidanceSettingRead,
		UpdateContext: resourceGuidanceSettingUpdate,
		DeleteContext: resourceGuidanceSettingDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the guidance settings",
				Computed:    true,
			},
			"power_settings_average_cpu": {
				Type:        schema.TypeInt,
				Description: "Shutdown will be recommended if the average CPU usage is below this value",
				Optional:    true,
				Computed:    true,
			},
			"power_settings_maximum_cpu": {
				Type:        schema.TypeInt,
				Description: "Shutdown will be recommended if the CPU usage never exceeds this value",
				Optional:    true,
				Computed:    true,
			},
			"power_settings_network_threshold": {
				Type:        schema.TypeInt,
				Description: "Shutdown will be recommended if the average network bandwidth is below this value",
				Optional:    true,
				Computed:    true,
			},
			"cpu_upsize_average_cpu": {
				Type:        schema.TypeInt,
				Description: "CPU up-size is recommended if the average CPU percentage exceeds this value",
				Optional:    true,
				Computed:    true,
			},
			"cpu_upsize_maximum_cpu": {
				Type:        schema.TypeInt,
				Description: "CPU up-size is recommended if the maximum CPU percentage exceeds this value",
				Optional:    true,
				Computed:    true,
			},
			"memory_upsize_minimum_free_memory": {
				Type:        schema.TypeInt,
				Description: "Memory up-size will be recommended if free memory dips below this value",
				Optional:    true,
				Computed:    true,
			},
			"memory_downsize_average_free_memory": {
				Type:        schema.TypeInt,
				Description: "Memory down-size is recommended if the average free memory is above this value",
				Optional:    true,
				Computed:    true,
			},
			"memory_downsize_maximum_free_memory": {
				Type:        schema.TypeInt,
				Description: "Memory down-size is recommended if free memory has never dipped below this value",
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceGuidanceSettingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	guidanceSettings := make(map[string]interface{})

	powerAverageCpu, powerAverageCpuok := d.GetOk("power_settings_average_cpu")
	if powerAverageCpuok {
		guidanceSettings["cpuAvgCutoffPower"] = powerAverageCpu
	}

	powerMaxCpu, powerMaxCpuok := d.GetOk("power_settings_maximum_cpu")
	if powerMaxCpuok {
		guidanceSettings["cpuMaxCutoffPower"] = powerMaxCpu
	}

	powerNetwork, powerNetworkok := d.GetOk("power_settings_network_threshold")
	if powerNetworkok {
		guidanceSettings["networkCutoffPower"] = powerNetwork
	}

	upsizeAverageCpu, upsizeAverageCpuok := d.GetOk("cpu_upsize_average_cpu")
	if upsizeAverageCpuok {
		guidanceSettings["cpuUpAvgStandardCutoffRightSize"] = upsizeAverageCpu
	}

	upsizeMaxCpu, upsizeMaxCpuok := d.GetOk("cpu_upsize_maximum_cpu")
	if upsizeMaxCpuok {
		guidanceSettings["cpuUpMaxStandardCutoffRightSize"] = upsizeMaxCpu
	}

	upsizeMaxMemory, upsizeMaxMemoryok := d.GetOk("memory_upsize_minimum_free_memory")
	if upsizeMaxMemoryok {
		guidanceSettings["memoryUpAvgStandardCutoffRightSize"] = upsizeMaxMemory
	}

	downsizeAverageMemory, downsizeAverageMemoryok := d.GetOk("memory_downsize_average_free_memory")
	if downsizeAverageMemoryok {
		guidanceSettings["memoryDownAvgStandardCutoffRightSize"] = downsizeAverageMemory
	}

	downsizeMaxMemory, downsizeMaxMemoryok := d.GetOk("memory_downsize_maximum_free_memory")
	if downsizeMaxMemoryok {
		guidanceSettings["memoryDownMaxStandardCutoffRightSize"] = downsizeMaxMemory
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"guidanceSettings": guidanceSettings,
		},
	}

	resp, err := client.UpdateGuidanceSettings(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.UpdateGuidanceSettingsResult)
	_ = result.GuidanceSettings
	// Successfully created resource, now set id
	d.SetId(int64ToString(1))

	resourceGuidanceSettingRead(ctx, d, meta)
	return diags
}

func resourceGuidanceSettingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error

	resp, err = client.GetGuidanceSettings(&morpheus.Request{})

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
	result := resp.Result.(*morpheus.GetGuidanceSettingsResult)
	guidanceSetting := result.GuidanceSettings
	d.SetId(int64ToString(1))
	d.Set("power_settings_average_cpu", guidanceSetting.CpuAvgCutoffPower)
	d.Set("power_settings_maximum_cpu", guidanceSetting.CpuMaxCutoffPower)
	d.Set("power_settings_network_threshold", guidanceSetting.NetworkCutoffPower)
	d.Set("cpu_upsize_average_cpu", guidanceSetting.CpuUpAvgStandardCutoffRightSize)
	d.Set("cpu_upsize_maximum_cpu", guidanceSetting.CpuUpMaxStandardCutoffRightSize)
	d.Set("memory_upsize_minimum_free_memory", guidanceSetting.MemoryUpAvgStandardCutoffRightSize)
	d.Set("memory_downsize_average_free_memory", guidanceSetting.MemoryDownAvgStandardCutoffRightSize)
	d.Set("memory_downsize_maximum_free_memory", guidanceSetting.MemoryDownMaxStandardCutoffRightSize)

	return diags
}

func resourceGuidanceSettingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	guidanceSettings := make(map[string]interface{})

	if d.HasChange("power_settings_average_cpu") {
		guidanceSettings["cpuAvgCutoffPower"] = d.Get("power_settings_average_cpu")
	}

	if d.HasChange("power_settings_maximum_cpu") {
		guidanceSettings["cpuMaxCutoffPower"] = d.Get("power_settings_maximum_cpu")
	}

	if d.HasChange("power_settings_network_threshold") {
		guidanceSettings["networkCutoffPower"] = d.Get("power_settings_network_threshold")
	}

	if d.HasChange("cpu_upsize_average_cpu") {
		guidanceSettings["cpuUpAvgStandardCutoffRightSize"] = d.Get("cpu_upsize_average_cpu")
	}

	if d.HasChange("cpu_upsize_maximum_cpu") {
		guidanceSettings["cpuUpMaxStandardCutoffRightSize"] = d.Get("cpu_upsize_maximum_cpu")
	}

	if d.HasChange("memory_upsize_minimum_free_memory") {
		guidanceSettings["memoryUpAvgStandardCutoffRightSize"] = d.Get("memory_upsize_minimum_free_memory")
	}

	if d.HasChange("memory_downsize_average_free_memory") {
		guidanceSettings["memoryDownAvgStandardCutoffRightSize"] = d.Get("memory_downsize_average_free_memory")
	}

	if d.HasChange("memory_downsize_maximum_free_memory") {
		guidanceSettings["memoryDownMaxStandardCutoffRightSize"] = d.Get("memory_downsize_maximum_free_memory")
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"guidanceSettings": guidanceSettings,
		},
	}

	resp, err := client.UpdateGuidanceSettings(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateGuidanceSettingsResult)
	_ = result.GuidanceSettings
	// Successfully created resource, now set id
	d.SetId(int64ToString(1))

	return resourceGuidanceSettingRead(ctx, d, meta)
}

func resourceGuidanceSettingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	d.SetId("")
	return diags
}
