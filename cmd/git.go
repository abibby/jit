package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"bitbucket.org/zombiezen/cardcpx/natsort"
	"github.com/andygrunwald/go-jira"
	"github.com/manifoldco/promptui"
	"golang.org/x/exp/constraints"
)

func git(options ...string) error {
	fmt.Printf("git %s\n", strings.Join(options, " "))
	return execRaw("git", os.Stdout, os.Stderr, options...)
}

func gitRoot() (string, error) {
	return "./.git", nil
}

func gitOutput(options ...string) (string, string, error) {
	return execOutput("git", options...)
}
func execOutput(command string, options ...string) (string, string, error) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	err := execRaw(command, stdout, stderr, options...)
	if err != nil {
		return stdout.String(), stderr.String(), err
	}
	return stdout.String(), stderr.String(), nil
}

func execRaw(command string, stdout, stderr io.Writer, options ...string) error {
	cmd := exec.Command(command, options...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}

func branchName(issue *jira.Issue, message string) string {
	if message == "" {
		message = issue.Fields.Summary
	}
	return configGetString("branch_prefix") + prepBranchName(issue.Key+" "+message)
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

	natsort.Strings(branches)

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
		Label: "Select branch",
		Items: selectedBranches,
	}

	_, selected, err := prompt.Run()

	if err != nil {
		return "", err
	}

	return selected, nil

}

func anyHasPrefix(branch string, prefixs ...string) bool {
	for _, prefix := range prefixs {
		if strings.HasPrefix(strings.ToLower(branch), strings.ToLower(prefix)) {
			return true
		}
	}
	return false
}

func defaultBranch(ctx context.Context) (string, error) {
	masterBranch, err := masterBranchName(ctx)
	if err != nil {
		return "", err
	}

	branches, err := allBranches()
	if err != nil {
		return "", err
	}

	releaseBranches := []string{}
	otherBranches := []string{}

	for _, branch := range reverseStringSlice(branches) {
		if anyHasPrefix(branch, "release/") {
			releaseBranches = append(releaseBranches, branch)
		}
		if anyHasPrefix(branch, "feature/", "epic/") {
			otherBranches = append(otherBranches, branch)
		}
	}

	selectedBranches := []string{masterBranch}
	if len(releaseBranches) >= 3 {
		selectedBranches = append(selectedBranches, releaseBranches[:3]...)
	} else {
		selectedBranches = append(selectedBranches, releaseBranches...)
	}
	selectedBranches = append(selectedBranches, otherBranches...)
	if len(releaseBranches) > 3 {
		selectedBranches = append(selectedBranches, releaseBranches[3:]...)
	}

	if len(selectedBranches) == 1 {
		return masterBranch, nil
	}

	prompt := promptui.Select{
		Label: "Select branch",
		Items: selectedBranches,
	}

	_, selected, err := prompt.Run()

	if err != nil {
		return "", err
	}

	return selected, nil
}

func reverseStringSlice(s []string) []string {
	a := make([]string, len(s))
	copy(a, s)

	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}

	return a
}

func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}
