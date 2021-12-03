package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/jenkins-x/go-scm/scm"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"github.com/microsoft/azure-devops-go-api/azuredevops/servicehooks"
	"io/ioutil"
	"testing"
)

type mockGitClient struct {
	git.Client
}

type mockHooksClient struct {
	servicehooks.Client
}

func (m mockGitClient) GetRepository(context.Context, git.GetRepositoryArgs) (*git.GitRepository, error) {
	want := git.GitRepository{}
	raw, err := ioutil.ReadFile("testdata/repo.json")
	json.Unmarshal(raw, &want)
	if err != nil {
		return nil, err
	}
	fmt.Println(want)
	return &want, nil
}

func (m mockGitClient) GetRepositories(context.Context, git.GetRepositoriesArgs) (*[]git.GitRepository, error) {
	var want []git.GitRepository
	raw, err := ioutil.ReadFile("testdata/repos.json")
	json.Unmarshal(raw, &want)
	if err != nil {
		return nil, err
	}
	fmt.Println(want)
	return &want, nil
}

func (m mockGitClient) GetStatuses(context.Context, git.GetStatusesArgs) (*[]git.GitStatus, error) {
	var want []git.GitStatus
	raw, err := ioutil.ReadFile("testdata/statuses.json")
	json.Unmarshal(raw, &want)
	if err != nil {
		return nil, err
	}
	fmt.Println(want)
	return &want, nil
}

func TestRepoFind(t *testing.T) {
	repo := &repositoryService{&wrapper{
		Project: "test-project",
	}, mockGitClient{}, mockHooksClient{}}

	ctx := context.Background()
	got, _, err := repo.Find(ctx, "test-repo")
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Repository)
	raw, _ := ioutil.ReadFile("testdata/repo.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestRepoList(t *testing.T) {
	repo := &repositoryService{&wrapper{
		Project: "test-project",
	}, mockGitClient{}, mockHooksClient{}}

	ctx := context.Background()
	got, _, err := repo.List(ctx, scm.ListOptions{})
	if err != nil {
		t.Error(err)
	}

	var want []*scm.Repository
	raw, _ := ioutil.ReadFile("testdata/repos.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestStatusList(t *testing.T) {
	repo := &repositoryService{&wrapper{
		Project: "test-project",
	}, mockGitClient{}, mockHooksClient{}}

	ctx := context.Background()
	got, _, err := repo.ListStatus(ctx, "test-repo", "6dcb09b5b57875f334f61aebed695e2e4193db5e", scm.ListOptions{})
	if err != nil {
		t.Error(err)
	}

	var want []*scm.Status
	raw, _ := ioutil.ReadFile("testdata/statuses.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}
