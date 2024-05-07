package bitbucket

type Repositories struct {
	c *Client
}

// https://developer.atlassian.com/cloud/bitbucket/rest/api-group-downloads/#api-repositories-workspace-repo-slug-downloads-get
func (pr *Repositories) ListDownloads(workspace, repoSlug string) (*PaginatedResponse[any], error) {
	resp := &PaginatedResponse[any]{}
	err := pr.c.GetJSON(makePath("repositories", workspace, repoSlug, "downloads"), resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
