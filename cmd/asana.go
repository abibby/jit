package cmd

import (
	"fmt"
	"net/http"

	"github.com/range-labs/go-asana/asana"
)

type User struct {
	GID        string            `json:"gid,omitempty"`
	Email      string            `json:"email,omitempty"`
	Name       string            `json:"name,omitempty"`
	Photo      map[string]string `json:"photo,omitempty"`
	Workspaces []Workspace       `json:"workspaces,omitempty"`
}

type Workspace struct {
	GID          int64  `json:"gid,omitempty"`
	Name         string `json:"name,omitempty"`
	Organization bool   `json:"is_organization,omitempty"`
}

func asanaClient() *asana.Client {
	return asana.NewClient(asana.DoerFunc(func(req *http.Request) (resp *http.Response, err error) {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", configGetString("asana.access_token")))
		return http.DefaultClient.Do(req)
	}))
}
