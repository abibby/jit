/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/abibby/jit/git"
	"github.com/spf13/cobra"
)

// bitbucketCmd represents the bitbucket command
var bitbucketCmd = &cobra.Command{
	Use:     "bitbucket",
	Aliases: []string{"bb"},
	Short:   "",
	Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		bb := git.NewBitbucket(cmd.Context())
		members, err := bb.ListUsers(cmd.Context())
		if err != nil {
			return err
		}
		for _, m := range members.Members {
			fmt.Printf("%30s %s\n", m.DisplayName, m.Uuid)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(bitbucketCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bitbucketCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bitbucketCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
