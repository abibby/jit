/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-github/v32/github"
	"github.com/spf13/cobra"
)

// pullrequestCmd represents the pullrequest command
var pullrequestCmd = &cobra.Command{
	Use:     "pull-request",
	Aliases: []string{"pr"},
	Short:   "create a pull request for from this branch",
	Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		root := "./"

		base, err := defaultBranch(cmd.Context())
		if err != nil {
			return err
		}

		templateBytes, err := os.ReadFile(path.Join(root, ".github/pull_request_template.md"))
		if errors.Is(err, os.ErrNotExist) {
		} else if err != nil {
			return err
		}

		issueTag, err := getIssueTag()
		if err != nil {
			return err
		}

		template := string(templateBytes)

		template = fmt.Sprintf("# %s: %s\n\n%s", issueTag, "title", strings.ReplaceAll(template, "MDASH-xxxx", issueTag))

		msgFile := "/tmp/jit-pull-request.md"
		err = os.WriteFile(msgFile, []byte(template), 0644)
		if err != nil {
			return err
		}

		c := exec.Command("code", "--wait", msgFile)
		c.Run()

		commitMsgBytes, err := os.ReadFile(msgFile)
		if err != nil {
			return err
		}

		commitMsg := string(commitMsgBytes)
		title := ""

		if strings.HasPrefix(commitMsg, "#") {
			parts := strings.SplitN(commitMsg, "\n", 2)
			title = strings.TrimSpace(parts[0][1:])
			commitMsg = strings.TrimSpace(parts[1])
		}

		gh := GitHubClient(cmd.Context())
		owner, repo, err := ownerAndRepo()
		if err != nil {
			return err
		}

		branch, err := currentBranch()
		if err != nil {
			return err
		}

		pr, _, err := gh.PullRequests.Create(
			cmd.Context(),
			owner,
			repo,
			&github.NewPullRequest{
				Title: ptr(title),
				Body:  ptr(commitMsg),
				Head:  ptr(branch),
				Base:  ptr(base),
				// Issue
				// MaintainerCanModify
				Draft: ptr(true),
			},
		)
		if err != nil {
			return err
		}
		spew.Dump(pr)
		// exec.Command("code", )
		return nil
	},
}

func ptr[T any](v T) *T {
	return &v
}

func init() {
	rootCmd.AddCommand(pullrequestCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pullrequestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pullrequestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
