package azure

import (
	"code.gitea.io/sdk/gitea"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jenkins-x/go-scm/scm"
	"github.com/microsoft/azure-devops-go-api/azuredevops/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"github.com/microsoft/azure-devops-go-api/azuredevops/servicehooks"
)

/*
 * https://docs.microsoft.com/en-us/rest/api/azure/devops/git/?view=azure-devops-rest-5.1
 * https://docs.microsoft.com/en-us/rest/api/azure/devops/git/?view=azure-devops-rest-6.1
 */

type repositoryService struct {
	client      *wrapper
	gitClient   git.Client
	hooksClient servicehooks.Client
}

// Find necessary
func (s *repositoryService) Find(ctx context.Context, repo string) (*scm.Repository, *scm.Response, error) {
	out, err := s.gitClient.GetRepository(ctx, git.GetRepositoryArgs{
		RepositoryId: &repo,
		Project:      &s.client.Project,
	})
	if err != nil {
		return nil, nil, err
	}
	return convertRepository(out), nil, nil
}

func (s *repositoryService) FindHook(ctx context.Context, repo string, id string) (*scm.Hook, *scm.Response, error) {
	gitRepo, err := s.gitClient.GetRepository(ctx, git.GetRepositoryArgs{
		RepositoryId: &repo,
		Project:      &s.client.Project,
	})
	if err != nil {
		return nil, nil, err
	}
	subs, err := s.hooksClient.ListSubscriptions(ctx, servicehooks.ListSubscriptionsArgs{})
	if err != nil {
		return nil, nil, err
	}

	repoId := gitRepo.Id.String()
	for _, subscription := range *subs {
		if (*subscription.PublisherInputs)["repository"] == repoId && subscription.Id.String() == id {
			return convertSubscription(subscription), nil, nil
		}
	}
	return nil, nil, fmt.Errorf("Couldn't find servicehook")
}

