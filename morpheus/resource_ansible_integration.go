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

func resourceAnsibleIntegration() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides an ansible integration resource",
		CreateContext: resourceAnsibleIntegrationCreate,
		ReadContext:   resourceAnsibleIntegrationRead,
		UpdateContext: resourceAnsibleIntegrationUpdate,
		DeleteContext: resourceAnsibleIntegrationDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the ansible integration",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the ansible integration",
				Required:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the ansible integration is enabled",
				Optional:    true,
				Computed:    true,
			},
			"url": {
				Type:        schema.TypeString,
				Description: "The url of the ansible repository",
				Required:    true,
			},
			"default_branch": {
				Type:        schema.TypeString,
				Description: "The default branch of the ansible repository",
				Optional:    true,
				Computed:    true,
			},
			"playbooks_path": {
				Type:        schema.TypeString,
				Description: "The path in the repository of the Ansible playbooks relative to the Git url",
				Optional:    true,
				Computed:    true,
			},
			"roles_path": {
				Type:        schema.TypeString,
				Description: "The path in the repository of the Ansible roles relative to the Git url",
				Optional:    true,
				Computed:    true,
			},
			"group_variables_path": {
				Type:        schema.TypeString,
				Description: "The path in the repository of the Ansible group variables relative to the Git url",
				Optional:    true,
				Computed:    true,
			},
			"host_variables_path": {
				Type:        schema.TypeString,
				Description: "The path in the repository of the Ansible host variables relative to the Git url",
				Optional:    true,
				Computed:    true,
			},
			"enable_ansible_galaxy_install": {
				Type:        schema.TypeBool,
				Description: "Whether to install the Ansible roles defined in the requirements.yml",
				Optional:    true,
				Computed:    true,
			},
			"enable_verbose_logging": {
				Type:        schema.TypeBool,
				Description: "Whether verbose logging is used during the execution of the ansible playbook",
				Optional:    true,
				Computed:    true,
			},
			"enable_agent_command_bus": {
				Type:        schema.TypeBool,
				Description: "Whether the agent command bus is used to execute the ansible playbook",
				Optional:    true,
				Computed:    true,
			},
			"username": {
				Type:        schema.TypeString,
				Description: "The username of the account used to authenticate to the ansible repository",
				Optional:    true,
				Computed:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "The password of the account used to authenticate to the ansible repository",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256_hash := hex.EncodeToString(h.Sum(nil))
					return strings.ToLower(old) == strings.ToLower(sha256_hash)
				},
				DiffSuppressOnRefresh: true,
			},
			"access_token": {
				Type:        schema.TypeString,
				Description: "The access token of the account used to authenticate to the ansible repository",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256_hash := hex.EncodeToString(h.Sum(nil))
					return strings.ToLower(old) == strings.ToLower(sha256_hash)
				},
			},
			"key_pair_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the key pair used to authenticate to the ansible repository",
				Optional:    true,
				Computed:    true,
			},
			"enable_git_caching": {
				Type:        schema.TypeBool,
				Description: "Whether the git repository is cached",
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceAnsibleIntegrationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	integration := make(map[string]interface{})

	integration["name"] = d.Get("name").(string)
	integration["enabled"] = d.Get("enabled").(bool)
	integration["type"] = "ansible"
	integration["serviceUrl"] = d.Get("url").(string)
	integration["serviceUsername"] = d.Get("username").(string)
	integration["servicePassword"] = d.Get("password").(string)
	integration["serviceToken"] = d.Get("access_token").(string)
	integration["serviceKey"] = d.Get("key_pair_id").(int)

	config := make(map[string]interface{})
	config["defaultBranch"] = d.Get("default_branch").(string)
	config["cacheEnabled"] = d.Get("enable_git_caching").(bool)
	config["ansiblePlaybooks"] = d.Get("playbooks_path").(string)
	config["ansibleRoles"] = d.Get("roles_path").(string)
	config["ansibleGroupVars"] = d.Get("group_variables_path").(string)
	config["ansibleHostVars"] = d.Get("host_variables_path").(string)
	config["ansibleGalaxyEnabled"] = d.Get("enable_ansible_galaxy_install").(bool)
	config["ansibleVerbose"] = d.Get("enable_verbose_logging").(bool)
	config["ansibleCommandBus"] = d.Get("enable_agent_command_bus").(bool)

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

	resourceAnsibleIntegrationRead(ctx, d, meta)
	return diags
}

func resourceAnsibleIntegrationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
			return diag.FromErr(err)
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
	d.Set("playbooks_path", integration.Config.AnsiblePlaybooks)
	d.Set("roles_path", integration.Config.AnsibleRoles)
	d.Set("group_variables_path", integration.Config.AnsibleGroupVars)
	d.Set("host_variables_path", integration.Config.AnsibleHostVars)
	d.Set("enable_ansible_galaxy_install", integration.Config.AnsibleGalaxyEnabled)
	d.Set("enable_verbose_logging", integration.Config.AnsibleVerbose)
	d.Set("enable_agent_command_bus", integration.Config.AnsibleCommandBus)
	return diags
}

func resourceAnsibleIntegrationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	integration := make(map[string]interface{})

	integration["name"] = d.Get("name").(string)
	integration["enabled"] = d.Get("enabled").(bool)
	integration["type"] = "ansible"
	integration["serviceUrl"] = d.Get("url").(string)
	integration["serviceUsername"] = d.Get("username").(string)
	integration["servicePassword"] = d.Get("password").(string)
	integration["serviceToken"] = d.Get("access_token").(string)
	integration["serviceKey"] = d.Get("key_pair_id").(int)

	config := make(map[string]interface{})
	config["defaultBranch"] = d.Get("default_branch").(string)
	config["cacheEnabled"] = d.Get("enable_git_caching").(bool)
	config["ansiblePlaybooks"] = d.Get("playbooks_path").(string)
	config["ansibleRoles"] = d.Get("roles_path").(string)
	config["ansibleGroupVars"] = d.Get("group_variables_path").(string)
	config["ansibleHostVars"] = d.Get("host_variables_path").(string)
	config["ansibleGalaxyEnabled"] = d.Get("enable_ansible_galaxy_install").(bool)
	config["ansibleVerbose"] = d.Get("enable_verbose_logging").(bool)
	config["ansibleCommandBus"] = d.Get("enable_agent_command_bus").(bool)

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
	return resourceAnsibleIntegrationRead(ctx, d, meta)
}

func resourceAnsibleIntegrationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
