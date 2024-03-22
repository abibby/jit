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

	"github.com/abibby/jit/cfg"
	"github.com/abibby/jit/git"
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
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		issueID := prepIssueID(args[0])

		message := ""
		if len(args) >= 2 && args[1] != "-" {
			message = strings.Join(args[1:], " ")
		}

		c, err := jiraClient()
		if err != nil {
			return err
		}

		issue, _, err := c.Issue.Get(issueID, nil)
		if err != nil {
			return err
		}

		branch := git.BranchName(issue, message)

		max := cfg.GetIntDefault("branch.max_name", 47)
		for len(branch) > max {
			p := &promptui.Prompt{
				Label:     fmt.Sprintf("Branch name too long (max %d), %s (%d)", max, branch, len(branch)),
				Default:   message,
				AllowEdit: true,
			}
			v, err := p.Run()
			if err != nil {
				return err
			}

			oldBranch := branch
			message = v
			branch = git.BranchName(issue, message)
			if branch == oldBranch {
				break
			}
		}

		if err := checkoutDefaultBranch(cmd.Context()); err != nil {
			return err
		}

		err = git.Run("branch", branch)
		if err != nil {
			return err
		}
		err = git.Run("checkout", branch)
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

		// 	err = SetStatus(c, issue.ID, cfg.GetString("in_progress_status"))
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
	branch, err := git.DefaultBranch(ctx)
	if err != nil {
		return err
	}
	if err = git.Run("checkout", branch); err != nil {
		return err
	}

	if err = git.Run("pull"); err != nil {
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
	board := cfg.GetString("board")
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

	branch := branchType + "/" + git.PrepBranchName(strings.Join(args, " "))

	err := git.Run("branch", branch)
	if err != nil {
		return err
	}
	err = git.Run("checkout", branch)
	if err != nil {
		return err
	}
	return nil
}
