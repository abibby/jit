package git

import (
	"context"
	"fmt"

	"github.com/abibby/jit/cfg"
	"github.com/abibby/jit/lodash"
	"github.com/abibby/salusa/slices"
	"github.com/ktrysmt/go-bitbucket"

	jitbb "github.com/abibby/jit/bitbucket"
)

type Bitbucket struct {
	client    *bitbucket.Client
	jitClient *jitbb.Client
}

func NewBitbucket(ctx context.Context) *Bitbucket {
	return &Bitbucket{
		client: bitbucket.NewBasicAuth(
			cfg.GetString("bitbucket.username"),
			cfg.GetString("bitbucket.password"),
		),
		jitClient: jitbb.NewClient(jitbb.NewBasicAuth(
			cfg.GetString("bitbucket.username"),
			cfg.GetString("bitbucket.password"),
		)),
	}
}

func (bb *Bitbucket) MainBranchName(ctx context.Context) (string, error) {
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

func (bb *Bitbucket) CreatePR(ctx context.Context, opt *PullRequestOptions) (PullRequest, error) {
	u, err := UrlParts()
	if err != nil {
		return nil, err
	}

	reviewers := make([]string, 0, len(opt.Reviewers))

	members, err := bb.client.Workspaces.Members(u.Owner)
	if err != nil {
		return nil, err
	}
	for _, r := range opt.Reviewers {
		u, ok := slices.Find(members.Members, func(u bitbucket.User) bool {
			return r == u.Uuid || r == u.DisplayName || r == u.Nickname || r == u.Username || r == u.AccountId
		})
		if ok {
			reviewers = append(reviewers, u.Uuid)
		}
	}

	pr, err := bb.client.Repositories.PullRequests.Create(&bitbucket.PullRequestsOptions{
		Message:  opt.Description,
		Owner:    u.Owner,
		RepoSlug: u.Repo,

		Title:             opt.Title,
		Description:       opt.Description,
		SourceBranch:      opt.SourceBranch,
		DestinationBranch: opt.BaseBranch,

		Reviewers: reviewers,
	})
	if err != nil {
		if err, ok := err.(*bitbucket.UnexpectedResponseStatusError); ok {
			return nil, fmt.Errorf("could not create pull request: %w", err.ErrorWithBody())
		}
		return nil, fmt.Errorf("could not create pull request: %w", err)
	}

	return bb.translatePR(pr)
}

func (bb *Bitbucket) ListPRs(ctx context.Context) ([]PullRequest, error) {
	prs, err := bb.jitClient.PullRequests.User("adambibby")
	if err != nil {
		return nil, err
	}
	outPRs := make([]PullRequest, len(prs.Values))
	for i, pr := range prs.Values {
		outPRs[i] = pr
	}
	return outPRs, nil
}
func (bb *Bitbucket) ListUsers(ctx context.Context) (*bitbucket.WorkspaceMembers, error) {
	workspace, err := bb.client.Workspaces.Members("ownersbox")
	if err != nil {
		return nil, err
	}
	return workspace, nil
}

func (bb *Bitbucket) translatePR(bbPR any) (*SimplePullRequest, error) {
	url, err := lodash.Get[string](bbPR, "links.html.href")
	if err != nil {
		return nil, fmt.Errorf("could not extract the url: %w", err)
	}

	return &SimplePullRequest{
		url: url,
	}, nil
}

// func (bb *Bitbucket) translateJitPR(bbPR *jitbb.PullRequest) *PullRequest {
// 	return &PullRequest{
// 		URL:          bbPR.Links["html"].Href,
// 		CommentCount: bbPR.CommentCount,
// 	}
// }

func (bb *Bitbucket) DiffURL(ctx context.Context) (string, error) {
	parts, err := UrlParts()
	if err != nil {
		return "", err
	}

	branch, err := CurrentBranch()
	if err != nil {
		return "", err
	}

	main, err := bb.MainBranchName(ctx)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://bitbucket.org/%s/%s/branch/%s?dest=%s", parts.Owner, parts.Repo, branch, main), nil
}
