package provider

import (
	"context"

	Harness "github.com/eu-evops/terraform-provider-harness/harness"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudProviderKubernetes() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"token_secret_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		CreateContext: resourceCloudProviderKubernetesCreate,
		ReadContext:   resourceCloudProviderKubernetesRead,
		UpdateContext: resourceCloudProviderKubernetesUpdate,
		DeleteContext: resourceCloudProviderKubernetesDelete,
	}
}

func resourceCloudProviderKubernetesCreate(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Harness.Client)
	app, err := client.NewCloudProviderKubernetes(
		d.Get("name").(string),
		d.Get("token_secret_id").(string),
		d.Get("url").(string),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(app.ID)
	return resourceCloudProviderKubernetesRead(c, d, meta)
}

func resourceCloudProviderKubernetesRead(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Harness.Client)
	app, err := client.GetCloudProviderKubernetes(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", app.Name)

	return nil
}

func resourceCloudProviderKubernetesUpdate(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Harness.Client)
	app, err := client.UpdateCloudProviderKubernetes(
		d.Id(),
		d.Get("name").(string),
		d.Get("url").(string),
		d.Get("token_secret_id").(string),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", app.Name)

	return nil
}

func resourceCloudProviderKubernetesDelete(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Harness.Client)

	err := client.DeleteCloudProviderKubernetes(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
