package git

import (
	"context"
	"fmt"

	"github.com/abibby/jit/cfg"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

type Github struct {
	client *github.Client
}

func NewGithub(ctx context.Context) *Github {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GetString("github.access_token")},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &Github{
		client: github.NewClient(tc),
	}
}

func (gh *Github) MainBranchName(ctx context.Context) (string, error) {
	u, err := UrlParts()
	if err != nil {
		return "", err
	}

	rep, _, err := gh.client.Repositories.Get(ctx, u.Owner, u.Repo)
	if err != nil {
		return "", err
	}
	if rep.DefaultBranch == nil {
		return "", fmt.Errorf("could not find default branch")
	}

	return *rep.DefaultBranch, nil
}

func (gh *Github) CreatePR(ctx context.Context, opt *PullRequestOptions) (PullRequest, error) {
	u, err := UrlParts()
	if err != nil {
		return nil, err
	}
	pr, _, err := gh.client.PullRequests.Create(
		ctx,
		u.Owner,
		u.Repo,
		&github.NewPullRequest{
			Title: ptr(opt.Title),
			Body:  ptr(opt.Description),
			Head:  ptr(opt.SourceBranch),
			Base:  ptr(opt.BaseBranch),
			// Issue
			// MaintainerCanModify
			Draft: ptr(true),
		},
	)
	if err != nil {
		return nil, err
	}

	return gh.translatePR(pr), nil
}

func (gh *Github) ListPRs(ctx context.Context) ([]PullRequest, error) {
	return []PullRequest{}, nil
}
func (gh *Github) translatePR(ghPR *github.PullRequest) *SimplePullRequest {
	return &SimplePullRequest{
		url: ghPR.GetHTMLURL(),
	}
}
func ptr[T any](v T) *T {
	return &v
}

func (gh *Github) DiffURL(ctx context.Context) (string, error) {
	return "", nil
}
