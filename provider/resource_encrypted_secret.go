package provider

import (
	"context"
	"log"

	Harness "github.com/eu-evops/terraform-provider-harness/harness"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceEncryptedSecret() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"secret_manager_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"scoped_to_account": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"scope": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"application_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"environment_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"environment_type": {
							Type:        schema.TypeString,
							Description: "Either NON_PRODUCTION_ENVIRONMENTS or PRODUCTION_ENVIRONMENTS",
							Optional:    true,
							ValidateFunc: validation.StringInSlice([]string{
								"NON_PRODUCTION_ENVIRONMENTS",
								"PRODUCTION_ENVIRONMENTS",
							}, false),
						},
					},
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		CreateContext: resourceSecretCreate,
		ReadContext:   resourceSecretRead,
		UpdateContext: resourceSecretUpdate,
		DeleteContext: resourceSecretDelete,
	}
}

func resourceSecretCreate(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Harness.Client)

	secret := &Harness.EncryptedSecret{
		Name:            d.Get("name").(string),
		Value:           d.Get("value").(string),
		SecretManagerID: d.Get("secret_manager_id").(string),
		ScopedToAccount: d.Get("scoped_to_account").(bool),
		UsageScope: &Harness.UsageScope{
			AppEnvScopes: make([]*Harness.AppEnvScope, 0),
		},
	}

	if !secret.ScopedToAccount {
		for _, usageScope := range d.Get("scope").([]interface{}) {
			uu := usageScope.(map[string]interface{})
			u := &Harness.AppEnvScope{
				Application: &Harness.ApplicationScope{
					AppId:      uu["application_id"].(string),
					FilterType: uu["application_type"].(string),
				},
				Environment: &Harness.EnvironmentScope{
					EnvId:      uu["environment_id"].(string),
					FilterType: uu["environment_type"].(string),
				},
			}
			secret.UsageScope.AppEnvScopes = append(secret.UsageScope.AppEnvScopes, u)
		}
	}

	app, err := client.NewEncryptedSecret(secret)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(app.ID)
	d.Set("name", app.Name)

	return resourceSecretRead(c, d, meta)
}

func resourceSecretRead(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Harness.Client)
	app, err := client.GetEncryptedSecret(d.Id())

	_, notFound := err.(*Harness.NotFound)
	if notFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", app.Name)
	d.Set("scoped_to_account", app.ScopedToAccount)

	return nil
}

func resourceSecretUpdate(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Harness.Client)

	secret := &Harness.EncryptedSecret{
		ID:              d.Id(),
		Name:            d.Get("name").(string),
		Value:           d.Get("value").(string),
		ScopedToAccount: d.Get("scoped_to_account").(bool),
		UsageScope: &Harness.UsageScope{
			AppEnvScopes: make([]*Harness.AppEnvScope, 0),
		},
	}

	log.Printf("[DEBUG] Updating secret with id: %s", secret.ID)

	if !secret.ScopedToAccount {
		for _, usageScope := range d.Get("scope").([]interface{}) {
			uu := usageScope.(map[string]interface{})
			u := &Harness.AppEnvScope{
				Application: &Harness.ApplicationScope{
					AppId:      uu["application_id"].(string),
					FilterType: uu["application_type"].(string),
				},
				Environment: &Harness.EnvironmentScope{
					EnvId:      uu["environment_id"].(string),
					FilterType: uu["environment_type"].(string),
				},
			}
			secret.UsageScope.AppEnvScopes = append(secret.UsageScope.AppEnvScopes, u)
		}
	}

	updatedSecret, err := client.UpdateEncryptedSecret(secret)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", updatedSecret.Name)
	d.Set("scoped_to_account", updatedSecret.ScopedToAccount)

	return nil
}

func resourceSecretDelete(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Harness.Client)

	err := client.DeleteEncryptedSecret(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
