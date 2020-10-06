package cmd

import (
	"context"

	"github.com/google/go-github/v32/github"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func GitHubClient(ctx context.Context) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: viper.GetString("github.access_token")},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}
