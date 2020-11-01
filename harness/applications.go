package harness

import (
	"fmt"
)

type CreateApplicationWrapper struct {
	Application *Application `json:"application"`
}

type CreateApplicationOperationWrapper struct {
	CreateApplication *CreateApplicationWrapper `json:"createApplication"`
}

type CreateApplicationApiResponse struct {
	Errors []Error                           `json:"errors"`
	Data   CreateApplicationOperationWrapper `json:"data"`
}

type UpdateApplicationOperationWrapper struct {
	UpdateApplication *ApplicationWrapper `json:"updateApplication"`
}

type UpdateApplicationApiResponse struct {
	Errors []Error                           `json:"errors"`
	Data   UpdateApplicationOperationWrapper `json:"data"`
}

type ApplicationWrapper struct {
	Application *Application `json:"application"`
}

type GetApplicationApiResponse struct {
	Errors []Error             `json:"errors"`
	Data   *ApplicationWrapper `json:"data"`
}

type DeleteApplicationApiResponse struct {
	Errors []Error             `json:"errors"`
	Data   *ApplicationWrapper `json:"data"`
}

type Application struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Client) GetApplication(id string) (*Application, error) {
	fmt.Print("Getting a Harness.io application with id '%s'", id)

	query := `query {
		application(applicationId: "%s"){
			id
			name
		}
	}
	`

	graphQlQuery := &GraphQLQuery{
		Query: fmt.Sprintf(query, id),
	}
	apiResponse := &GetApplicationApiResponse{}
	err := h.query(graphQlQuery, &apiResponse)
	if err != nil {
		return nil, err
	}

	if len(apiResponse.Errors) > 0 {
		return nil, fmt.Errorf("Errors: %#v", apiResponse.Errors)
	}

	app := &Application{}
	app.ID = apiResponse.Data.Application.ID
	app.Name = apiResponse.Data.Application.Name

	return app, nil
}

func (h *Client) DeleteApplication(id string) error {
	fmt.Print("Deleting a Harness.io application with id '%s'", id)

	query := `mutation($app: DeleteApplicationInput!){
		deleteApplication(input: $app) {
			clientMutationId
		}
	}
	`

	graphQlQuery := &GraphQLQuery{
		Query: query,
		Variables: map[string]interface{}{
			"app": map[string]string{
				"applicationId": id,
			},
		},
	}

	apiResponse := &DeleteApplicationApiResponse{}
	err := h.query(graphQlQuery, &apiResponse)

	if err != nil {
		return err
	}

	if len(apiResponse.Errors) > 0 {
		return fmt.Errorf("Errors: %#v", apiResponse.Errors)
	}

	return nil
}

func (h *Client) NewApplication(name string) (*Application, error) {
	fmt.Print("Creating a Harness.io application with name '%s'", name)

	query := `mutation createApp($app: CreateApplicationInput!){
		createApplication(input: $app){
			application {
				id
				name
			}
		}
	}
	`

	graphQlQuery := &GraphQLQuery{
		OperationName: "createApp",
		Query:         query,
		Variables: map[string]interface{}{
			"app": map[string]string{
				"name": name,
			},
		},
	}

	apiResponse := &CreateApplicationApiResponse{}
	err := h.query(graphQlQuery, &apiResponse)

	if err != nil {
		return nil, err
	}

	if len(apiResponse.Errors) > 0 {
		return nil, fmt.Errorf("Errors: %#v", apiResponse.Errors)
	}

	app := &Application{}
	app.ID = apiResponse.Data.CreateApplication.Application.ID
	app.Name = apiResponse.Data.CreateApplication.Application.Name

	return app, nil
}

func (h *Client) UpdateApplication(id string, name string) (*Application, error) {
	fmt.Print("Updating a Harness.io application with id '%s'", id)

	query := `mutation updateApp($app: UpdateApplicationInput!){
		updateApplication(input: $app){
			application {
				id
				name
			}
		}
	}
	`

	graphQlQuery := &GraphQLQuery{
		OperationName: "updateApp",
		Query:         query,
		Variables: map[string]interface{}{
			"app": map[string]string{
				"applicationId": id,
				"name":          name,
			},
		},
	}

	apiResponse := &UpdateApplicationApiResponse{}
	err := h.query(graphQlQuery, &apiResponse)
	if err != nil {
		return nil, err
	}

	if len(apiResponse.Errors) > 0 {
		return nil, fmt.Errorf("Errors: %#v", apiResponse.Errors)
	}

	app := &Application{}
	app.ID = apiResponse.Data.UpdateApplication.Application.ID
	app.Name = apiResponse.Data.UpdateApplication.Application.Name

	return app, nil
}
