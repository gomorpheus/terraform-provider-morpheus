package morpheus

import (
	"context"
	"log"
	"strconv"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMorpheusUser() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus user resource.",
		CreateContext: resourceMorpheusUserCreate,
		ReadContext:   resourceMorpheusUserRead,
		UpdateContext: resourceMorpheusUserUpdate,
		DeleteContext: resourceMorpheusUserDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the user account",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"tenant_id": {
				Description: "The ID of the tenant to create the user account in",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
			},
			"first_name": {
				Description: "The first name of the user account",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"last_name": {
				Description: "The last name of the user account",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"username": {
				Description: "The username of the user account",
				Type:        schema.TypeString,
				Required:    true,
			},
			"email": {
				Description: "The email address of the user account",
				Type:        schema.TypeString,
				Required:    true,
			},
			"password": {
				Description: "The Morpheus password for the user account (external password changes are not detected)",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
			"password_expired": {
				Description: "Set user password expiration. After the first login you will be prompted to create a new password. This attribute only works during the initial user creation and will force the user to be deleted and recreated if the attribute is changed.",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
			},
			"receive_notifications": {
				Description: "Whether notification emails will be sent to the email address associated with the user account or not",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"role_ids": {
				Description: "A list of user role ids associated with the user account",
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"linux_username": {
				Description: "The username assigned to linux instances for this user account",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"linux_password": {
				Description: "The password assigned to linux instances for this user account (external password changes are not detected)",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"linux_keypair_id": {
				Description: "The private key pair id associated with the user account for accessing linux instances",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"windows_username": {
				Description: "The username assigned to windows instances for this user account",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"windows_password": {
				Description: "The password assigned to windows instances for this user account (external password changes are not detected)",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceMorpheusUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// roles
	var roles []map[string]interface{}
	if d.Get("role_ids") != nil {
		roleList := d.Get("role_ids").([]interface{})
		// iterate over the array of roles
		for i := 0; i < len(roleList); i++ {
			row := make(map[string]interface{})
			row["id"] = roleList[i]
			roles = append(roles, row)
		}
	}

	var tenantId string
	if d.Get("tenant_id").(int) > 1 {
		tenantId = strconv.Itoa(d.Get("tenant_id").(int))
	} else {
		tenantId = "1"
	}

	req := &morpheus.Request{
		QueryParams: map[string]string{
			"accountId": tenantId,
		},
		Body: map[string]interface{}{
			"user": map[string]interface{}{
				"firstName":            d.Get("first_name").(string),
				"lastName":             d.Get("last_name").(string),
				"username":             d.Get("username").(string),
				"email":                d.Get("email").(string),
				"password":             d.Get("password").(string),
				"passwordExpired":      d.Get("password_expired").(bool),
				"receiveNotifications": d.Get("receive_notifications").(bool),
				"linuxUsername":        d.Get("linux_username").(string),
				"linuxPassword":        d.Get("linux_password").(string),
				"linuxKeyPairId":       d.Get("linux_keypair_id").(int),
				"windowsUsername":      d.Get("windows_username").(string),
				"windowsPassword":      d.Get("windows_password").(string),
				"roles":                roles,
			},
		},
	}

	resp, err := client.CreateUser(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.CreateUserResult)
	user := result.User

	// Successfully created resource, now set id
	d.SetId(int64ToString(user.ID))
	resourceMorpheusUserRead(ctx, d, meta)
	return diags
}

func resourceMorpheusUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	username := d.Get("username").(string)

	// lookup by username if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && username != "" {
		resp, err = client.FindUserByName(username)
	} else if id != "" {
		resp, err = client.GetUser(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("User cannot be read without username or id")
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
	result := resp.Result.(*morpheus.GetUserResult)
	user := result.User
	if user != nil {
		d.SetId(int64ToString(user.ID))
		d.Set("username", user.Username)
		d.Set("first_name", user.FirstName)
		d.Set("last_name", user.LastName)
		d.Set("email", user.Email)
		d.Set("receive_notifications", user.ReceiveNotifications)
		var roleIds []int
		for _, role := range user.Roles {
			roleIds = append(roleIds, int(role.ID))
		}
		d.Set("role_ids", matchUserRoleIdsWithSchema(roleIds, d.Get("role_ids").([]interface{})))
		d.Set("linux_keypair_id", user.LinuxKeyPairID)
		d.Set("linux_username", user.LinuxUsername)
		d.Set("windows_username", user.WindowsUsername)
	} else {
		return diag.Errorf("User not found in response data.") // should not happen
	}

	return diags
}

func resourceMorpheusUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	// roles
	var roles []map[string]interface{}
	if d.Get("role_ids") != nil {
		roleList := d.Get("role_ids").([]interface{})
		// iterate over the array of roles
		for i := 0; i < len(roleList); i++ {
			row := make(map[string]interface{})
			row["id"] = roleList[i]
			roles = append(roles, row)
		}
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"user": map[string]interface{}{
				"firstName":            d.Get("first_name").(string),
				"lastName":             d.Get("last_name").(string),
				"username":             d.Get("username").(string),
				"email":                d.Get("email").(string),
				"password":             d.Get("password").(string),
				"passwordExpired":      d.Get("password_expired").(bool),
				"receiveNotifications": d.Get("receive_notifications").(bool),
				"linuxUsername":        d.Get("linux_username").(string),
				"linuxPassword":        d.Get("linux_password").(string),
				"linuxKeyPairId":       d.Get("linux_keypair_id").(int),
				"windowsUsername":      d.Get("windows_username").(string),
				"windowsPassword":      d.Get("windows_password").(string),
				"roles":                roles,
			},
		},
	}

	resp, err := client.UpdateUser(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateUserResult)
	user := result.User

	// Successfully updated resource, now set id
	d.SetId(int64ToString(user.ID))
	return resourceMorpheusUserRead(ctx, d, meta)
}

func resourceMorpheusUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteUserResult(toInt64(id), req)
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
func matchUserRoleIdsWithSchema(roleIds []int, declaredRoleIds []interface{}) []int {
	result := make([]int, len(declaredRoleIds))

	rMap := make(map[int]int, len(roleIds))
	for _, roleId := range roleIds {
		rMap[roleId] = roleId
	}

	for i, declaredPrice := range declaredRoleIds {
		declaredPrice := declaredPrice.(int)

		if v, ok := rMap[declaredPrice]; ok {
			// matched price declared by ID
			result[i] = v
			delete(rMap, v)
		}
	}
	// append unmatched price set to the result
	for _, rcpt := range rMap {
		result = append(result, rcpt)
	}
	return result
}
