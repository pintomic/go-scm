package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/jenkins-x/go-scm/scm"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"io/ioutil"
	"testing"
)

func (m mockGitClient) GetBranch(context.Context, git.GetBranchArgs) (*git.GitBranchStats, error) {
	want := git.GitBranchStats{}
	raw, err := ioutil.ReadFile("testdata/branch.json")
	json.Unmarshal(raw, &want)
	if err != nil {
		return nil, err
	}
	fmt.Println(want)
	return &want, nil
}

func (m mockGitClient) GetCommit(context.Context, git.GetCommitArgs) (*git.GitCommit, error) {
	want := git.GitCommit{}
	raw, err := ioutil.ReadFile("testdata/commit.json")
	json.Unmarshal(raw, &want)
	if err != nil {
		return nil, err
	}
	fmt.Println(want)
	return &want, nil
}

// Tests
func TestGitFindBranch(t *testing.T) {
	gitService := &gitService{&wrapper{Project: "test-project"}, mockGitClient{}}

	got, _, err := gitService.FindBranch(context.Background(), "test-repo", "master")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Reference)
	raw, _ := ioutil.ReadFile("testdata/branch.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestGitFindCommit(t *testing.T) {
	gitService := &gitService{&wrapper{Project: "test-project"}, mockGitClient{}}

	got, _, err := gitService.FindCommit(context.Background(), "test-repo", "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Commit)
	raw, _ := ioutil.ReadFile("testdata/commit.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}
