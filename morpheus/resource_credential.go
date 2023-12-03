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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceCredential() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus credential resource.",
		CreateContext: resourceCredentialCreate,
		ReadContext:   resourceCredentialRead,
		UpdateContext: resourceCredentialUpdate,
		DeleteContext: resourceCredentialDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the credential",
				Computed:    true,
			},
			"credential_store_integration_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the credential store integration",
				Optional:    true,
				ForceNew:    true,
			},
			"type": {
				Type:         schema.TypeString,
				Description:  "The credential type (access-key-secret, api-key, client-id-secret, email-private-key, tenant-username-keypair, username-password, username-api-key, username-keypair, username-password-keypair)",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"access-key-secret", "api-key", "client-id-secret", "email-private-key", "tenant-username-keypair", "username-password", "username-api-key", "username-keypair", "username-password-keypair"}, false),
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the credential",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the credential",
				Optional:    true,
				Computed:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the credential is enabled",
				Optional:    true,
			},
			"access_key": {
				Type:        schema.TypeString,
				Description: "The credential access key",
				Optional:    true,
			},
			"secret_key": {
				Type:        schema.TypeString,
				Description: "The credential secret key",
				Optional:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256_hash := hex.EncodeToString(h.Sum(nil))
					return strings.EqualFold(old, sha256_hash)
				},
			},
			"client_id": {
				Type:        schema.TypeString,
				Description: "The credential client id",
				Optional:    true,
			},
			"client_secret": {
				Type:        schema.TypeString,
				Description: "The credential client secret",
				Optional:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256_hash := hex.EncodeToString(h.Sum(nil))
					return strings.EqualFold(old, sha256_hash)
				},
			},
			"tenant": {
				Type:        schema.TypeString,
				Description: "The credential tenant",
				Optional:    true,
			},
			"email": {
				Type:        schema.TypeString,
				Description: "The credential email address",
				Optional:    true,
			},
			"username": {
				Type:        schema.TypeString,
				Description: "The credential username",
				Optional:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "The credential password",
				Optional:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256_hash := hex.EncodeToString(h.Sum(nil))
					return strings.EqualFold(old, sha256_hash)
				},
			},
			"api_key": {
				Type:        schema.TypeString,
				Description: "The credential api key",
				Optional:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256_hash := hex.EncodeToString(h.Sum(nil))
					return strings.EqualFold(old, sha256_hash)
				},
			},
			"key_pair_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the credential key pair",
				Optional:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCredentialCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	credential := make(map[string]interface{})
	credential["name"] = d.Get("name").(string)
	credential["description"] = d.Get("description").(string)
	credential["enabled"] = d.Get("enabled").(bool)

	integration := make(map[string]interface{})
	if d.Get("credential_store_integration_id").(int) != 0 {
		integration["id"] = d.Get("credential_store_integration_id").(int)
	}
	credential["integration"] = integration
	switch d.Get("type").(string) {
	case "access-key-secret":
		credential["type"] = "access-key-secret"
		credential["username"] = d.Get("access_key").(string)
		credential["password"] = d.Get("secret_key").(string)
	case "api-key":
		credential["type"] = "api-key"
		credential["password"] = d.Get("api_key").(string)
	case "client-id-secret":
		credential["type"] = "client-id-secret"
		credential["username"] = d.Get("client_id").(string)
		credential["password"] = d.Get("client_secret").(string)
	case "email-private-key":
		credential["type"] = "email-private-key"
		credential["username"] = d.Get("email").(string)
		keypair := make(map[string]interface{})
		keypair["id"] = d.Get("key_pair_id").(int)
		credential["authKey"] = keypair
	case "tenant-username-keypair":
		credential["type"] = "tenant-username-keypair"
		credential["authPath"] = d.Get("tenant").(string)
		credential["username"] = d.Get("username").(string)
		keypair := make(map[string]interface{})
		keypair["id"] = d.Get("key_pair_id").(int)
		credential["authKey"] = keypair
	case "username-api-key":
		credential["type"] = "username-api-key"
		credential["username"] = d.Get("username").(string)
		credential["password"] = d.Get("api_key").(string)
	case "username-keypair":
		credential["type"] = "username-keypair"
		credential["username"] = d.Get("username").(string)
		keypair := make(map[string]interface{})
		keypair["id"] = d.Get("key_pair_id").(int)
		credential["authKey"] = keypair
	case "username-password":
		credential["type"] = "username-password"
		credential["username"] = d.Get("username").(string)
		credential["password"] = d.Get("password").(string)
	case "username-password-keypair":
		credential["type"] = "username-password-keypair"
		credential["username"] = d.Get("username").(string)
		credential["password"] = d.Get("password").(string)
		keypair := make(map[string]interface{})
		keypair["id"] = d.Get("key_pair_id").(int)
		credential["authKey"] = keypair
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"credential": credential,
		},
	}

	resp, err := client.CreateCredential(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateCredentialResult)
	contact := result.Credential
	// Successfully created resource, now set id
	d.SetId(int64ToString(contact.ID))

	resourceCredentialRead(ctx, d, meta)
	return diags
}

func resourceCredentialRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindCredentialByName(name)
	} else if id != "" {
		resp, err = client.GetCredential(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Credential cannot be read without name or id")
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
	result := resp.Result.(*morpheus.GetCredentialResult)
	credential := result.Credential
	if credential != nil {
		d.SetId(int64ToString(credential.ID))
		d.Set("name", credential.Name)
		d.Set("description", credential.Description)
		d.Set("enabled", credential.Enabled)
		if credential.Integration.ID != 0 {
			d.Set("credential_store_integration_id", credential.Integration.ID)
		}
		switch credential.Type.Code {
		case "access-key-secret":
			d.Set("access_key", credential.Username)
			d.Set("secret_key", credential.PasswordHash)
		case "api-key":
			d.Set("api_key", credential.PasswordHash)
		case "client-id-secret":
			d.Set("client_id", credential.Username)
			d.Set("client_secret", credential.PasswordHash)
		case "email-private-key":
			d.Set("email", credential.Username)
			d.Set("key_pair_id", credential.AuthKey.ID)
		case "tenant-username-keypair":
			d.Set("tenant", credential.AuthPath)
			d.Set("username", credential.Username)
			d.Set("key_pair_id", credential.AuthKey.ID)
		case "username-api-key":
			d.Set("username", credential.Username)
			d.Set("api_key", credential.PasswordHash)
		case "username-keypair":
			d.Set("username", credential.Username)
			d.Set("key_pair_id", credential.AuthKey.ID)
		case "username-password":
			d.Set("username", credential.Username)
			d.Set("password", credential.PasswordHash)
		case "username-password-keypair":
			d.Set("username", credential.Username)
			d.Set("password", credential.PasswordHash)
			d.Set("key_pair_id", credential.AuthKey.ID)
		}
	} else {
		return diag.Errorf("read operation: credential not found in response data") // should not happen
	}

	return diags
}

func resourceCredentialUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	credential := make(map[string]interface{})
	credential["name"] = d.Get("name").(string)
	credential["description"] = d.Get("description").(string)
	credential["enabled"] = d.Get("enabled").(bool)

	switch d.Get("type").(string) {
	case "access-key-secret":
		credential["type"] = "access-key-secret"
		credential["username"] = d.Get("access_key").(string)
		credential["password"] = d.Get("secret_key").(string)
	case "api-key":
		credential["type"] = "api-key"
		credential["password"] = d.Get("api_key").(string)
	case "client-id-secret":
		credential["type"] = "client-id-secret"
		credential["username"] = d.Get("client_id").(string)
		credential["password"] = d.Get("client_secret").(string)
	case "email-private-key":
		credential["type"] = "email-private-key"
		credential["username"] = d.Get("email").(string)
		keypair := make(map[string]interface{})
		keypair["id"] = d.Get("key_pair_id").(int)
		credential["authKey"] = keypair
	case "tenant-username-keypair":
		credential["type"] = "tenant-username-keypair"
		credential["authPath"] = d.Get("tenant").(string)
		credential["username"] = d.Get("username").(string)
		keypair := make(map[string]interface{})
		keypair["id"] = d.Get("key_pair_id").(int)
		credential["authKey"] = keypair
	case "username-api-key":
		credential["type"] = "username-api-key"
		credential["username"] = d.Get("username").(string)
		credential["password"] = d.Get("api_key").(string)
	case "username-keypair":
		credential["type"] = "username-keypair"
		credential["username"] = d.Get("username").(string)
		keypair := make(map[string]interface{})
		keypair["id"] = d.Get("key_pair_id").(int)
		credential["authKey"] = keypair
	case "username-password":
		credential["type"] = "username-password"
		credential["username"] = d.Get("username").(string)
		credential["password"] = d.Get("password").(string)
	case "username-password-keypair":
		credential["type"] = "username-password-keypair"
		credential["username"] = d.Get("username").(string)
		credential["password"] = d.Get("password").(string)
		keypair := make(map[string]interface{})
		keypair["id"] = d.Get("key_pair_id").(int)
		credential["authKey"] = keypair
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"credential": credential,
		},
	}
	resp, err := client.UpdateCredential(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateCredentialResult)
	contact := result.Credential
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(contact.ID))
	return resourceCredentialRead(ctx, d, meta)
}

func resourceCredentialDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteCredential(toInt64(id), req)
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
