package azure

import (
	"context"
	"github.com/jenkins-x/go-scm/scm"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
)

type gitService struct {
	client    *wrapper
	gitClient git.Client
}

func (g gitService) FindBranch(ctx context.Context, repo, name string) (*scm.Reference, *scm.Response, error) {
	branch, err := g.gitClient.GetBranch(ctx, git.GetBranchArgs{
		RepositoryId: &repo,
		Name:         &name,
		Project:      &g.client.Project,
	})
	if err != nil {
		return nil, nil, err
	}
	return &scm.Reference{
		Name: *branch.Name,
		Path: scm.ExpandRef(*branch.Name, "refs/heads/"),
		Sha:  *branch.Commit.CommitId,
	}, nil, nil
}

func (g gitService) FindCommit(ctx context.Context, repo, ref string) (*scm.Commit, *scm.Response, error) {
	commit, err := g.gitClient.GetCommit(ctx, git.GetCommitArgs{
		CommitId:     &ref,
		RepositoryId: &repo,
		Project:      &g.client.Project,
	})
	if err != nil {
		return nil, nil, err
	}
	return convertCommit(*commit), nil, nil
}

func (g gitService) FindTag(ctx context.Context, repo, name string) (*scm.Reference, *scm.Response, error) {
	panic("implement me")
}

func (g gitService) ListBranches(ctx context.Context, repo string, opts scm.ListOptions) ([]*scm.Reference, *scm.Response, error) {
	panic("implement me")
}

func (g gitService) ListCommits(ctx context.Context, repo string, opts scm.CommitListOptions) ([]*scm.Commit, *scm.Response, error) {
	panic("implement me")
}

func (g gitService) ListChanges(ctx context.Context, repo, ref string, opts scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	panic("implement me")
}

func (g gitService) CompareCommits(ctx context.Context, repo, ref1, ref2 string, opts scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	panic("implement me")
}

func (g gitService) ListTags(ctx context.Context, repo string, opts scm.ListOptions) ([]*scm.Reference, *scm.Response, error) {
	panic("implement me")
}

func (g gitService) FindRef(ctx context.Context, repo, ref string) (string, *scm.Response, error) {
	commit, err := g.gitClient.GetCommit(ctx, git.GetCommitArgs{
		CommitId:     &ref,
		RepositoryId: &repo,
	})
	if err != nil {
		return "", nil, err
	}
	return *commit.CommitId, nil, nil
}

func (g gitService) DeleteRef(ctx context.Context, repo, ref string) (*scm.Response, error) {
	panic("implement me")
}

func (g gitService) CreateRef(ctx context.Context, repo, ref, sha string) (*scm.Reference, *scm.Response, error) {
	panic("implement me")
}

func convertCommit(from git.GitCommit) *scm.Commit {
	return &scm.Commit{
		Sha:     *from.CommitId,
		Message: *from.Comment,
		Author: scm.Signature{
			Name:   *from.Author.Name,
			Email:  *from.Author.Email,
			Date:   from.Author.Date.Time,
			Avatar: *from.Author.ImageUrl,
		},
		Link: *from.RemoteUrl,
		Committer: scm.Signature{
			Name:   *from.Committer.Name,
			Email:  *from.Committer.Email,
			Date:   from.Committer.Date.Time,
			Avatar: *from.Committer.ImageUrl,
		},
		Tree: scm.CommitTree{
			Sha: *from.TreeId,
		},
	}
}
