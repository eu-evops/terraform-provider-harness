package harness

type CloudProvider struct {
	ID          string
	Name        string
	Description string
}

type CloudProviderWrapper struct {
	CloudProvider *CloudProvider
}

type CreateCloudProviderWrapper struct {
	CreateCloudProvider *CloudProviderWrapper
}

type DeleteCloudProviderWrapper struct {
	DeleteCloudProvider *CloudProviderWrapper
}

type UpdateCloudProviderWrapper struct {
	UpdateCloudProvider *CloudProviderWrapper
}

type CreateCloudProviderResponse struct {
	Errors []Error
	Data   *CreateCloudProviderWrapper
}

type UpdateCloudProviderResponse struct {
	Errors []Error
	Data   *UpdateCloudProviderWrapper
}

type DeleteCloudProviderApiResponse struct {
	Errors []Error
	Data   *DeleteCloudProviderWrapper
}

type GetCloudProviderResponse struct {
	Errors []Error
	Data   *CloudProviderWrapper
}
