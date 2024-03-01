package git

import (
	"context"
	"fmt"

	"github.com/abibby/jit/cfg"
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
