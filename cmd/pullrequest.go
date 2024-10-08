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
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/abibby/jit/cfg"
	"github.com/abibby/jit/editor"
	"github.com/abibby/jit/git"
	"github.com/spf13/cobra"
)

var tagRE = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`[A-Z]+-(\d+|x+)`)
})

// pullrequestCmd represents the pullrequest command
var pullrequestCmd = &cobra.Command{
	Use:     "pull-request",
	Aliases: []string{"pr"},
	Short:   "create a pull request for from this branch",
	Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		root := "./"

		branch, err := git.CurrentBranch()
		if err != nil {
			return err
		}

		issueTag, err := git.GetIssueID()
		if err != nil {
			return err
		}

		err = git.Run("push")
		if err != nil {
			return err
		}

		title := strings.ReplaceAll(branch[len(issueTag):], "-", " ")
		title = strings.TrimSpace(title)
		title = strings.ToUpper(title[:1]) + title[1:]
		if issueTag != "" {
			title = issueTag + ": " + title
		}
		title = regexp.MustCompile(` +`).ReplaceAllString(title, " ")

		base, err := git.DefaultBranch(cmd.Context())
		if err != nil {
			return err
		}

		template, err := readFirst(
			path.Join(root, ".github/pull_request_template.md"),
			path.Join(cfgDir, "pull_request_template.md"),
		)
		if errors.Is(err, os.ErrNotExist) {
			// empty
		} else if err != nil {
			return err
		}

		template = append(
			[]byte(fmt.Sprintf("# %s\n\n", title)),
			tagRE().ReplaceAll(template, []byte(issueTag))...,
		)

		msgFile := "/tmp/jit-pull-request.md"
		err = os.WriteFile(msgFile, template, 0644)
		if err != nil {
			return err
		}

		err = editor.File(msgFile)
		if err != nil {
			return err
		}

		commitMsgBytes, err := os.ReadFile(msgFile)
		if err != nil {
			return err
		}

		commitMsg := string(commitMsgBytes)

		if strings.HasPrefix(commitMsg, "#") {
			parts := strings.SplitN(commitMsg, "\n", 2)
			title = strings.TrimSpace(parts[0][1:])
			commitMsg = strings.TrimSpace(parts[1])
		}

		if !ask("Are you ready to create this PR?") {
			return nil
		}

		var reviewers []string
		if ask("Do you want to add default reviewers?") {
			reviewers = cfg.GetStringSlice("pullrequest.default_reviewers")
		}
		pr, err := git.CreatePR(cmd.Context(), &git.PullRequestOptions{
			Title:        title,
			Description:  commitMsg,
			SourceBranch: branch,
			BaseBranch:   base,
			Reviewers:    reviewers,
		})
		if err != nil {
			return err
		}

		fmt.Printf("Your PR os up at %s\n", pr.GetURL())
		return nil
	},
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

func readFirst(files ...string) ([]byte, error) {
	for _, file := range files {
		if file == "" {
			continue
		}
		b, err := os.ReadFile(file)
		if err == nil {
			return b, nil
		}
	}
	return nil, os.ErrNotExist
}
