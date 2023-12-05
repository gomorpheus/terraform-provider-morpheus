package morpheus

import (
	"context"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBackupSetting() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus backup setting resource.",
		CreateContext: resourceBackupSettingCreate,
		ReadContext:   resourceBackupSettingRead,
		UpdateContext: resourceBackupSettingUpdate,
		DeleteContext: resourceBackupSettingDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the backup settings",
				Computed:    true,
			},
			"scheduled_backups": {
				Type:        schema.TypeBool,
				Description: "Whether automatic backups will be scheduled for provisioned instances",
				Optional:    true,
				Computed:    true,
			},
			"create_backups": {
				Type:        schema.TypeBool,
				Description: "Whether morpheus will automatically configure instances for manual or scheduled backups",
				Optional:    true,
				Computed:    true,
			},
			"backup_appliance": {
				Type:        schema.TypeBool,
				Description: "Whether a backup will be created for the Morpheus appliance database",
				Optional:    true,
				Computed:    true,
			},
			"default_backup_storage_bucket_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the storage bucket to set as the default for backups",
				Optional:    true,
				Computed:    true,
			},
			"default_backup_schedule_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the execution schedule used as the default backup schedule",
				Optional:    true,
				Computed:    true,
			},
			"retention_days": {
				Type:        schema.TypeInt,
				Description: "The number of days to retain backups",
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceBackupSettingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"backupSettings": map[string]interface{}{
				"backupsEnabled":  d.Get("scheduled_backups").(bool),
				"createBackups":   d.Get("create_backups").(bool),
				"backupAppliance": d.Get("backup_appliance").(bool),
				"retentionCount":  d.Get("retention_days").(int),
			},
		},
	}

	var defaultStorageBucketId = d.Get("default_backup_storage_bucket_id").(int)
	if defaultStorageBucketId != 0 {
		req.Body["defaultStorageBucket"] = map[string]interface{}{
			"id": defaultStorageBucketId,
		}
	}

	var defaultBackupScheduleId = d.Get("default_backup_schedule_id").(int)
	if defaultBackupScheduleId != 0 {
		req.Body["defaultSchedule"] = map[string]interface{}{
			"id": defaultBackupScheduleId,
		}
	}

	resp, err := client.UpdateBackupSettings(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.UpdateBackupSettingsResult)
	_ = result.BackupSettings
	// Successfully created resource, now set id
	d.SetId(int64ToString(1))

	resourceBackupSettingRead(ctx, d, meta)
	return diags
}

func resourceBackupSettingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error

	resp, err = client.GetBackupSettings(&morpheus.Request{})

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
	result := resp.Result.(*morpheus.GetBackupSettingsResult)
	backupSetting := result.BackupSettings
	d.SetId(int64ToString(1))
	d.Set("scheduled_backups", backupSetting.BackupsEnabled)
	d.Set("create_backups", backupSetting.CreateBackups)
	d.Set("backup_appliance", backupSetting.BackupAppliance)
	d.Set("default_backup_storage_bucket_id", backupSetting.DefaultStorageBucket.ID)
	d.Set("default_backup_schedule_id", backupSetting.DefaultSchedule.ID)
	d.Set("retention_days", backupSetting.RetentionCount)

	return diags
}

func resourceBackupSettingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"backupSettings": map[string]interface{}{
				"backupsEnabled":  d.Get("scheduled_backups").(bool),
				"createBackups":   d.Get("create_backups").(bool),
				"backupAppliance": d.Get("backup_appliance").(bool),
				"retentionCount":  d.Get("retention_days").(int),
			},
		},
	}

	var defaultStorageBucketId = d.Get("default_backup_storage_bucket_id").(int)
	if defaultStorageBucketId != 0 {
		req.Body["defaultStorageBucket"] = map[string]interface{}{
			"id": defaultStorageBucketId,
		}
	}

	var defaultBackupScheduleId = d.Get("default_backup_schedule_id").(int)
	if defaultBackupScheduleId != 0 {
		req.Body["defaultSchedule"] = map[string]interface{}{
			"id": defaultBackupScheduleId,
		}
	}

	log.Printf("API Update: %s", req)

	resp, err := client.UpdateBackupSettings(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateBackupSettingsResult)
	_ = result.BackupSettings
	// Successfully created resource, now set id
	d.SetId(int64ToString(1))

	return resourceBackupSettingRead(ctx, d, meta)
}

func resourceBackupSettingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	d.SetId("")
	return diags
}
