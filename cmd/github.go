package cmd

import (
	"context"
	"fmt"
	"regexp"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

func GitHubClient(ctx context.Context) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: configGetString("github.access_token")},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func masterBranchName(ctx context.Context) (string, error) {
	gh := GitHubClient(ctx)
	owner, repo, err := ownerAndRepo()
	if err != nil {
		return "", err
	}

	rep, _, err := gh.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return "", err
	}
	if rep.DefaultBranch == nil {
		return "", fmt.Errorf("could not find default branch")
	}

	return *rep.DefaultBranch, nil
}

func ownerAndRepo() (string, string, error) {
	url, _, err := gitOutput("remote", "get-url", "origin")
	if err != nil {
		return "", "", err
	}
	re := regexp.MustCompile(`(?:https?:\/\/github\.com\/|git@github\.com[:/])([^\/]+)\/(.+)\.git`)
	matches := re.FindStringSubmatch(url)
	if len(matches) <= 2 {
		return "", "", fmt.Errorf("not a github repo")
	}
	return matches[1], matches[2], nil
}
