package jirahelper

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/abibby/jit/cfg"
	"github.com/abibby/jit/git"
	"github.com/andygrunwald/go-jira"
)

func NewClient() (*jira.Client, error) {
	tp := jira.BasicAuthTransport{
		Username: cfg.GetString("jira.username"),
		Password: cfg.GetString("jira.password"),
	}

	return jira.NewClient(tp.Client(), cfg.GetString("jira.base_url"))
}

func SetStatus(c *jira.Client, issueKey, status string) error {
	transitions, _, err := c.Issue.GetTransitions(issueKey)
	if err != nil {
		return err
	}

	for _, transition := range transitions {
		if transition.To.Name == status {
			_, err = c.Issue.DoTransition(issueKey, transition.ID)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return fmt.Errorf("no transition to status %s", status)
}

func GetIssueID() (string, error) {
	branch, err := git.CurrentBranch()
	if err != nil {
		return "", err
	}

	matches := regexp.MustCompile(`[A-Za-z]{2,}-\d+`).FindStringSubmatch(branch)
	if matches == nil {
		return "", nil
	}
	return strings.ToUpper(matches[0]), err
}
