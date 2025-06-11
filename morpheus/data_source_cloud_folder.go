package morpheus

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMorpheusCloudFolder() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus Cloud Folder data source.",
		ReadContext: dataSourceMorpheusCloudFolderRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description:   "The name of the Morpheus Cloud Folder, supply either this or the id.",
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
			"cloud_id": {
				Description: "The ID of the Morpheus Cloud.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"id": {
				Description:   "The ID of the Morpheus Cloud Folder, supply either this or the name.",
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"name"},
			},
			"external_id": {
				Description: "The external ID of the Morpheus Cloud Folder.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"type": {
				Description: "The type of the Morpheus Cloud Folder.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceMorpheusCloudFolderRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	name := d.Get("name").(string)
	id := d.Get("id").(int)
	cloudId := d.Get("cloud_id").(int)

	var resp *morpheus.Response
	var err error
	if id == 0 && name != "" {
		resp, err = getCloudFolderFromName(client, cloudId, name)
	} else if id != 0 && name == "" {
		resp, err = client.GetCloudResourceFolder(int64(cloudId), int64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Virtual image cannot be read without name or id")
	}
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %v", resp, err)
			return nil
		} else {
			return diag.FromErr(err)
		}
	}
	cloudFolder := resp.Result.(*morpheus.Folder)

	d.SetId(strconv.FormatInt(cloudFolder.ID, 10))
	d.Set("external_id", cloudFolder.ExternalId)
	d.Set("type", cloudFolder.Type)
	d.Set("name", cloudFolder.Name)

	return nil
}

func getCloudFolderFromName(client *morpheus.Client, cloudId int, name string) (*morpheus.Response, error) {
	resp, err := client.ListCloudResourceFolders(int64(cloudId), &morpheus.Request{})
	if err != nil {
		return nil, err
	}

	result := resp.Result.(*morpheus.ListCloudResourceFoldersResult)
	for _, folder := range *result.Folders {
		if folder.Name == name {
			ret := &morpheus.Response{Result: &folder}

			return ret, nil
		}
	}

	return nil, fmt.Errorf("Cloud Folder not found")
}