func (s *repositoryService) FindPerms(ctx context.Context, repo string) (*scm.Perm, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *repositoryService) ListOrganisation(ctx context.Context, org string, opts scm.ListOptions) ([]*scm.Repository, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *repositoryService) ListUser(ctx context.Context, username string, opts scm.ListOptions) ([]*scm.Repository, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

// ListLabels necessary
func (s *repositoryService) ListLabels(ctx context.Context, repo string, opts scm.ListOptions) ([]*scm.Label, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *repositoryService) ListHooks(ctx context.Context, repo string, opts scm.ListOptions) ([]*scm.Hook, *scm.Response, error) {
	gitRepo, err := s.gitClient.GetRepository(ctx, git.GetRepositoryArgs{
		RepositoryId: &repo,
		Project:      &s.client.Project,
	})
	if err != nil {
		return nil, nil, err
	}
	subs, err := s.hooksClient.ListSubscriptions(ctx, servicehooks.ListSubscriptionsArgs{})
	if err != nil {
		return nil, nil, err
	}

	repoId := gitRepo.Id.String()
	var dst []*scm.Hook
	for _, subscription := range *subs {
		if (*subscription.PublisherInputs)["repository"] == repoId {
			dst = append(dst, convertSubscription(subscription))
		}
	}
	return dst, nil, nil
}

// ListStatus necessary
func (s *repositoryService) ListStatus(ctx context.Context, repo string, ref string, _ scm.ListOptions) ([]*scm.Status, *scm.Response, error) {
	latestOnly := false
	statuses, err := s.gitClient.GetStatuses(ctx, git.GetStatusesArgs{
		CommitId:     &ref,
		RepositoryId: &repo,
		Project:      &s.client.Project,
		LatestOnly:   &latestOnly,
	})
	if err != nil {
		return nil, nil, err
	}
	return convertStatusList(statuses), nil, nil
}

// CreateStatus necessary
func (s *repositoryService) CreateStatus(ctx context.Context, repo string, ref string, input *scm.StatusInput) (*scm.Status, *scm.Response, error) {
	gitStatus := convertFromState(input.State)
	genre := "go-scm"
	status, err := s.gitClient.CreateCommitStatus(ctx, git.CreateCommitStatusArgs{
		GitCommitStatusToCreate: &git.GitStatus{
			State:       &gitStatus,
			Description: &input.Desc,
			TargetUrl:   &input.Target,
			Context: &git.GitStatusContext{
				Genre: &genre,
				Name:  &input.Label,
			},
		},
		CommitId:     &ref,
		RepositoryId: &repo,
	})
	if err != nil {
		return nil, nil, err
	}
	return convertStatus(*status), nil, nil
}

// FindCombinedStatus necessary
func (s *repositoryService) FindCombinedStatus(ctx context.Context, repo, ref string) (*scm.CombinedStatus, *scm.Response, error) {
	latestOnly := false
	gitStatuses, err := s.gitClient.GetStatuses(ctx, git.GetStatusesArgs{
		CommitId:     &ref,
		RepositoryId: &repo,
		Project:      &s.client.Project,
		LatestOnly:   &latestOnly,
	})
	if err != nil {
		return nil, nil, err
	}
	statuses := convertStatusList(gitStatuses)
	latest := scm.StateUnknown
	if len(statuses) > 0 {
		latest = statuses[0].State
	}
	return &scm.CombinedStatus{
		State:    latest,
		Sha:      ref,
		Statuses: statuses,
	}, nil, nil
}

func (s *repositoryService) Create(ctx context.Context, input *scm.RepositoryInput) (*scm.Repository, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *repositoryService) Fork(ctx context.Context, input *scm.RepositoryInput, origRepo string) (*scm.Repository, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

// CreateHook necessary
func (s *repositoryService) CreateHook(ctx context.Context, repo string, input *scm.HookInput) (*scm.Hook, *scm.Response, error) {
	panic("implement me")

}

func (s *repositoryService) UpdateHook(ctx context.Context, repo string, input *scm.HookInput) (*scm.Hook, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

// DeleteHook necessary
func (s *repositoryService) DeleteHook(ctx context.Context, repo string, id string) (*scm.Response, error) {
	panic("implement me")
}

// IsCollaborator necessary
func (s *repositoryService) IsCollaborator(ctx context.Context, repo string, user string) (bool, *scm.Response, error) {
	return false, nil, scm.ErrNotSupported
}

// AddCollaborator necessary
func (s *repositoryService) AddCollaborator(ctx context.Context, repo, user, permission string) (bool, bool, *scm.Response, error) {
	return false, false, nil, scm.ErrNotSupported
}

// ListCollaborators necessary
func (s *repositoryService) ListCollaborators(ctx context.Context, repo string, ops scm.ListOptions) ([]scm.User, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

// FindUserPermission necessary
func (s *repositoryService) FindUserPermission(ctx context.Context, repo string, user string) (string, *scm.Response, error) {
	return "", nil, scm.ErrNotSupported
}

func (s *repositoryService) Delete(ctx context.Context, repo string) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *repositoryService) List(ctx context.Context, _ scm.ListOptions) ([]*scm.Repository, *scm.Response, error) {

	repos, err := s.gitClient.GetRepositories(ctx, git.GetRepositoriesArgs{
		Project: &s.client.Project,
	})
	if err == nil {
		return convertRepositoryList(repos), nil, err
	}
	return nil, nil, err
}

//
// native data structure conversion
//

func convertRepositoryList(src *[]git.GitRepository) []*scm.Repository {
	var dst []*scm.Repository
	for _, v := range *src {
		dst = append(dst, convertRepository(&v))
	}
	return dst
}

func convertRepository(src *git.GitRepository) *scm.Repository {
	if src == nil {
		return nil
	}
	var defaultBranch = ""
	if src.DefaultBranch != nil {
		defaultBranch = *src.DefaultBranch
	}
	return &scm.Repository{
		ID:        src.Id.String(),
		Namespace: *src.Project.Name,
		Name:      *src.Name,
		FullName:  *src.Name,
		Perm: &scm.Perm{ // TODO: get permissions
			Push:  true,
			Pull:  true,
			Admin: true,
		},
		Branch:   defaultBranch,
		Private:  *src.Project.Visibility == core.ProjectVisibilityValues.Private,
		Clone:    *src.RemoteUrl,
		CloneSSH: *src.SshUrl,
		Link:     *src.WebUrl,
		Created:  src.Project.LastUpdateTime.Time, // TODO: find created time from repo
		Updated:  src.Project.LastUpdateTime.Time,
	}
}

func convertPerm(src *gitea.Permission) *scm.Perm {
	if src == nil {
		return nil
	}
	return &scm.Perm{
		Push:  src.Push,
		Pull:  src.Pull,
		Admin: src.Admin,
	}
}

func convertStatusList(src *[]git.GitStatus) []*scm.Status {
	var dst []*scm.Status
	for _, v := range *src {
		dst = append(dst, convertStatus(v))
	}
	return dst
}

func convertStatus(from git.GitStatus) *scm.Status {
	return &scm.Status{
		State:  convertState(*from.State),
		Label:  *from.Context.Name,
		Desc:   *from.Description,
		Target: *from.TargetUrl,
	}
}

func convertState(from git.GitStatusState) scm.State {
	switch from {
	case git.GitStatusStateValues.Error:
		return scm.StateError
	case git.GitStatusStateValues.Failed:
		return scm.StateFailure
	case git.GitStatusStateValues.Pending:
		return scm.StatePending
	case git.GitStatusStateValues.Succeeded:
		return scm.StateSuccess
	default:
		return scm.StateUnknown
	}
}

func convertFromState(from scm.State) git.GitStatusState {
	switch from {
	case scm.StatePending, scm.StateRunning:
		return git.GitStatusStateValues.Pending
	case scm.StateSuccess:
		return git.GitStatusStateValues.Succeeded
	case scm.StateFailure:
		return git.GitStatusStateValues.Failed
	default:
		return git.GitStatusStateValues.Error
	}
}

func convertSubscription(from servicehooks.Subscription) *scm.Hook {
	return &scm.Hook{
		ID:         from.Id.String(),
		Name:       *from.EventDescription,
		Target:     (*from.ConsumerInputs)["url"],
		Events:     []string{*from.EventType},
		Active:     true,
		SkipVerify: false,
	}
}
