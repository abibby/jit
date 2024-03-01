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
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/range-labs/go-asana/asana"
	"github.com/spf13/cobra"
)

// branchCmd represents the branch command
var branchCmd = &cobra.Command{
	Use:     "branch <issue id> [title]",
	Aliases: []string{"b"},
	Short:   "Create a new branch from a Jira issue",
	Long:    ``,
	Args:    cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {

		issueID := prepIssueID(args[0])

		message := ""
		if len(args) >= 2 && args[1] != "-" {
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

		err = git("branch", branch)
		if err != nil {
			return err
		}
		err = git("checkout", branch)
		if err != nil {
			return err
		}

		// if confirm("Do you want to move this issue to in progress on Jira?", false) {
		// 	u, _, err := c.User.GetSelf()
		// 	if err != nil {
		// 		return err
		// 	}
		// 	if issue.Fields.Assignee.AccountID != u.AccountID {
		// 		_, err = c.Issue.UpdateAssignee(issue.ID, u)
		// 		if err != nil {
		// 			return err
		// 		}
		// 	}

		// 	err = SetStatus(c, issue.ID, configGetString("in_progress_status"))
		// 	if err != nil {
		// 		return err
		// 	}
		// }
		return nil
	},
}

func init() {
	rootCmd.AddCommand(branchCmd)
}

func checkoutDefaultBranch(ctx context.Context) error {
	branch, err := defaultBranch(ctx)
	if err != nil {
		return err
	}
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

func prepIssueID(rawID string) string {
	board := configGetString("board")
	if isNumeric(rawID) && board != "" {
		return board + "-" + rawID
	}
	return rawID
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func creteBranch(ctx context.Context, branchType string, args []string) error {
	if len(args) == 0 {
		client := asanaClient()
		me, err := client.GetAuthenticatedUser(ctx, nil)
		if err != nil {
			return err
		}
		tasks, err := client.ListTasks(ctx, &asana.Filter{
			WorkspaceGID: me.Workspaces[0].GID,
			AssigneeGID:  me.GID,
		})
		if err != nil {
			if asanaErr, ok := err.(*asana.RequestError); ok {
				fmt.Println(asanaErr.Body)
			}
			return err
		}

		taskNames := make([]string, 0, len(tasks))

		for _, t := range tasks {
			if t.Completed {
				continue
			}
			taskNames = append(taskNames, t.Name)
		}

		prompt := promptui.Select{
			Label: "Select branch",
			Items: taskNames,
		}

		_, selected, err := prompt.Run()
		if err != nil {
			return err
		}
		args = append(args, selected)
	}
	if err := checkoutDefaultBranch(ctx); err != nil {
		return err
	}

	branch := branchType + "/" + prepBranchName(strings.Join(args, " "))

	err := git("branch", branch)
	if err != nil {
		return err
	}
	err = git("checkout", branch)
	if err != nil {
		return err
	}
	return nil
}
