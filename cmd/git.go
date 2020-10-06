package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
)

func git(options ...string) error {
	fmt.Printf("git %s\n", strings.Join(options, " "))
	return gitRaw(os.Stdout, os.Stderr, options...)
}

func gitOutput(options ...string) (string, string, error) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	err := gitRaw(stdout, stderr, options...)
	if err != nil {
		return stdout.String(), stderr.String(), err
	}
	return stdout.String(), stderr.String(), nil
}

func gitRaw(stdout, stderr io.Writer, options ...string) error {
	cmd := exec.Command("git", options...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}

func branchName(issue *jira.Issue, message string) string {
	if message == "" {
		message = issue.Fields.Summary
	}
	return viper.GetString("branch_prefix") + prepBranchName(issue.Key+" "+message)
}

func findBranch(issueID string) (string, error) {
	result, _, err := gitOutput("branch", "--format", "%(refname)")
	if err != nil {
		return "", err
	}

	branches := []string{}
	branchPrefix := viper.GetString("branch_prefix") + prepBranchName(issueID)
	ref := "refs/heads/"
	for _, branch := range strings.Split(result, "\n") {
		if strings.HasPrefix(branch, ref+branchPrefix) {
			branches = append(branches, branch[len(ref):])
		}
	}
	if len(branches) == 0 {
		return "", fmt.Errorf("could not find branch for jira issue %s", issueID)
	}
	if len(branches) == 1 {
		return branches[0], nil
	}
	prompt := promptui.Select{
		Label: "Select Day",
		Items: branches,
	}

	_, selected, err := prompt.Run()

	if err != nil {
		return "", err
	}

	return selected, nil

}
