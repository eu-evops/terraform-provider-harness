package harness

import (
	"fmt"
	"strings"
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
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ApplicationRetrievalError struct{}

func (e *ApplicationRetrievalError) Error() string {
	return fmt.Sprintf("ApplicationRetrievalError")
}

type UserNotAuthorisedError struct{}

func (e *UserNotAuthorisedError) Error() string {
	return fmt.Sprintf("UserNotAuthorisedError")
}

type ApplicationNotFound struct{}

func (e *ApplicationNotFound) Error() string {
	return fmt.Sprintf("ApplicationNotFound")
}

func (h *Client) GetApplication(id string) (*Application, error) {
	fmt.Print("Getting a Harness.io application with id '%s'", id)

	query := `query {
		application(applicationId: "%s"){
			id
			name
			description
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

		if strings.Contains(apiResponse.Errors[0].Message, "User not authorized") {
			return nil, &UserNotAuthorisedError{}
		}

		return nil, fmt.Errorf("Errors: %#v", apiResponse.Errors)
	}

	return apiResponse.Data.Application, nil
}

func (h *Client) GetApplicationByName(name string) (*Application, error) {
	fmt.Print("Getting a Harness.io application with name '%s'", name)

	query := `query {
		applicationByName(name: "%s"){
			id
			name
			description
		}
	}
	`

	graphQlQuery := &GraphQLQuery{
		Query: fmt.Sprintf(query, name),
	}
	apiResponse := &GetApplicationApiResponse{}
	err := h.query(graphQlQuery, &apiResponse)
	if err != nil {
		return nil, err
	}

	if len(apiResponse.Errors) > 0 {
		if strings.Contains(apiResponse.Errors[0].Message, "Application does not exist") {
			return nil, &ApplicationNotFound{}
		}

		return nil, fmt.Errorf("Errors: %#v", apiResponse.Errors)
	}

	return apiResponse.Data.Application, nil
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

func (h *Client) NewApplication(a *Application) (*Application, error) {
	fmt.Print("Creating a Harness.io application with name '%s'", a.Name)

	query := `mutation createApp($app: CreateApplicationInput!){
		createApplication(input: $app){
			application {
				id
				name
				description
			}
		}
	}
	`

	graphQlQuery := &GraphQLQuery{
		OperationName: "createApp",
		Query:         query,
		Variables: map[string]interface{}{
			"app": map[string]string{
				"name":        a.Name,
				"description": a.Description,
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

	return apiResponse.Data.CreateApplication.Application, nil
}

func (h *Client) UpdateApplication(a *Application) (*Application, error) {
	fmt.Print("Updating a Harness.io application with id '%s'", a.ID)

	query := `mutation updateApp($app: UpdateApplicationInput!){
		updateApplication(input: $app){
			application {
				id
				name
				description
			}
		}
	}
	`

	graphQlQuery := &GraphQLQuery{
		OperationName: "updateApp",
		Query:         query,
		Variables: map[string]interface{}{
			"app": map[string]string{
				"applicationId": a.ID,
				"name":          a.Name,
				"description":   a.Description,
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
