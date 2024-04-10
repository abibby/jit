package bitbucket

import (
	"context"
	"net/http"
)

type Authenticator interface {
	Authenticate(ctx context.Context, req *http.Request) error
}

type BasicAuth struct {
	username string
	password string
}

func (b *BasicAuth) Authenticate(ctx context.Context, req *http.Request) error {
	req.SetBasicAuth(b.username, b.password)
	return nil
}
