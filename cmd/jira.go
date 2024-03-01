package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/andygrunwald/go-jira"
)

func jiraClient() (*jira.Client, error) {
	tp := jira.BasicAuthTransport{
		Username: configGetString("jira.username"),
		Password: configGetString("jira.password"),
	}

	return jira.NewClient(tp.Client(), "https://merotechnologies.atlassian.net")
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

func getIssueTag() (string, error) {
	branch, err := currentBranch()
	if err != nil {
		return "", err
	}

	matches := regexp.MustCompile(`[A-Za-z]{2,}-\d+`).FindStringSubmatch(branch)
	if matches == nil {
		return "", nil
	}
	return strings.ToUpper(matches[0]), err
}
