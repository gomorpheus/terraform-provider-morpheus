package morpheus

import (
	"context"
	"strings"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceScriptTemplate() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus script template resource",
		CreateContext: resourceScriptTemplateCreate,
		ReadContext:   resourceScriptTemplateRead,
		UpdateContext: resourceScriptTemplateUpdate,
		DeleteContext: resourceScriptTemplateDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the script template",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the script template",
				Required:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "The organization labels associated with the script template (Only supported on Morpheus 5.5.3 or higher)",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"script_type": {
				Type:         schema.TypeString,
				Description:  "The type of the script template (powershell, bash)",
				ValidateFunc: validation.StringInSlice([]string{"powershell", "bash"}, false),
				Required:     true,
			},
			"script_phase": {
				Type:         schema.TypeString,
				Description:  "The phase that the script should be run during (start, stop, preProvision, provision, postProvision, preDeploy, deploy, reconfigure, teardown)",
				ValidateFunc: validation.StringInSlice([]string{"start", "stop", "preProvision", "provision", "postProvision", "preDeploy", "deploy", "reconfigure", "teardown"}, false),
				Required:     true,
			},
			"script_content": {
				Type:        schema.TypeString,
				Description: "The content of the script template",
				Optional:    true,
				StateFunc: func(v interface{}) string {
					payload := strings.TrimSuffix(v.(string), "\n")
					return payload
				},
			},
			"run_as_user": {
				Type:        schema.TypeString,
				Description: "The name of the user account the script should run as",
				Optional:    true,
			},
			"sudo": {
				Type:        schema.TypeBool,
				Description: "Whether the script should run with sudo privileges",
				Optional:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceScriptTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		for _, s := range attr.(*schema.Set).List() {
			labelsPayload = append(labelsPayload, s.(string))
		}
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"containerScript": map[string]interface{}{
				"name":        name,
				"labels":      labelsPayload,
				"scriptType":  d.Get("script_type").(string),
				"scriptPhase": d.Get("script_phase").(string),
				"script":      d.Get("script_content").(string),
				"runAsUser":   d.Get("run_as_user").(string),
				"sudoUser":    d.Get("sudo").(bool),
			},
		},
	}
	resp, err := client.CreateScriptTemplate(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateScriptTemplateResult)
	scriptTemplate := result.ScriptTemplate
	// Successfully created resource, now set id
	d.SetId(int64ToString(scriptTemplate.ID))

	resourceScriptTemplateRead(ctx, d, meta)
	return diags
}

func resourceScriptTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindScriptTemplateByName(name)
	} else if id != "" {
		resp, err = client.GetScriptTemplate(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Script template cannot be read without name or id")
	}

	if err != nil {
		// 404 is ok?
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
	result := resp.Result.(*morpheus.GetScriptTemplateResult)
	scriptTemplate := result.ScriptTemplate
	d.SetId(int64ToString(scriptTemplate.ID))
	d.Set("name", scriptTemplate.Name)
	d.Set("labels", scriptTemplate.Labels)
	d.Set("script_phase", scriptTemplate.ScriptPhase)
	d.Set("script_type", scriptTemplate.ScriptType)
	d.Set("script_content", scriptTemplate.Script)
	d.Set("run_as_user", scriptTemplate.RunAsUser)
	d.Set("sudo", scriptTemplate.SudoUser)
	return diags
}

func resourceScriptTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()

	name := d.Get("name").(string)

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		for _, s := range attr.(*schema.Set).List() {
			labelsPayload = append(labelsPayload, s.(string))
		}
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"containerScript": map[string]interface{}{
				"name":        name,
				"labels":      labelsPayload,
				"scriptType":  d.Get("script_type").(string),
				"scriptPhase": d.Get("script_phase").(string),
				"script":      d.Get("script_content").(string),
				"runAsUser":   d.Get("run_as_user").(string),
				"sudoUser":    d.Get("sudo").(bool),
			},
		},
	}

	resp, err := client.UpdateScriptTemplate(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateScriptTemplateResult)
	scriptTemplate := result.ScriptTemplate
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(scriptTemplate.ID))
	return resourceScriptTemplateRead(ctx, d, meta)
}

func resourceScriptTemplateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteScriptTemplate(toInt64(id), req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
			return nil
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)
			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)
	d.SetId("")
	return diags
}
