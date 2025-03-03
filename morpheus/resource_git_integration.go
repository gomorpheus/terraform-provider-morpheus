package morpheus

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGitIntegration() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a git integration resource",
		CreateContext: resourceGitIntegrationCreate,
		ReadContext:   resourceGitIntegrationRead,
		UpdateContext: resourceGitIntegrationUpdate,
		DeleteContext: resourceGitIntegrationDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the git integration",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the git integration",
				Required:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the git integration is enabled",
				Optional:    true,
				Computed:    true,
			},
			"url": {
				Type:        schema.TypeString,
				Description: "The url of the git repository",
				Required:    true,
			},
			"default_branch": {
				Type:        schema.TypeString,
				Description: "The default branch of the git repository",
				Optional:    true,
				Computed:    true,
			},
			"username": {
				Type:        schema.TypeString,
				Description: "The username of the account used to authenticate to the git repository",
				Optional:    true,
				Computed:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "The password of the account used to authenticate to the git repository",
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
			"access_token": {
				Type:        schema.TypeString,
				Description: "The access token of the account used to authenticate to the git repository",
				Optional:    true,
				Computed:    true,
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
				Description: "The ID of the key pair used to authenticate to the git repository",
				Optional:    true,
				Computed:    true,
			},
			"enable_git_caching": {
				Type:        schema.TypeBool,
				Description: "Whether the git repository is cached",
				Optional:    true,
				Computed:    true,
			},
			"repository_ids": {
				Computed:    true,
				Type:        schema.TypeMap,
				Description: "A map of git repository ids for use with integrations that reference a git repository",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceGitIntegrationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	integration := make(map[string]interface{})

	integration["name"] = d.Get("name").(string)
	integration["enabled"] = d.Get("enabled").(bool)
	integration["type"] = "git"
	integration["serviceUrl"] = d.Get("url").(string)
	integration["serviceUsername"] = d.Get("username").(string)
	integration["servicePassword"] = d.Get("password").(string)
	integration["serviceToken"] = d.Get("access_token").(string)
	integration["serviceKey"] = d.Get("key_pair_id").(int)

	config := make(map[string]interface{})
	config["defaultBranch"] = d.Get("default_branch").(string)
	config["cacheEnabled"] = d.Get("enable_git_caching").(bool)

	integration["config"] = config

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"integration": integration,
		},
	}

	resp, err := client.CreateIntegration(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateIntegrationResult)
	integrationResult := result.Integration
	// Successfully created resource, now set id
	d.SetId(int64ToString(integrationResult.ID))

	if err := retry.RetryContext(ctx, 1*time.Minute, func() *retry.RetryError {
		resp, err := client.Execute(&morpheus.Request{
			Method:      "GET",
			Path:        fmt.Sprintf("/api/options/codeRepositories?integrationId=%d", integrationResult.ID),
			QueryParams: map[string]string{},
		})
		if err != nil {
			tflog.Error(ctx, "API", map[string]any{"resp": resp.String(), "err": err})
			return retry.NonRetryableError(err)
		}
		tflog.Info(ctx, "API", map[string]any{"resp": resp.String()})

		var itemResponsePayload CodeRepositories
		if err := json.Unmarshal(resp.Body, &itemResponsePayload); err != nil {
			return retry.NonRetryableError(err)
		}
		repo_ids := make(map[string]int)
		for _, v := range itemResponsePayload.Data {
			repo_ids[v.Name] = v.Value
		}
		if len(repo_ids) == 0 {
			return retry.RetryableError(errors.New("expected codeRepositories to be created"))
		}
		return nil
	}); err != nil {
		return diag.FromErr(err)
	}

	return resourceGitIntegrationRead(ctx, d, meta)
}

func resourceGitIntegrationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindIntegrationByName(name)
	} else if id != "" {
		resp, err = client.GetIntegration(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Integration cannot be read without name or id")
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
	result := resp.Result.(*morpheus.GetIntegrationResult)
	integration := result.Integration
	d.SetId(int64ToString(integration.ID))
	d.Set("name", integration.Name)
	d.Set("enabled", integration.Enabled)
	d.Set("url", integration.URL)
	d.Set("username", integration.Username)
	d.Set("password", integration.PasswordHash)
	d.Set("access_token", integration.TokenHash)
	d.Set("key_pair_id", integration.ServiceKey.ID)
	d.Set("default_branch", integration.Config.DefaultBranch)
	d.Set("enable_git_caching", integration.Config.CacheEnabled)

	resp, err = client.Execute(&morpheus.Request{
		Method:      "GET",
		Path:        fmt.Sprintf("/api/options/codeRepositories?integrationId=%d", integration.ID),
		QueryParams: map[string]string{},
	})
	if err != nil {
		log.Println("API ERROR: ", err)
	}
	log.Println("API RESPONSE:", resp)
	repo_ids := make(map[string]int)

	var itemResponsePayload CodeRepositories
	json.Unmarshal(resp.Body, &itemResponsePayload)
	for _, v := range itemResponsePayload.Data {
		repo_ids[v.Name] = v.Value
	}
	d.Set("repository_ids", repo_ids)

	return diags
}

func resourceGitIntegrationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	integration := make(map[string]interface{})

	integration["name"] = d.Get("name").(string)
	integration["enabled"] = d.Get("enabled").(bool)
	integration["type"] = "git"
	integration["serviceUrl"] = d.Get("url").(string)
	integration["serviceUsername"] = d.Get("username").(string)
	integration["servicePassword"] = d.Get("password").(string)
	integration["serviceToken"] = d.Get("access_token").(string)
	integration["serviceKey"] = d.Get("key_pair_id").(int)

	config := make(map[string]interface{})
	config["defaultBranch"] = d.Get("default_branch").(string)
	config["cacheEnabled"] = d.Get("enable_git_caching").(bool)

	integration["config"] = config

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"integration": integration,
		},
	}

	resp, err := client.UpdateIntegration(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateIntegrationResult)
	integrationResult := result.Integration

	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(integrationResult.ID))
	return resourceGitIntegrationRead(ctx, d, meta)
}

func resourceGitIntegrationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteIntegration(toInt64(id), req)
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

type CodeRepositories struct {
	Success bool `json:"success"`
	Data    []struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	} `json:"data"`
}
