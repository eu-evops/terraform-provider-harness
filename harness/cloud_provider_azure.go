package harness

import (
	"fmt"
)

func (h *Client) GetCloudProviderAzure(id string) (*CloudProvider, error) {
	query := `query { cloudProvider(cloudProviderId: "%s") { id name } }`
	graphQLQuery := &GraphQLQuery{
		Query: fmt.Sprintf(query, id),
	}

	response := &GetCloudProviderResponse{}
	err := h.query(graphQLQuery, response)
	if err != nil {
		return nil, err
	}

	if len(response.Errors) > 0 {
		return nil, fmt.Errorf("Error retrieving cloud provider: %#v", response.Errors)
	}

	return response.Data.CloudProvider, nil
}

func (h *Client) NewCloudProviderAzure(name string, secretId string, clientId string, tenantId string) (*CloudProvider, error) {
	query := `mutation($cloudProvider: CreateCloudProviderInput!) {
		createCloudProvider(input: $cloudProvider){
			cloudProvider {
				id
				name
				description
			}
		}
	}
	`

	graphQLQuery := &GraphQLQuery{
		Query: query,
		Variables: map[string]interface{}{
			"cloudProvider": map[string]interface{}{
				"cloudProviderType": "AZURE",
				"azureCloudProvider": map[string]interface{}{
					"name":        name,
					"clientId":    clientId,
					"tenantId":    tenantId,
					"keySecretId": secretId,
				},
			},
		},
	}

	response := &CreateCloudProviderResponse{}
	err := h.query(graphQLQuery, response)
	if err != nil {
		return nil, err
	}

	if len(response.Errors) > 0 {
		return nil, fmt.Errorf("Error creating cloud provider: %#v", response.Errors)
	}

	return response.Data.CreateCloudProvider.CloudProvider, nil
}

func (h *Client) DeleteCloudProviderAzure(id string) error {
	fmt.Printf("Deleting a Harness.io cloud provider with id '%s'", id)

	query := `mutation($cp: DeleteCloudProviderInput!){
		deleteCloudProvider(input: $cp) {
			clientMutationId
		}
	}
	`

	graphQlQuery := &GraphQLQuery{
		Query: query,
		Variables: map[string]interface{}{
			"cp": map[string]string{
				"cloudProviderId": id,
			},
		},
	}

	apiResponse := &DeleteCloudProviderApiResponse{}
	err := h.query(graphQlQuery, &apiResponse)

	if err != nil {
		return err
	}

	if len(apiResponse.Errors) > 0 {
		return fmt.Errorf("Errors: %#v", apiResponse.Errors)
	}

	return nil
}

func (h *Client) UpdateCloudProviderAzure(id string, name string, clientId string, tenantId string) (*CloudProvider, error) {
	{
		query := `mutation($cloudProvider: UpdateCloudProviderInput!) {
			updateCloudProvider(input: $cloudProvider){
				clientMutationId
				cloudProvider {
					id
					name
					description
				}
			}
		}
		`

		graphQLQuery := &GraphQLQuery{
			Query: query,
			Variables: map[string]interface{}{
				"cloudProvider": map[string]interface{}{
					"cloudProviderType": "AZURE",
					"cloudProviderId":   id,
					"azureCloudProvider": map[string]interface{}{
						"name":     name,
						"clientId": clientId,
						"tenantId": tenantId,
					},
				},
			},
		}

		response := &UpdateCloudProviderResponse{}
		err := h.query(graphQLQuery, response)
		if err != nil {
			return nil, err
		}

		if len(response.Errors) > 0 {
			return nil, fmt.Errorf("Error creating cloud provider: %#v", response.Errors)
		}

		return response.Data.UpdateCloudProvider.CloudProvider, nil
	}
}
