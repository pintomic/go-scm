package azure

import (
	"context"
	"fmt"
	"github.com/jenkins-x/go-scm/scm"
	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"github.com/microsoft/azure-devops-go-api/azuredevops/servicehooks"
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
	urlPart := strings.Split(base.Path, "/")
	if len(urlPart) < 3 {
		return nil, fmt.Errorf("azure server url %s /organization/project", uri)
	}
	organization := urlPart[1]
	project := urlPart[2]
	base.Path = organization + "/"
	// Create a connection to your organization
	connection := azuredevops.NewPatConnection(base.String(), token)

	ctx := context.Background()

	client := &wrapper{
		Client:       new(scm.Client),
		Organization: organization,
		Project:      project,
	}
	// Create a client to interact with the Core area
	coreClient, err := core.NewClient(ctx, connection)
	if err != nil {
		return nil, err
	}
	client.AzureClient = &coreClient
	client.BaseURL = base
	client.Driver = scm.DriverAzure
	gitClient, err := git.NewClient(ctx, connection)
	if err != nil {
		return nil, err
	}
	hooksClient := servicehooks.NewClient(ctx, connection)
	client.Repositories = &repositoryService{client, gitClient, hooksClient}
	client.Webhooks = &webhookService{client, hooksClient}
	client.Git = &gitService{client, gitClient}
	client.Reviews = &reviewService{client}

	return client.Client, nil
}

// wrapper wraps the Client to provide high level helper functions
// for making http requests and unmarshalling the response.
type wrapper struct {
	*scm.Client
	AzureClient  *core.Client
	Organization string
	Project      string
}
