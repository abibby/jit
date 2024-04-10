package bitbucket

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	auth       Authenticator
	httpClient *http.Client
	baseURL    string

	PullRequests *PullRequests
}

func NewClient(auth Authenticator) *Client {
	c := &Client{
		auth:       auth,
		httpClient: http.DefaultClient,
		baseURL:    "https://api.bitbucket.org/2.0/",
	}

	c.PullRequests = &PullRequests{c: c}
	return c
}
func NewBasicAuth(username, password string) *BasicAuth {
	return &BasicAuth{
		username: username,
		password: password,
	}
}

func (c *Client) newRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, c.baseURL+strings.TrimPrefix(url, "/"), body)
	if err != nil {
		return nil, err
	}
	err = c.auth.Authenticate(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (c *Client) do(r *http.Request) (*http.Response, error) {
	return c.httpClient.Do(r)
}

func (c *Client) Get(url string) (*http.Response, error) {
	req, err := c.newRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}
func (c *Client) GetJSON(url string, v any) error {
	httpResp, err := c.Get(url)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	return json.NewDecoder(httpResp.Body).Decode(v)
}

func makePath(parts ...string) string {
	b := &bytes.Buffer{}
	for i, p := range parts {
		if i > 0 {
			b.WriteString("/")
		}
		b.WriteString(url.PathEscape(p))
	}
	return b.String()
}
