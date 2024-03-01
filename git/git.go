package git

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
	"github.com/abibby/jit/cfg"
	"github.com/andygrunwald/go-jira"
	"github.com/manifoldco/promptui"
	"golang.org/x/exp/constraints"
)

func Run(options ...string) error {
	fmt.Printf("git %s\n", strings.Join(options, " "))
	return execRaw("git", os.Stdout, os.Stderr, options...)
}

func Root() (string, error) {
	return "./.git", nil
}

func Output(options ...string) (string, string, error) {
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

func BranchName(issue *jira.Issue, message string) string {
	if message == "" {
		message = issue.Fields.Summary
	}
	return cfg.GetString("branch_prefix") + PrepBranchName(issue.Key+" "+message)
}
func PrepBranchName(str string) string {
	str = strings.ReplaceAll(str, " ", "-")
	str = regexp.MustCompile(`[^A-Za-z0-9\-]`).ReplaceAllString(str, "")
	str = strings.ToLower(str)
	str = removeRepeats(str)
	return str
}

func removeRepeats(s string) string {
	result := make([]rune, 0, len(s))
	last := rune(0)
	for _, c := range s {
		if c == last {
			continue
		}
		result = append(result, c)
		last = c
	}
	return string(result)
}

func allBranches() ([]string, error) {
	result, _, err := Output("branch", "-r", "--format", "%(refname)", "--all")
	if err != nil {
		return nil, err
	}

	brancheMap := map[string]struct{}{}
	a := `^(?:refs\/heads\/|refs\/remotes\/[^\/]+\/)(.*)$`
	re := regexp.MustCompile(a)
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

func FindBranch(issueID string) (string, error) {
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

func DefaultBranch(ctx context.Context) (string, error) {
	p, err := GetProvider(ctx)
	if err != nil {
		return "", err
	}

	masterBranch, err := p.MainBranchName(ctx)
	if err != nil {
		return "", err
	}

	branches, err := allBranches()
	if err != nil {
		return "", err
	}

	devBranches := []string{}
	releaseBranches := []string{}
	otherBranches := []string{}

	for _, branch := range reverseStringSlice(branches) {
		if branch == "develop" {
			devBranches = append(devBranches, branch)
		}
		if anyHasPrefix(branch, "release/") {
			releaseBranches = append(releaseBranches, branch)
		}
		if anyHasPrefix(branch, "feature/", "epic/") {
			otherBranches = append(otherBranches, branch)
		}
	}

	selectedBranches := []string{masterBranch}
	selectedBranches = append(selectedBranches, devBranches...)
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

func CurrentBranch() (string, error) {
	branch, _, err := Output("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(branch), nil
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

type GitUrl struct {
	Host  string
	Owner string
	Repo  string
	SSH   bool
}

func UrlParts() (*GitUrl, error) {
	url, _, err := Output("remote", "get-url", "origin")
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`(?:https?:\/\/([^/]+)\/|git@([^:]+):)([^\/]+)\/(.+)\.git`)
	matches := re.FindStringSubmatch(url)
	if len(matches) <= 3 {
		return nil, fmt.Errorf("invalid url")
	}
	host := matches[1]
	ssh := false
	if host == "" {
		host = matches[2]
		ssh = true
	}
	return &GitUrl{
		Host:  host,
		Owner: matches[3],
		Repo:  matches[4],
		SSH:   ssh,
	}, nil
}

func (u *GitUrl) String() string {
	if u.SSH {
		return fmt.Sprintf("git@%s:%s/%s.git", u.Host, u.Owner, u.Repo)
	}
	return fmt.Sprintf("https://%s/%s/%s.git", u.Host, u.Owner, u.Repo)
}
