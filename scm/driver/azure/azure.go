package azure

import (
	"context"
	"github.com/jenkins-x/go-scm/scm"
	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/core"
	"net/url"
	"strings"
)

func NewWithToken(uri string, token string) (*scm.Client, error) {
	base, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(base.Path, "/") {
		base.Path = base.Path + "/"
	}
	// Create a connection to your organization
	connection := azuredevops.NewPatConnection(uri, token)

	ctx := context.Background()

	client := &wrapper{Client: new(scm.Client)}
	// Create a client to interact with the Core area
	coreClient, err := core.NewClient(ctx, connection)
	if err != nil {
		return nil, err
	}
	client.AzureClient = &coreClient
	client.BaseURL = base
	client.Driver = scm.DriverAzure
	return client.Client, nil
}

// wrapper wraps the Client to provide high level helper functions
// for making http requests and unmarshalling the response.
type wrapper struct {
	*scm.Client
	AzureClient *core.Client
}