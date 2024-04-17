package hooks

import (
	"bytes"
	"fmt"
	"os"

	"github.com/abibby/jit/git"
	"github.com/abibby/jit/jitlog"
)

type PrepareCommitMsg func(msgFile, commitType string) error

var PrepareCommitMsgHooks = []PrepareCommitMsg{
	AddIssueTagToCommit,
	LogCommits,
}

func AddIssueTagToCommit(msgFile, commitType string) error {

	if commitType == "merge" || commitType == "commit" {
		return nil
	}

	issueTag, err := git.GetIssueID()
	if err != nil {
		return err
	}
	if issueTag == "" {
		return nil
	}

	commitMsg, err := os.ReadFile(msgFile)
	if err != nil {
		return err
	}

	if bytes.HasPrefix(commitMsg, []byte(issueTag+": ")) {
		return nil
	}

	f, err := os.OpenFile(msgFile, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "%s: %s", issueTag, commitMsg)
	if err != nil {
		return err
	}
	return nil
}

func LogCommits(msgFile, commitType string) error {
	branch, err := git.CurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %v", err)
	}

	parts, err := git.UrlParts()
	if err != nil {
		return fmt.Errorf("failed to get url parts: %v", err)
	}

	logger, err := jitlog.Logger("commit")
	if err != nil {
		return fmt.Errorf("failed to open logger: %v", err)
	}

	msg, err := os.ReadFile(msgFile)
	if err != nil {
		return fmt.Errorf("failed to open log: %v", err)
	}
	msg = bytes.TrimSpace(msg)

	logger.Info(string(msg), "repo", parts.String(), "branch", branch)

	return nil
}
