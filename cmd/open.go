/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// openCmd represents the open command
var openCmd = &cobra.Command{
	Use:     "open",
	Aliases: []string{"o"},
	Short:   "Open the jira issue in the browser",
	Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		branch, _, err := execOutput("git", "rev-parse", "--abbrev-ref", "HEAD")

		if err != nil {
			return err
		}

		matches := regexp.MustCompile(regexp.QuoteMeta(viper.GetString("branch_prefix")) + `([A-Za-z]{2,}-[0-9]+)`).FindStringSubmatch(branch)

		c, err := jiraClient()
		if err != nil {
			return err
		}

		issue, _, err := c.Issue.Get(matches[1], nil)
		if err != nil {
			return err
		}

		baseUrl := c.GetBaseURL()

		execOutput("xdg-open", baseUrl.String()+"browse/"+issue.Key)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(openCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// openCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// openCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}