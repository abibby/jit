package cmd

import (
	"fmt"

	"github.com/andygrunwald/go-jira"
	"github.com/spf13/viper"
)

func jiraClient() (*jira.Client, error) {
	tp := jira.BasicAuthTransport{
		Username: viper.GetString("jira.username"),
		Password: viper.GetString("jira.password"),
	}

	return jira.NewClient(tp.Client(), "https://merotechnologies.atlassian.net")
}

func SetStatus(c *jira.Client, issueKey, status string) error {
	transitions, _, err := c.Issue.GetTransitions(issueKey)
	if err != nil {
		return err
	}

	for _, transition := range transitions {
		if transition.Name == status {
			_, err = c.Issue.DoTransition(issueKey, transition.ID)
			if err != nil {
				return err
			}
		}
	}

	return fmt.Errorf("no transition to status %s", status)
}
