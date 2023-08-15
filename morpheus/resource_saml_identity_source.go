package morpheus

import (
	"context"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSAMLIdentitySource() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a saml identity source resource",
		CreateContext: resourceSAMLIdentitySourceCreate,
		ReadContext:   resourceSAMLIdentitySourceRead,
		UpdateContext: resourceSAMLIdentitySourceUpdate,
		DeleteContext: resourceSAMLIdentitySourceDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the SAML identity source",
				Computed:    true,
			},
			"tenant_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the Morpheus tenant to associate the identity source with",
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the SAML identity source",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the SAML identity source",
				Optional:    true,
				Computed:    true,
			},
			"login_redirect_url": {
				Type:        schema.TypeString,
				Description: "This is the SAML endpoint Morpheus will redirect to when a user signs into Morpheus via SAML",
				Optional:    true,
				Computed:    true,
			},
			"logout_redirect_url": {
				Type:        schema.TypeString,
				Description: "The URL Morpheus will POST to when a SAML user logs out of Morpheus",
				Optional:    true,
				Computed:    true,
			},
			"include_saml_request_parameter": {
				Type:        schema.TypeBool,
				Description: "Whether to include the SAML request as a parameter",
				Optional:    true,
				Computed:    true,
			},
			"saml_request": {
				Type:         schema.TypeString,
				Description:  "The SAML request configuration (NoSignature, SelfSigned, CustomSignature)",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"NoSignature", "SelfSigned", "CustomSignature"}, false),
			},
			"validate_assertion_signature": {
				Type:        schema.TypeBool,
				Description: "Whether to validate the assertion signature (SAML RESPONSE field in the UI)",
				Optional:    true,
				Computed:    true,
			},
			"given_name_attribute": {
				Type:        schema.TypeString,
				Description: "SAML SP field value to map to Morpheus user First Name",
				Optional:    true,
				Computed:    true,
			},
			"surname_attribute": {
				Type:        schema.TypeString,
				Description: "SAML SP field value to map to Morpheus user Last Name",
				Optional:    true,
				Computed:    true,
			},
			"email_attribute": {
				Type:        schema.TypeString,
				Description: "SAML SP field value to map to Morpheus user email address",
				Optional:    true,
				Computed:    true,
			},
			"default_account_role_id": {
				Type:        schema.TypeInt,
				Description: "The id of the default role a user is assigned when they are in the required group or if no specific group mapping applies to the user",
				Required:    true,
			},
			"role_attribute_name": {
				Type:        schema.TypeString,
				Description: "The name of the attribute/assertion field that will map to Morpheus roles, such a MemberOf",
				Optional:    true,
				Computed:    true,
			},
			"required_role_attribute_value": {
				Type:        schema.TypeString,
				Description: "The name of the attribute/assertion field that maps to the required role",
				Optional:    true,
				Computed:    true,
			},
			"role_mapping": {
				Description: "The SAML to Morpheus Role mapping",
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
						"assertion_attribute": {
							Description: "The assertion attribute to map the role to",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"enable_role_mapping_permission": {
				Type:        schema.TypeBool,
				Description: "When enabled, Tenant users with appropriate rights to view and edit Roles will have the ability to set role mapping for the Identity Source integration",
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceSAMLIdentitySourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	identitySource := make(map[string]interface{})

	identitySource["name"] = d.Get("name").(string)
	identitySource["description"] = d.Get("description").(string)
	identitySource["type"] = "saml"

	config := make(map[string]interface{})
	config["url"] = d.Get("login_redirect_url").(string)
	config["logoutUrl"] = d.Get("logout_redirect_url").(string)
	if d.Get("include_saml_request_parameter").(bool) {
		config["doNotIncludeSAMLRequest"] = false
	} else {
		config["doNotIncludeSAMLRequest"] = true
	}

	config["SAMLSignatureMode"] = d.Get("saml_request").(string)

	config["givenNameAttribute"] = d.Get("given_name_attribute").(string)
	config["surnameAttribute"] = d.Get("surname_attribute").(string)
	config["emailAttribute"] = d.Get("email_attribute").(string)
	config["roleAttributeName"] = d.Get("role_attribute_name").(string)
	config["requiredAttributeValue"] = d.Get("required_role_attribute_value").(string)

	if d.Get("validate_assertion_signature").(bool) {
		config["doNotValidateSignature"] = false
	} else {
		config["doNotValidateSignature"] = true
	}
	identitySource["config"] = config

	defaultAccountRole := make(map[string]interface{})
	defaultAccountRole["id"] = d.Get("default_account_role_id").(int)
	identitySource["defaultAccountRole"] = defaultAccountRole

	// Role Mappings
	if d.Get("role_mapping") != "" {
		identitySource["roleMappings"] = parseSAMLRoleMappings(d.Get("role_mapping").(*schema.Set))
	}
	identitySource["allowCustomMappings"] = d.Get("enable_role_mapping_permission").(bool)

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

	resourceSAMLIdentitySourceRead(ctx, d, meta)
	return diags
}

func resourceSAMLIdentitySourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	d.Set("login_redirect_url", identitySource.Config.URL)
	d.Set("logout_redirect_url", identitySource.Config.LogoutURL)
	if identitySource.Config.DoNotIncludeSAMLRequest {
		d.Set("include_saml_request_parameter", false)
	} else {
		d.Set("include_saml_request_parameter", true)
	}
	d.Set("saml_request", identitySource.Config.SAMLSignatureMode)
	if identitySource.Config.DoNotValidateSignature {
		d.Set("validate_assertion_signature", false)
	} else {
		d.Set("validate_assertion_signature", true)
	}
	d.Set("given_name_attribute", identitySource.Config.GivenNameAttribute)
	d.Set("surname_attribute", identitySource.Config.SurnameAttribute)
	d.Set("email_attribute", identitySource.Config.EmailAttribute)
	d.Set("default_account_role_id", identitySource.DefaultAccountRole.ID)
	d.Set("role_attribute_name", identitySource.Config.RoleAttributeName)
	d.Set("required_role_attribute_value", identitySource.Config.RequiredAttributeValue)
	d.Set("enable_role_mapping_permission", identitySource.AllowCustomMappings)

	var roleMappingPayload []map[string]interface{}

	for _, roleMapping := range identitySource.RoleMappings {
		roleOutput := make(map[string]interface{})
		roleOutput["assertion_attribute"] = roleMapping.SourceRoleName
		roleOutput["role_id"] = roleMapping.MappedRole.ID
		roleOutput["role_name"] = roleMapping.MappedRole.Authority
		roleMappingPayload = append(roleMappingPayload, roleOutput)
	}
	d.Set("role_mapping", roleMappingPayload)
	return diags
}

func resourceSAMLIdentitySourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	identitySource := make(map[string]interface{})

	identitySource["name"] = d.Get("name").(string)
	identitySource["description"] = d.Get("description").(string)
	identitySource["type"] = "saml"

	config := make(map[string]interface{})
	config["url"] = d.Get("login_redirect_url").(string)
	config["logoutUrl"] = d.Get("logout_redirect_url").(string)
	if d.Get("include_saml_request_parameter").(bool) {
		config["doNotIncludeSAMLRequest"] = false
	} else {
		config["doNotIncludeSAMLRequest"] = true
	}
	config["SAMLSignatureMode"] = d.Get("saml_request").(string)

	config["givenNameAttribute"] = d.Get("given_name_attribute").(string)
	config["surnameAttribute"] = d.Get("surname_attribute").(string)
	config["emailAttribute"] = d.Get("email_attribute").(string)
	config["roleAttributeName"] = d.Get("role_attribute_name").(string)
	config["requiredAttributeValue"] = d.Get("required_role_attribute_value").(string)

	if d.Get("validate_assertion_signature").(bool) {
		config["doNotValidateSignature"] = false
	} else {
		config["doNotValidateSignature"] = true
	}
	identitySource["config"] = config

	defaultAccountRole := make(map[string]interface{})
	defaultAccountRole["id"] = d.Get("default_account_role_id").(int)
	identitySource["defaultAccountRole"] = defaultAccountRole

	// Role Mappings
	if d.Get("role_mapping") != "" {
		identitySource["roleMappings"] = parseSAMLRoleMappings(d.Get("role_mapping").(*schema.Set))
	}
	identitySource["allowCustomMappings"] = d.Get("enable_role_mapping_permission").(bool)

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
	return resourceSAMLIdentitySourceRead(ctx, d, meta)
}

func resourceSAMLIdentitySourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func parseSAMLRoleMappings(mappings *schema.Set) []map[string]interface{} {
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
			case "assertion_attribute":
				row["sourceRoleName"] = v.(string)
			}
		}
		row["mappedRole"] = mappedRole
		roleMappings = append(roleMappings, row)
	}
	return roleMappings
}
