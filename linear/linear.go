package linear

import (
	"net/http"
	"os"

	"github.com/Khan/genqlient/graphql"
)

//go:generate curl https://raw.githubusercontent.com/linear/linear/master/packages/sdk/src/schema.graphql -o schema.graphql
//go:generate go run github.com/Khan/genqlient@latest

type Doer func(*http.Request) (*http.Response, error)

func (d Doer) Do(r *http.Request) (*http.Response, error) {
	return d(r)
}

func New() graphql.Client {
	client := http.Client{}
	return graphql.NewClient("https://api.linear.app/graphql", Doer(func(r *http.Request) (*http.Response, error) {
		r.Header.Add("Content-Type", "application/json")
		r.Header.Add("Authorization", os.Getenv("LINEAR_API_KEY"))
		return client.Do(r)
	}))
}
