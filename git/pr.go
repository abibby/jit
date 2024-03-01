package git

import "context"

type PullRequestOptions struct {
	Title        string
	Description  string
	SourceBranch string
	BaseBranch   string
}
type PullRequest struct {
	URL string
}

func CreatePR(ctx context.Context, opt *PullRequestOptions) (*PullRequest, error) {
	p, err := GetProvider(ctx)
	if err != nil {
		return nil, err
	}
	return p.CreatePR(ctx, opt)
}
