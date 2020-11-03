package harness

import (
	"fmt"
)

type ApplicationScope struct {
	AppId      string `json:"appId,omitempty"`
	FilterType string `json:"filterType,omitempty"`
}

type EnvironmentScope struct {
	EnvId      string `json:"envId,omitempty"`
	FilterType string `json:"filterType,omitempty"`
}

type AppEnvScope struct {
	Application *ApplicationScope `json:"application"`
	Environment *EnvironmentScope `json:"environment"`
}

type UsageScope struct {
	AppEnvScopes []*AppEnvScope `json:"appEnvScopes"`
}

type EncryptedSecret struct {
	ID              string      `json:"id"`
	Name            string      `json:"name"`
	Value           string      `json:"value"`
	ScopedToAccount bool        `json:"scopedToAccount"`
	UsageScope      *UsageScope `json:"usageScope"`
}

type EncryptedSecretWrapper struct {
	Secret *EncryptedSecret
}

type GetEncryptedSecretResponse struct {
	Errors []Error
	Data   *EncryptedSecretWrapper
}

type CreateEncryptedSecretWrapper struct {
	CreateSecret *EncryptedSecretWrapper
}

type UpdateEncryptedSecretWrapper struct {
	UpdateSecret *EncryptedSecretWrapper
}

type NewEncryptedSecretResponse struct {
	Errors []Error
	Data   *CreateEncryptedSecretWrapper
}

type DeleteEncryptedSecretApiResponse struct {
	Errors []Error
	Data   *CreateEncryptedSecretWrapper
}

type UpdateEncryptedSecretApiResponse struct {
	Errors []Error
	Data   *UpdateEncryptedSecretWrapper
}

func (h *Client) GetEncryptedSecret(id string) (*EncryptedSecret, error) {
	query := `query { secret(secretId: "%s", secretType: ENCRYPTED_TEXT) { id name } }`
	graphQLQuery := &GraphQLQuery{
		Query: fmt.Sprintf(query, id),
	}

	response := &GetEncryptedSecretResponse{}
	err := h.query(graphQLQuery, response)
	if err != nil {
		return nil, err
	}

	if len(response.Errors) > 0 {
		return nil, fmt.Errorf("Error retrieving secret: %#v", response.Errors)
	}

	return response.Data.Secret, nil
}

func (h *Client) NewEncryptedSecret(s *EncryptedSecret) (*EncryptedSecret, error) {
	query := `mutation($secret: CreateSecretInput!) {
		createSecret(input: $secret){
			secret {
				id
				name
			}
		}
	}
	`

	graphQLQuery := &GraphQLQuery{
		Query: query,
		Variables: map[string]interface{}{
			"secret": map[string]interface{}{
				"secretType": "ENCRYPTED_TEXT",
				"encryptedText": map[string]interface{}{
					"scopedToAccount": s.ScopedToAccount,
					"name":            s.Name,
					"value":           s.Value,
					"usageScope":      s.UsageScope,
				},
			},
		},
	}

	response := &NewEncryptedSecretResponse{}
	err := h.query(graphQLQuery, response)
	if err != nil {
		return nil, err
	}

	if len(response.Errors) > 0 {
		return nil, fmt.Errorf("Error retrieving secret: %#v", response.Errors)
	}

	return response.Data.CreateSecret.Secret, nil
}

func (h *Client) DeleteEncryptedSecret(id string) error {
	fmt.Printf("Deleting a Harness.io secret with id '%s'", id)

	query := `mutation($secret: DeleteSecretInput!){
		deleteSecret(input: $secret) {
			clientMutationId
		}
	}
	`

	graphQlQuery := &GraphQLQuery{
		Query: query,
		Variables: map[string]interface{}{
			"secret": map[string]string{
				"secretId":   id,
				"secretType": "ENCRYPTED_TEXT",
			},
		},
	}

	apiResponse := &DeleteEncryptedSecretApiResponse{}
	err := h.query(graphQlQuery, &apiResponse)

	if err != nil {
		return err
	}

	if len(apiResponse.Errors) > 0 {
		return fmt.Errorf("Errors: %#v", apiResponse.Errors)
	}

	return nil
}

func (h *Client) UpdateEncryptedSecret(s *EncryptedSecret) (*EncryptedSecret, error) {
	fmt.Printf("Updating Harness.io secret with id '%s'", s.ID)

	query := `mutation($secret: UpdateSecretInput!){
		updateSecret(input: $secret) {
			clientMutationId
			secret {
				id
				name
			}
		}
	}
	`

	graphQlQuery := &GraphQLQuery{
		Query: query,
		Variables: map[string]interface{}{
			"secret": map[string]interface{}{
				"secretId":   s.ID,
				"secretType": "ENCRYPTED_TEXT",
				"encryptedText": map[string]interface{}{
					"name":            s.Name,
					"value":           s.Value,
					"scopedToAccount": s.ScopedToAccount,
					"usageScope":      s.UsageScope,
				},
			},
		},
	}

	apiResponse := &UpdateEncryptedSecretApiResponse{}
	err := h.query(graphQlQuery, &apiResponse)

	if err != nil {
		return nil, err
	}

	if len(apiResponse.Errors) > 0 {
		return nil, fmt.Errorf("Errors: %#v", apiResponse.Errors)
	}

	return apiResponse.Data.UpdateSecret.Secret, nil
}
