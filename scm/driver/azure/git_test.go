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
