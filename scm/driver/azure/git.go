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
	return convertBranchStats(*branch), nil, nil
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
	return nil, nil, scm.ErrNotSupported
}

func (g gitService) ListBranches(ctx context.Context, repo string, opts scm.ListOptions) ([]*scm.Reference, *scm.Response, error) {
	branches, err := g.gitClient.GetBranches(ctx, git.GetBranchesArgs{
		RepositoryId: &repo,
		Project:      &g.client.Project,
	})
	if err != nil {
		return nil, nil, err
	}
	return convertBranchStatsList(branches), nil, nil
}

func (g gitService) ListCommits(ctx context.Context, repo string, opts scm.CommitListOptions) ([]*scm.Commit, *scm.Response, error) {
	itemVersion := ""
	if opts.Ref != "" {
		itemVersion = opts.Ref
	}
	if opts.Sha != "" {
		itemVersion = opts.Sha
	}
	commits, err := g.gitClient.GetCommits(ctx, git.GetCommitsArgs{
		RepositoryId: &repo,
		Project:      &g.client.Project,
		SearchCriteria: &git.GitQueryCommitsCriteria{
			ItemVersion: &git.GitVersionDescriptor{
				Version: &itemVersion,
			},
		},
	})
	if err != nil {
		return nil, nil, err
	}
	return convertCommitList(commits), nil, nil
}

func (g gitService) ListChanges(ctx context.Context, repo, ref string, opts scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (g gitService) CompareCommits(ctx context.Context, repo, ref1, ref2 string, opts scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (g gitService) ListTags(ctx context.Context, repo string, opts scm.ListOptions) ([]*scm.Reference, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
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
	return nil, scm.ErrNotSupported
}

func (g gitService) CreateRef(ctx context.Context, repo, ref, sha string) (*scm.Reference, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func convertCommitList(src *[]git.GitCommitRef) []*scm.Commit {
	var dst []*scm.Commit
	for _, v := range *src {
		dst = append(dst, convertCommitRef(v))
	}
	return dst
}

func convertCommitRef(from git.GitCommitRef) *scm.Commit {
	return &scm.Commit{
		Sha:     *from.CommitId,
		Message: *from.Comment,
		Author: scm.Signature{
			Name:  *from.Author.Name,
			Email: *from.Author.Email,
			Date:  from.Author.Date.Time,
		},
		Link: *from.RemoteUrl,
		Committer: scm.Signature{
			Name:  *from.Committer.Name,
			Email: *from.Committer.Email,
			Date:  from.Committer.Date.Time,
		},
	}
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

func convertBranchStatsList(src *[]git.GitBranchStats) []*scm.Reference {
	var dst []*scm.Reference
	for _, v := range *src {
		dst = append(dst, convertBranchStats(v))
	}
	return dst
}

func convertBranchStats(from git.GitBranchStats) *scm.Reference {
	return &scm.Reference{
		Name: *from.Name,
		Path: scm.ExpandRef(*from.Name, "refs/heads/"),
		Sha:  *from.Commit.CommitId,
	}
}