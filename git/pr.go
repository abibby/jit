package git

import "context"

type PullRequestOptions struct {
	Title        string
	Description  string
	SourceBranch string
	BaseBranch   string
}
type PullRequest interface {
	GetURL() string
	GetCommentCount() int
}
type SimplePullRequest struct {
	url          string
	commentCount int
}

func (pr *SimplePullRequest) GetURL() string {
	return pr.url
}
func (pr *SimplePullRequest) GetCommentCount() int {
	return pr.commentCount
}

func CreatePR(ctx context.Context, opt *PullRequestOptions) (PullRequest, error) {
	p, err := GetProvider(ctx)
	if err != nil {
		return nil, err
	}
	return p.CreatePR(ctx, opt)
}
