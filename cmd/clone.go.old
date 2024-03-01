/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-github/v32/github"
	"github.com/spf13/cobra"
)

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "clone a repo",
	Long:  ``,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getRepos(cmd.Context(), toComplete), cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("clone called")

		spew.Dump(getRepos(cmd.Context(), args[0]))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cloneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cloneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getRepos(ctx context.Context, search string) []string {

	repos, _, err := GitHubClient(ctx).Search.Repositories(ctx, search, &github.SearchOptions{})
	if err != nil {
		log.Print(err)
	}
	results := make([]string, len(repos.Repositories))

	for i, r := range repos.Repositories {
		results[i] = *r.CloneURL
	}

	return results
}
