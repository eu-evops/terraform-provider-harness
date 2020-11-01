package provider

import (
	"context"

	Harness "github.com/eu-evops/terraform-provider-harness/harness"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEncryptedSecret() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
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
	app, err := client.NewEncryptedSecret(d.Get("name").(string), d.Get("value").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(app.ID)
	return resourceSecretRead(c, d, meta)
}

func resourceSecretRead(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Harness.Client)
	app, err := client.GetEncryptedSecret(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", app.Name)

	return nil
}

func resourceSecretUpdate(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Harness.Client)
	app, err := client.UpdateEncryptedSecret(d.Id(), d.Get("name").(string), d.Get("value").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", app.Name)

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
