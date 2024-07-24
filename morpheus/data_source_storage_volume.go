package morpheus

import (
	"context"
	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMorpheusStorageVolume() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus storage volume data source.",
		ReadContext: dataSourceMorpheusStorageVolumeRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Description: "The ID of the storage volume",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the storage volume",
				Computed:    true,
			},
			"active": {
				Type:        schema.TypeBool,
				Description: "Whether the storage volume is enabled or not",
				Computed:    true,
			},
			"category": {
				Type:        schema.TypeString,
				Description: "The storage volume category",
				Computed:    true,
			},
			"cloud_name": {
				Type:        schema.TypeString,
				Description: "The storage volume cloud name",
				Computed:    true,
			},
			"cloud_id": {
				Type:        schema.TypeInt,
				Description: "The storage volume cloud id",
				Computed:    true,
			},
			"datastore_name": {
				Type:        schema.TypeString,
				Description: "The storage volume datastore name",
				Computed:    true,
			},
			"datastore_id": {
				Type:        schema.TypeInt,
				Description: "The storage volume datastore id",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "The status of the storage volume",
				Computed:    true,
			},
			"source": {
				Type:        schema.TypeString,
				Description: "The associated cloud name",
				Computed:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "The storage volume type name",
				Computed:    true,
			},
			"type_id": {
				Type:        schema.TypeInt,
				Description: "The storage volume type id",
				Computed:    true,
			},
			"uuid": {
				Type:        schema.TypeString,
				Description: "The storage volume uuid",
				Computed:    true,
			},
		},
	}
}

func dataSourceMorpheusStorageVolumeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Get("id").(int)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error

	resp, err = client.GetStorageVolume(int64(id), &morpheus.Request{})
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %v", resp, err)
			return nil
		} else {
			log.Printf("API FAILURE: %s - %v", resp, err)
			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)

	// store resource data
	result := resp.Result.(*morpheus.GetStorageVolumeResult)
	storageVolume := result.StorageVolume
	if storageVolume != nil {
		d.SetId(int64ToString(int64(storageVolume.ID.(float64))))
		d.Set("name", storageVolume.Name)
		d.Set("active", storageVolume.Active)
		d.Set("category", storageVolume.Category)
		d.Set("cloud_name", storageVolume.Zone.Name)
		d.Set("cloud_id", storageVolume.Zone.ID)
		d.Set("datastore_id", storageVolume.Datastore.ID)
		d.Set("datastore_name", storageVolume.Datastore.Name)
		d.Set("status", storageVolume.Status)
		d.Set("source", storageVolume.Source)
		d.Set("type", storageVolume.Type.Name)
		d.Set("type_id", storageVolume.TypeId)
		d.Set("uuid", storageVolume.Uuid)
	} else {
		return diag.Errorf("Storage volume not found in response data.") // should not happen
	}
	return diags
}
