package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"sort"
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

func allBranches() ([]string, error) {
	result, _, err := gitOutput("branch", "--format", "%(refname)", "--all")
	if err != nil {
		return nil, err
	}

	brancheMap := map[string]struct{}{}
	re := regexp.MustCompile(`^(?:refs\/heads\/|refs\/remotes\/[^\/]+\/)(.*)$`)
	for _, branch := range strings.Split(result, "\n") {
		matches := re.FindStringSubmatch(branch)
		if len(matches) > 1 {
			brancheMap[matches[1]] = struct{}{}
		}
	}
	branches := []string{}
	for branch := range brancheMap {
		branches = append(branches, branch)
	}

	sort.Strings(branches)

	return branches, nil
}

func findBranch(issueID string) (string, error) {
	branches, err := allBranches()
	if err != nil {
		return "", err
	}
	selectedBranches := []string{}
	for _, branch := range branches {
		if strings.Contains(strings.ToLower(branch), strings.ToLower(issueID)) {
			selectedBranches = append(selectedBranches, branch)
		}
	}
	if len(selectedBranches) == 0 {
		return "", fmt.Errorf("could not find branch for jira issue %s", issueID)
	}
	if len(selectedBranches) == 1 {
		return selectedBranches[0], nil
	}
	prompt := promptui.Select{
		Label: "Select Day",
		Items: selectedBranches,
	}

	_, selected, err := prompt.Run()

	if err != nil {
		return "", err
	}

	return selected, nil

}
