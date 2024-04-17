package pm

import (
	"fmt"
	"strings"

	"github.com/abibby/jit/cfg"
	"github.com/andygrunwald/go-jira"
)

type JiraProvider struct {
	auth    jira.BasicAuthTransport
	baseUrl string
}

func NewJira() *JiraProvider {
	return &JiraProvider{
		auth: jira.BasicAuthTransport{
			Username: cfg.GetString("jira.username"),
			Password: cfg.GetString("jira.password"),
		},
		baseUrl: cfg.GetString("jira.base_url"),
	}
}

func (j *JiraProvider) client() (*jira.Client, error) {
	return jira.NewClient(j.auth.Client(), j.baseUrl)
}
func (j *JiraProvider) GetIssue(issueID string) (*Issue, error) {
	c, err := j.client()
	if err != nil {
		return nil, err
	}

	i, _, err := c.Issue.Get(issueID, nil)
	if err != nil {
		return nil, err
	}

	return transformIssue(i), nil
}
func (j *JiraProvider) GetMyIssues() ([]*Issue, error) {
	c, err := j.client()
	if err != nil {
		return nil, err
	}

	jiraIssues, _, err := c.Issue.Search(`assignee = currentUser() AND statusCategory != "Done"`, nil)
	if err != nil {
		return nil, err
	}

	issues := make([]*Issue, len(jiraIssues))
	for i, ji := range jiraIssues {
		issues[i] = transformIssue(&ji)
	}
	return issues, nil
}

func transformIssue(i *jira.Issue) *Issue {
	return &Issue{
		ID:     i.Key,
		Title:  i.Fields.Summary,
		Status: transformStatus(i.Fields.Status.StatusCategory.Name),
	}
}
func transformStatus(s string) Status {
	switch strings.ToLower(s) {
	case "to do":
		return StatusToDo
	case "in progress":
		return StatusInProgress
	case "in review", "in testing", "ready to deploy":
		return StatusInReview
	case "closed":
		return StatusDone
	}
	return StatusUnknown
}

func (j *JiraProvider) SetStatus(issueKey, status string) error {
	c, err := j.client()
	if err != nil {
		return err
	}
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
