package morpheus

import (
	"context"
	"log"
	"strconv"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceMorpheusVirtualImages() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus virtual images data source.",
		ReadContext: dataSourceMorpheusVirtualImagesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"sort_ascending": {
				Type:        schema.TypeBool,
				Description: "Whether to sort the IDs in ascending order. Defaults to true",
				Default:     true,
				Optional:    true,
			},
			"source": {
				Type:         schema.TypeString,
				Description:  "The source of the Morpheus virtual image (User, System, Synced) (Default: User)",
				Optional:     true,
				Default:      "User",
				ValidateFunc: validation.StringInSlice([]string{"User", "System", "Synced"}, false),
			},
			"filter": {
				Type:        schema.TypeSet,
				Description: "Custom filter block as described below.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Description:  "The name of the filter. Filter names are case-sensitive. Valid names are (name, type)",
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"name", "type"}, false),
						},
						"values": {
							Type:        schema.TypeSet,
							Description: "The filter values. Filter values are case-sensitive. Filters values support the use of Golang regex and can be tested at https://regex101.com/",
							Required:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceMorpheusVirtualImagesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var sortOrder string
	var imageTypes []string
	var names []string

	if len(d.Get("filter").(*schema.Set).List()) > 0 {
		filters := d.Get("filter").(*schema.Set).List()
		for _, filter := range filters {
			filterData := filter.(map[string]interface{})
			if filterData["name"].(string) == "type" {
				for _, item := range filterData["values"].(*schema.Set).List() {
					imageTypes = append(imageTypes, item.(string))
				}
			}

			if filterData["name"].(string) == "name" {
				for _, item := range filterData["values"].(*schema.Set).List() {
					names = append(names, item.(string))
				}
			}
		}
	}

	if len(imageTypes) == 0 {
		imageTypes = append(imageTypes, "$")
	}

	if len(names) == 0 {
		names = append(names, "$")
	}

	// Sort virtual images in ascending or descending order
	if d.Get("sort_ascending").(bool) {
		sortOrder = "asc"
	} else {
		sortOrder = "desc"
	}

	output := ListAllVirtualImages(client, 200, sortOrder, d.Get("source").(string))

	var virtulImageIDs []string

	// store resource data
	for _, virtualImage := range output {
		if regexCheck(imageTypes, virtualImage.ImageType) && regexCheck(names, virtualImage.Name) {
			virtulImageIDs = append(virtulImageIDs, strconv.Itoa(int(virtualImage.ID)))
		}
	}

	d.SetId("1")
	d.Set("ids", virtulImageIDs)
	return diags
}

func ListAllVirtualImages(client *morpheus.Client, max int, sortOrder string, source string) (images []morpheus.VirtualImage) {
	// Fetch initial images
	params := make(map[string]string)
	params["max"] = strconv.Itoa(max)
	params["sort"] = "id"
	params["direction"] = sortOrder
	params["filterType"] = source
	resp, err := client.ListVirtualImages(&morpheus.Request{
		QueryParams: params,
	})
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %v", resp, err)
		} else {
			log.Printf("API FAILURE: %s - %v", resp, err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.ListVirtualImagesResult)
	pollIterations := result.Meta.Total / int64(max)
	// Add Page 1 virtual images
	images = append(images, *result.VirtualImages...)
	currentPage := 1
	for currentPage < int(pollIterations) {
		// Fetch initial images
		params := make(map[string]string)
		params["max"] = strconv.Itoa(max)
		params["sort"] = "id"
		params["direction"] = sortOrder
		params["filterType"] = source

		params["offset"] = strconv.Itoa(currentPage * max)
		resp, err := client.ListVirtualImages(&morpheus.Request{
			QueryParams: params,
		})
		if err != nil {
			if resp != nil && resp.StatusCode == 404 {
				log.Printf("API 404: %s - %v", resp, err)
			} else {
				log.Printf("API FAILURE: %s - %v", resp, err)
			}
		}
		log.Printf("API RESPONSE: %s", resp)

		result := resp.Result.(*morpheus.ListVirtualImagesResult)
		images = append(images, *result.VirtualImages...)
		currentPage++
	}
	return images
}
