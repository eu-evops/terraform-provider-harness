package provider

import (
	"context"

	Harness "github.com/eu-evops/terraform-provider-harness/harness"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudProviderAzure() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"encrypted_secret_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		CreateContext: resourceCloudProviderAzureCreate,
		ReadContext:   resourceCloudProviderAzureRead,
		UpdateContext: resourceCloudProviderAzureUpdate,
		DeleteContext: resourceCloudProviderAzureDelete,
	}
}

func resourceCloudProviderAzureCreate(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Harness.Client)
	app, err := client.NewCloudProviderAzure(
		d.Get("name").(string),
		d.Get("encrypted_secret_id").(string),
		d.Get("client_id").(string),
		d.Get("tenant_id").(string),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(app.ID)
	return resourceCloudProviderAzureRead(c, d, meta)
}

func resourceCloudProviderAzureRead(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Harness.Client)
	app, err := client.GetCloudProviderAzure(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", app.Name)

	return nil
}

func resourceCloudProviderAzureUpdate(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Harness.Client)
	app, err := client.UpdateCloudProviderAzure(
		d.Id(),
		d.Get("name").(string),
		d.Get("client_id").(string),
		d.Get("tenant_id").(string),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", app.Name)

	return nil
}

func resourceCloudProviderAzureDelete(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Harness.Client)

	err := client.DeleteCloudProviderAzure(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
