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

func resourceActiveDirectoryIdentitySource() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides an active directory identity source resource",
		CreateContext: resourceActiveDirectoryIdentitySourceCreate,
		ReadContext:   resourceActiveDirectoryIdentitySourceRead,
		UpdateContext: resourceActiveDirectoryIdentitySourceUpdate,
		DeleteContext: resourceActiveDirectoryIdentitySourceDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the active directory identity source",
				Computed:    true,
			},
			"tenant_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the Morpheus tenant to associate the identity source with",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the active directory identity source",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the active directory identity source",
				Optional:    true,
				Computed:    true,
			},
			"ad_server": {
				Type:        schema.TypeString,
				Description: "The IP address or hostname of the active directory domain controller",
				Required:    true,
			},
			"domain": {
				Type:        schema.TypeString,
				Description: "The name of the active directory domain",
				Required:    true,
			},
			"use_ssl": {
				Type:        schema.TypeBool,
				Description: "Whether to use SSL when connecting to the domain controller",
				Optional:    true,
				Computed:    true,
			},
			"binding_username": {
				Type:        schema.TypeString,
				Description: "The username of the account used to authenticate to the domain",
				Required:    true,
			},
			"binding_password": {
				Type:        schema.TypeString,
				Description: "The password of the account used to authenticate to the domain",
				Required:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256_hash := hex.EncodeToString(h.Sum(nil))
					return strings.ToLower(old) == strings.ToLower(sha256_hash)
				},
			},
			"required_group": {
				Type:        schema.TypeString,
				Description: "The active directory group users must be in to access Morpheus",
				Optional:    true,
				Computed:    true,
			},
			"search_member_groups": {
				Type:        schema.TypeBool,
				Description: "The path in the repository of the Ansible playbooks relative to the Git url",
				Optional:    true,
				Computed:    true,
			},
			"default_account_role_id": {
				Type:        schema.TypeInt,
				Description: "The id of the default role a user is assigned when they are in the required group or if no specific group mapping applies to the user",
				Required:    true,
			},
			"enable_role_mapping_permission": {
				Type:        schema.TypeBool,
				Description: "When enabled, Tenant users with appropriate rights to view and edit Roles will have the ability to set role mapping for the Identity Source integration",
				Optional:    true,
				Computed:    true,
			},
			"role_mapping": {
				Description: "The Active Directory to Morpheus Role mapping",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_id": {
							Description: "The id of the Morpheus role to map to",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"role_name": {
							Description: "The name or authority of the Morpheus role to map to",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"active_directory_group_name": {
							Description: "The name of the active directory role to map to",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"active_directory_group_fqn": {
							Description: "The fully qualified name of the active directory role to map to (i.e. - CN=Administrators,CN=Builtin,DC=contoso,DC=com)",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceActiveDirectoryIdentitySourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	identitySource := make(map[string]interface{})

	identitySource["name"] = d.Get("name").(string)
	identitySource["description"] = d.Get("description").(string)
	identitySource["type"] = "activeDirectory"

	config := make(map[string]interface{})
	config["url"] = d.Get("ad_server").(string)
	config["domain"] = d.Get("domain").(string)
	config["useSSL"] = d.Get("use_ssl").(bool)
	config["bindingUsername"] = d.Get("binding_username").(string)
	config["bindingPassword"] = d.Get("binding_password").(string)
	config["requiredGroup"] = d.Get("required_group").(string)
	config["searchMemberGroups"] = d.Get("search_member_groups").(bool)
	config["allowCustomMappings"] = d.Get("enable_role_mapping_permission").(bool)

	identitySource["config"] = config

	defaultAccountRole := make(map[string]interface{})
	defaultAccountRole["id"] = d.Get("default_account_role_id").(int)
	identitySource["defaultAccountRole"] = defaultAccountRole

	// Role Mappings
	if d.Get("role_mapping") != "" {
		identitySource["roleMappings"] = parseRoleMappings(d.Get("role_mapping").(*schema.Set))
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"userSource": identitySource,
		},
	}

	resp, err := client.CreateIdentitySource(int64(d.Get("tenant_id").(int)), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateIdentitySourceResult)
	identitySourceResult := result.IdentitySource
	// Successfully created resource, now set id
	d.SetId(int64ToString(identitySourceResult.ID))

	resourceActiveDirectoryIdentitySourceRead(ctx, d, meta)
	return diags
}

func resourceActiveDirectoryIdentitySourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindIdentitySourceByName(name)
	} else if id != "" {
		resp, err = client.GetIdentitySource(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Identity source cannot be read without name or id")
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
	result := resp.Result.(*morpheus.GetIdentitySourceResult)
	identitySource := result.IdentitySource
	d.SetId(int64ToString(identitySource.ID))
	d.Set("name", identitySource.Name)
	d.Set("description", identitySource.Description)
	d.Set("ad_server", identitySource.Config.URL)
	d.Set("domain", identitySource.Config.Domain)
	d.Set("use_ssl", identitySource.Config.UseSSL)
	d.Set("binding_username", identitySource.Config.BindingUsername)
	d.Set("binding_password", identitySource.Config.BindingPasswordHash)
	d.Set("required_group", identitySource.Config.RequiredGroup)
	d.Set("search_member_groups", identitySource.Config.SearchMemberGroups)
	d.Set("enable_role_mapping_permission", identitySource.AllowCustomMappings)
	d.Set("default_account_role_id", identitySource.DefaultAccountRole.ID)

	var roleMappingPayload []map[string]interface{}

	for _, roleMapping := range identitySource.RoleMappings {
		roleOutput := make(map[string]interface{})
		roleOutput["active_directory_group_fqn"] = roleMapping.SourceRoleFqn
		roleOutput["active_directory_group_name"] = roleMapping.SourceRoleName
		roleOutput["role_id"] = roleMapping.MappedRole.ID
		roleOutput["role_name"] = roleMapping.MappedRole.Authority
		roleMappingPayload = append(roleMappingPayload, roleOutput)
	}
	d.Set("role_mapping", roleMappingPayload)
	return diags
}

func resourceActiveDirectoryIdentitySourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	identitySource := make(map[string]interface{})

	identitySource["name"] = d.Get("name").(string)
	identitySource["description"] = d.Get("description").(string)
	identitySource["type"] = "activeDirectory"

	config := make(map[string]interface{})
	config["url"] = d.Get("ad_server").(string)
	config["domain"] = d.Get("domain").(string)
	config["useSSL"] = d.Get("use_ssl").(bool)
	config["bindingUsername"] = d.Get("binding_username").(string)
	if d.HasChange("binding_password") {
		config["bindingPassword"] = d.Get("binding_password").(string)
	}
	config["requiredGroup"] = d.Get("required_group").(string)
	config["searchMemberGroups"] = d.Get("search_member_groups").(bool)
	config["allowCustomMappings"] = d.Get("enable_role_mapping_permission").(bool)

	identitySource["config"] = config

	defaultAccountRole := make(map[string]interface{})
	defaultAccountRole["id"] = d.Get("default_account_role_id").(int)
	identitySource["defaultAccountRole"] = defaultAccountRole

	// Role Mappings
	if d.Get("role_mapping") != "" {
		identitySource["roleMappings"] = parseRoleMappings(d.Get("role_mapping").(*schema.Set))
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"userSource": identitySource,
		},
	}

	resp, err := client.UpdateIdentitySource(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateIdentitySourceResult)
	identitySourceResult := result.IdentitySource

	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(identitySourceResult.ID))
	return resourceActiveDirectoryIdentitySourceRead(ctx, d, meta)
}

func resourceActiveDirectoryIdentitySourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteIdentitySource(toInt64(id), req)
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

func parseRoleMappings(mappings *schema.Set) []map[string]interface{} {
	var roleMappings []map[string]interface{}
	// iterate over the array of roleMappings
	for _, mapping := range mappings.List() {
		row := make(map[string]interface{})
		mappedRole := make(map[string]interface{})
		mappingConfig := mapping.(map[string]interface{})
		for k, v := range mappingConfig {
			switch k {
			case "role_id":
				mappedRole["id"] = v.(int)
			case "role_name":
				mappedRole["authority"] = v.(string)
			case "active_directory_group_name":
				row["sourceRoleName"] = v.(string)
			case "active_directory_group_fqn":
				row["sourceRoleFqn"] = v.(string)
			}
		}
		row["mappedRole"] = mappedRole
		roleMappings = append(roleMappings, row)
	}
	return roleMappings
}
