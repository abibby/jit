/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"context"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// branchCmd represents the branch command
var branchCmd = &cobra.Command{
	Use:     "branch <issue id> [title]",
	Aliases: []string{"b"},
	Short:   "Create a new branch from a Jira issue",
	Long:    ``,
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {

		issueID := args[0]
		message := ""
		if len(args) >= 2 || args[1] == "-" {
			message = args[1]
		}

		if err := checkoutDefaultBranch(cmd.Context()); err != nil {
			return err
		}
		c, err := jiraClient()
		if err != nil {
			return err
		}

		issue, _, err := c.Issue.Get(issueID, nil)
		if err != nil {
			return err
		}

		branch := branchName(issue, message)

		if err = git("branch", branch); err != nil {
			return err
		}
		if err = git("checkout", branch); err != nil {
			return err
		}

		if confirm("Do you want to assign yourself to this issue on Jira?", false) {
			u, _, err := c.User.GetSelf()
			if err != nil {
				return err
			}

			_, err = c.Issue.UpdateAssignee(issue.ID, u)
			if err != nil {
				return err
			}
			err = SetStatus(c, issue.ID, viper.GetString("in_progress_status"))
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(branchCmd)
}

func prepBranchName(str string) string {
	str = strings.ReplaceAll(str, " ", "-")
	str = regexp.MustCompile("[^A-Za-z0-9\\-]").ReplaceAllString(str, "")
	str = strings.ToLower(str)
	return str
}

func checkoutDefaultBranch(ctx context.Context) error {
	branch, err := defaultBranch(ctx)

	if err = git("checkout", branch); err != nil {
		return err
	}

	if err = git("pull"); err != nil {
		return err
	}
	return nil
}

func confirm(message string, defaultValue bool) bool {
	cursorPos := 0
	if !defaultValue {
		cursorPos = 1
	}
	prompt := promptui.Select{
		Label:     message,
		CursorPos: cursorPos,
		Items:     []string{"yes", "no"},
	}

	_, selected, err := prompt.Run()
	return err == nil && selected == "yes"
}
