package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	Harness "github.com/eu-evops/terraform-provider-harness/harness"
)

// Provider for Harness.io
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Description: "Your Harness.io API key",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HARNESS_API_KEY", nil),
			},
			"account_id": {
				Type:        schema.TypeString,
				Description: "Your Harness.io account id",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HARNESS_ACCOUNT_ID", nil),
			},
			"endpoint": {
				Type:        schema.TypeString,
				Description: "Your Harness.io endpoint",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HARNESS_ENDPOINT", nil),
				Default:     "https://app.harness.io",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"harness_application":               resourceApplication(),
			"harness_cloud_provider_azure":      resourceCloudProviderAzure(),
			"harness_cloud_provider_kubernetes": resourceCloudProviderKubernetes(),
			"harness_encrypted_secret":          resourceEncryptedSecret(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: configureFunc,
	}
}

func configureFunc(c context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiKey := d.Get("api_key").(string)
	endpoint := d.Get("endpoint").(string)
	accountID := d.Get("account_id").(string)

	url := fmt.Sprintf("%s/gateway/api/graphql?accountId=%s", endpoint, accountID)

	return Harness.NewClient(apiKey, url), nil
}
