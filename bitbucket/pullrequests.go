package bitbucket

type Link struct {
	Href string `json:"href"`
	Name string `json:"name"`
}

type RenderedSection struct {
	Raw    string `json:"raw"`
	Markup string `json:"markup"`
	HTML   string `json:"html"`
}

type Rendered struct {
	Title       *RenderedSection `json:"title"`
	Description *RenderedSection `json:"description"`
	Reason      *RenderedSection `json:"reason"`
}

type PullRequestSourceBranch struct {
	Name string `json:"name"`
}
type PullRequestSourceBranchCommit struct {
	Hash  string           `json:"hash"`
	Links map[string]*Link `json:"links"`
	Type  string           `json:"type"`
}
type PullRequestSourceBranchRepository struct {
	Type     string           `json:"type"`
	FullName string           `json:"full_name"`
	Links    map[string]*Link `json:"links"`
}
type PullRequestSource struct {
	Branch     *PullRequestSourceBranch           `json:"branch"`
	Commit     *PullRequestSourceBranchCommit     `json:"commit"`
	Repository *PullRequestSourceBranchRepository `json:"repository"`
	Name       string                             `json:"name"`
	UUID       string                             `json:"uuid"`
}

type PullRequest struct {
	Type         string             `json:"type"`
	Links        map[string]*Link   `json:"links"`
	ID           int                `json:"id"`
	Title        string             `json:"title"`
	Rendered     *Rendered          `json:"rendered"`
	Summary      *RenderedSection   `json:"summary"`
	State        string             `json:"state"`
	CommentCount int                `json:"comment_count"`
	TaskCount    int                `json:"task_count"`
	Source       *PullRequestSource `json:"source"`
}

type Comment struct{}

func (pr *PullRequest) GetURL() string {
	return pr.Links["html"].Href
}
func (pr *PullRequest) GetCommentCount() int {
	return pr.CommentCount
}

type PullRequests struct {
	c *Client
}

// https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pullrequests/#api-pullrequests-selected-user-get
func (pr *PullRequests) User(username string) (*PaginatedResponse[*PullRequest], error) {
	resp := &PaginatedResponse[*PullRequest]{}
	err := pr.c.GetJSON(makePath("pullrequests", username), resp)
	return resp, err
}

// https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pullrequests/#api-repositories-workspace-repo-slug-pullrequests-pull-request-id-comments-get
func (pr *PullRequests) ListComments(workspace, repoSlug, prID string) (*PaginatedResponse[*Comment], error) {
	resp := &PaginatedResponse[*Comment]{}
	err := pr.c.GetJSON(makePath("pullrequests", workspace, repoSlug, "pullrequests", prID, "comments"), resp)
	return resp, err
}
