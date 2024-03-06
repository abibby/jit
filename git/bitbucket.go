package git

import (
	"context"
	"fmt"

	"github.com/abibby/jit/cfg"
	"github.com/abibby/jit/lodash"
	"github.com/ktrysmt/go-bitbucket"
)

type Bitbucket struct {
	client *bitbucket.Client
}

func NewBitbucket(ctx context.Context) *Bitbucket {
	client := bitbucket.NewBasicAuth(
		cfg.GetString("bitbucket.username"),
		cfg.GetString("bitbucket.password"),
	)

	return &Bitbucket{
		client: client,
	}
}

func (bb Bitbucket) MainBranchName(ctx context.Context) (string, error) {
	u, err := UrlParts()
	if err != nil {
		return "", err
	}

	repo, err := bb.client.Repositories.Repository.Get(&bitbucket.RepositoryOptions{
		Owner:    u.Owner,
		RepoSlug: u.Repo,
	})
	if err != nil {
		return "", fmt.Errorf("could not find repo %s: %w", u.String(), err)
	}

	return repo.Mainbranch.Name, nil
}

func (bb Bitbucket) CreatePR(ctx context.Context, opt *PullRequestOptions) (*PullRequest, error) {
	u, err := UrlParts()
	if err != nil {
		return nil, err
	}

	pr, err := bb.client.Repositories.PullRequests.Create(&bitbucket.PullRequestsOptions{
		Message:  opt.Description,
		Owner:    u.Owner,
		RepoSlug: u.Repo,

		Title:             opt.Title,
		Description:       opt.Description,
		SourceBranch:      opt.SourceBranch,
		DestinationBranch: opt.BaseBranch,
	})
	if err != nil {
		if err, ok := err.(*bitbucket.UnexpectedResponseStatusError); ok {
			return nil, fmt.Errorf("could not create pull request: %w", err.ErrorWithBody())
		}
		return nil, fmt.Errorf("could not create pull request: %w", err)
	}

	url, err := lodash.GetString(pr, "links.html.href")
	if err != nil {
		return nil, fmt.Errorf("could not extract the url: %w", err)
	}

	return &PullRequest{
		URL: url,
	}, nil
}
