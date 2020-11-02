package provider

import (
	"context"

	Harness "github.com/eu-evops/terraform-provider-harness/harness"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApplication() *schema.Resource {
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
		},
		CreateContext: resourceApplicationCreate,
		ReadContext:   resourceApplicationRead,
		UpdateContext: resourceApplicationUpdate,
		DeleteContext: resourceApplicationDelete,
	}
}

func resourceApplicationCreate(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Harness.Client)

	app := &Harness.Application{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	app, err := client.NewApplication(app)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(app.ID)
	return resourceApplicationRead(c, d, meta)
}

func resourceApplicationRead(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Harness.Client)
	app, err := client.GetApplication(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", app.Name)
	d.Set("description", app.Description)

	return nil
}

func resourceApplicationUpdate(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Harness.Client)

	app := &Harness.Application{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	app, err := client.UpdateApplication(app)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", app.Name)
	d.Set("description", app.Description)

	return nil
}

func resourceApplicationDelete(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Harness.Client)

	err := client.DeleteApplication(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
