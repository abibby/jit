package git

import (
	"context"
	"fmt"
)

type Provider interface {
	MainBranchName(context.Context) (string, error)
	CreatePR(context.Context, *PullRequestOptions) (PullRequest, error)
	DiffURL(ctx context.Context) (string, error)
	ListPRs(ctx context.Context) ([]PullRequest, error)
}

func GetProvider(ctx context.Context) (Provider, error) {
	u, err := UrlParts()
	if err != nil {
		return nil, err
	}
	switch u.Host {
	case "github.com":
		return NewGithub(ctx), nil
	case "bitbucket.org":
		return NewBitbucket(ctx), nil
	default:
		return nil, fmt.Errorf("no git provider for %s", u.Host)
	}
}
