package morpheus

import (
	"context"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUserGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus user group resource",
		CreateContext: resourceUserGroupCreate,
		ReadContext:   resourceUserGroupRead,
		UpdateContext: resourceUserGroupUpdate,
		DeleteContext: resourceUserGroupDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the user group",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the user group",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the user group",
				Optional:    true,
				Computed:    true,
			},
			"server_group": {
				Type:        schema.TypeString,
				Description: "The name of the Linux group to add the users to",
				Optional:    true,
				Computed:    true,
			},
			"sudo_access": {
				Type:        schema.TypeBool,
				Description: "Whether the users in the group are granted sudo permissions",
				Optional:    true,
				Computed:    true,
			},
			"user_ids": {
				Type:        schema.TypeList,
				Description: "A list of Morpheus user IDs to add to the user group",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceUserGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	userGroup := make(map[string]interface{})

	userGroup["name"] = d.Get("name").(string)
	userGroup["description"] = d.Get("description").(string)
	userGroup["sudoUser"] = d.Get("sudo_access").(bool)
	userGroup["serverGroup"] = d.Get("server_group").(string)
	userGroup["users"] = d.Get("user_ids")

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"userGroup": userGroup,
		},
	}
	resp, err := client.CreateUserGroup(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateUserGroupResult)
	userGroupResult := result.UserGroup
	// Successfully created resource, now set id
	d.SetId(int64ToString(userGroupResult.ID))

	resourceUserGroupRead(ctx, d, meta)
	return diags
}

func resourceUserGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindUserGroupByName(name)
	} else if id != "" {
		resp, err = client.GetUserGroup(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("User Group cannot be read without name or id")
	}

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
	result := resp.Result.(*morpheus.GetUserGroupResult)
	userGroup := result.UserGroup

	d.SetId(intToString(int(userGroup.ID)))
	d.Set("name", userGroup.Name)
	d.Set("description", userGroup.Description)
	d.Set("server_group", userGroup.ServerGroup)
	d.Set("sudo_access", userGroup.SudoUser)
	var users []int64
	if userGroup.Users != nil {
		// iterate over the array of tasks
		for i := 0; i < len(userGroup.Users); i++ {
			users = append(users, userGroup.Users[i].ID)
		}
	}
	userIds := matchUserIdsWithSchema(users, d.Get("user_ids").([]interface{}))
	d.Set("user_ids", userIds)

	return diags
}

func resourceUserGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	userGroup := make(map[string]interface{})

	userGroup["name"] = d.Get("name").(string)
	userGroup["description"] = d.Get("description").(string)
	userGroup["sudoUser"] = d.Get("sudo_access").(bool)
	userGroup["serverGroup"] = d.Get("server_group").(string)
	userGroup["users"] = d.Get("user_ids")

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"userGroup": userGroup,
		},
	}
	resp, err := client.UpdateUserGroup(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateUserGroupResult)
	userGroupResult := result.UserGroup

	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(userGroupResult.ID))
	return resourceUserGroupRead(ctx, d, meta)
}

func resourceUserGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteUserGroup(toInt64(id), req)
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

// This cannot currently be handled efficiently by a DiffSuppressFunc.
// See: https://github.com/hashicorp/terraform-plugin-sdk/issues/477
func matchUserIdsWithSchema(userIds []int64, declareduserIds []interface{}) []int64 {
	result := make([]int64, len(declareduserIds))

	rMap := make(map[int64]int64, len(userIds))
	for _, userId := range userIds {
		rMap[userId] = userId
	}

	for i, definedUserId := range declareduserIds {
		definedUserId := int64(definedUserId.(int))

		if v, ok := rMap[definedUserId]; ok {
			// matched node type declared by ID
			result[i] = v
			delete(rMap, v)
		}
	}
	// append unmatched node type to the result
	for _, rcpt := range rMap {
		result = append(result, rcpt)
	}
	return result
}
