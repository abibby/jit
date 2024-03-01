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

func (gh Github) MainBranchName(ctx context.Context) (string, error) {
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
