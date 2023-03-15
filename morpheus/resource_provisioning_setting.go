package morpheus

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProvisioningSetting() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus provisioning setting resource.",
		CreateContext: resourceProvisioningSettingCreate,
		ReadContext:   resourceProvisioningSettingRead,
		UpdateContext: resourceProvisioningSettingUpdate,
		DeleteContext: resourceProvisioningSettingDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the provisioning settings",
				Computed:    true,
			},
			"allow_zone_selection": {
				Type:        schema.TypeBool,
				Description: "Displays or hides Cloud Selection dropdown in Provisioning wizard.",
				Optional:    true,
				Computed:    true,
			},
			"allow_host_selection": {
				Type:        schema.TypeBool,
				Description: "Displays or hides Host Selection dropdown in Provisioning wizard.",
				Optional:    true,
				Computed:    true,
			},
			"require_environments": {
				Type:        schema.TypeBool,
				Description: "Forces users to select and Environment during provisioning",
				Optional:    true,
				Computed:    true,
			},
			"show_pricing": {
				Type:        schema.TypeBool,
				Description: "Displays or hides Pricing in Provisioning wizard and Instance and Host detail pages.",
				Optional:    true,
				Computed:    true,
			},
			"hide_datastore_stats": {
				Type:        schema.TypeBool,
				Description: "Hides Datastore utilization and size stats in provisioning and app wizards.",
				Optional:    true,
				Computed:    true,
			},
			"cross_tenant_naming_policies": {
				Type:        schema.TypeBool,
				Description: "Enable for the sequence value in naming policies to apply across tenants.",
				Optional:    true,
				Computed:    true,
			},
			"reuse_sequence": {
				Type:        schema.TypeBool,
				Description: "When enabled, sequence numbers can be reused when Instances are removed. Deselect this option and Morpheus will track issued sequence numbers and use the next available number each time.",
				Optional:    true,
				Computed:    true,
			},
			"show_console_keyboard_settings": {
				Type:        schema.TypeBool,
				Description: "",
				Optional:    true,
				Computed:    true,
			},
			"cloudinit_username": {
				Type:        schema.TypeString,
				Description: "User to be added to Linux Instances during provisioning.",
				Optional:    true,
				Computed:    true,
			},
			"cloudinit_password": {
				Type:        schema.TypeString,
				Description: "Password to be set for the Cloud-Init Linux user.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256_hash := hex.EncodeToString(h.Sum(nil))
					return strings.EqualFold(old, sha256_hash)
				},
				DiffSuppressOnRefresh: true,
			},
			// "cloudinit_keypair_id": {
			// 	Type:        schema.TypeInt,
			// 	Description: "ID of the keypair to be added for the Cloud-Init Linux user.",
			// 	Optional:    true,
			// 	Computed:    true,
			// 	Sensitive:   true,
			// },
			"windows_password": {
				Type:        schema.TypeString,
				Description: "Password to be set for the Windows Administrator User during provisioning.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256_hash := hex.EncodeToString(h.Sum(nil))
					return strings.EqualFold(old, sha256_hash)
				},
				DiffSuppressOnRefresh: true,
			},
			"pxe_root_password": {
				Type:        schema.TypeString,
				Description: "Password to be set for Root during PXE Boots.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256_hash := hex.EncodeToString(h.Sum(nil))
					return strings.EqualFold(old, sha256_hash)
				},
				DiffSuppressOnRefresh: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceProvisioningSettingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"provisioningSettings": map[string]interface{}{
				"allowZoneSelection":        d.Get("allow_zone_selection").(bool),
				"allowServerSelection":      d.Get("allow_host_selection").(bool),
				"requireEnvironments":       d.Get("require_environments").(bool),
				"showPricing":               d.Get("show_pricing").(bool),
				"hideDatastoreStats":        d.Get("hide_datastore_stats").(bool),
				"crossTenantNamingPolicies": d.Get("cross_tenant_naming_policies").(bool),
				"reuseSequence":             d.Get("reuse_sequence").(bool),
				"cloudInitUsername":         d.Get("cloudinit_username").(string),
				"cloudInitPassword":         d.Get("cloudinit_password").(string),
				"windowsPassword":           d.Get("windows_password").(string),
				"pxeRootPassword":           d.Get("pxe_root_password").(string),
			},
		},
	}

	// var cloudInitKeypairId = d.Get("cloudinit_keypair_id").(int)
	// if cloudInitKeypairId != 0 {
	// 	req.Body["cloudInitKeyPair"] = map[string]interface{}{
	// 		"id": cloudInitKeypairId,
	// 	}
	// }

	resp, err := client.UpdateProvisioningSettings(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.UpdateProvisioningSettingsResult)
	_ = result.ProvisioningSettings

	// Successfully created resource, now set id
	d.SetId(int64ToString(1))

	resourceProvisioningSettingRead(ctx, d, meta)
	return diags
}

func resourceProvisioningSettingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error

	resp, err = client.GetProvisioningSettings(&morpheus.Request{})

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
	result := resp.Result.(*morpheus.GetProvisioningSettingsResult)
	provisioningSetting := result.ProvisioningSettings
	d.SetId(int64ToString(1))

	d.Set("allow_zone_selection", provisioningSetting.AllowZoneSelection)
	d.Set("allow_host_selection", provisioningSetting.AllowServerSelection)
	d.Set("require_environments", provisioningSetting.RequireEnvironments)
	d.Set("show_pricing", provisioningSetting.ShowPricing)
	d.Set("hide_datastore_stats", provisioningSetting.HideDatastoreStats)
	d.Set("cross_tenant_naming_policies", provisioningSetting.CrossTenantNamingPolicies)
	d.Set("reuse_sequence", provisioningSetting.ReuseSequence)
	d.Set("show_console_keyboard_settings", provisioningSetting.ShowConsoleKeyboardSettings)
	d.Set("cloudinit_username", provisioningSetting.CloudInitUsername)
	d.Set("cloudinit_password", provisioningSetting.CloudInitPasswordHash)
	// d.Set("cloudinit_keypair_id", provisioningSetting.Cloudinitkeypair.ID)
	d.Set("windows_password", provisioningSetting.WindowsPasswordHash)
	d.Set("pxe_root_password", provisioningSetting.PXERootPasswordHash)

	return diags
}

func resourceProvisioningSettingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"provisioningSettings": map[string]interface{}{
				"allowZoneSelection":        d.Get("allow_zone_selection").(bool),
				"allowServerSelection":      d.Get("allow_host_selection").(bool),
				"requireEnvironments":       d.Get("require_environments").(bool),
				"showPricing":               d.Get("show_pricing").(bool),
				"hideDatastoreStats":        d.Get("hide_datastore_stats").(bool),
				"crossTenantNamingPolicies": d.Get("cross_tenant_naming_policies").(bool),
				"reuseSequence":             d.Get("reuse_sequence").(bool),
				"cloudInitUsername":         d.Get("cloudinit_username").(string),
				"cloudInitPassword":         d.Get("cloudinit_password").(string),
				"windowsPassword":           d.Get("windows_password").(string),
				"pxeRootPassword":           d.Get("pxe_root_password").(string),
			},
		},
	}

	// var cloudInitKeypairId = d.Get("cloudinit_keypair_id").(int)
	// if cloudInitKeypairId != 0 {
	// 	req.Body["cloudInitKeyPair"] = map[string]interface{}{
	// 		"id": cloudInitKeypairId,
	// 	}
	// }

	resp, err := client.UpdateProvisioningSettings(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateProvisioningSettingsResult)
	_ = result.ProvisioningSettings

	// Successfully created resource, now set id
	d.SetId(int64ToString(1))

	return resourceProvisioningSettingRead(ctx, d, meta)
}

func resourceProvisioningSettingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	d.SetId("")
	return diags
}
