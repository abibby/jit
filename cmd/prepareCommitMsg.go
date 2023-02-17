/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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
	"bytes"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// prepareCommitMsgCmd represents the prepareCommitMsg command
var prepareCommitMsgCmd = &cobra.Command{
	Use:   "prepareCommitMsg",
	Short: "",
	Long:  ``,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		msgFile := args[0]
		commitType := args[1]

		if commitType == "merge" || commitType == "commit" {
			return nil
		}

		issueTag, err := getIssueTag()
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
	},
}

func init() {
	rootCmd.AddCommand(prepareCommitMsgCmd)
}

// #!/bin/sh
// #
// # git prepare-commit-msg hook for automatically prepending an issue key
// # from the start of the current branch name to commit messages.
// # check if commit is merge commit or a commit ammend
// if [ $2 = "merge" ] || [ $2 = "commit" ]; then
//     exit
// fi
// ISSUE_TAG=`git rev-parse --abbrev-ref HEAD | grep -o '[A-Za-z]\{2,\}-[0-9]\+' | head -n 1`
// if [ $? -ne 0 ]; then
//     # no issue key in branch, use the default message
//     exit
// fi
// ISSUE_TAG=`echo "$ISSUE_TAG" | tr '[a-z]' '[A-Z]'`
// if grep "^$ISSUE_TAG" $1 -q; then
//     # dont add the tag if it is already in the commit
//     exit
// fi
// MESSAGE="$(cat $1)"
// # issue key matched from branch prefix, prepend to commit message
// MESSAGE="$ISSUE_TAG: $MESSAGE"
// echo "$MESSAGE" > $1
