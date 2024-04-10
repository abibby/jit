/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/abibby/jit/git"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
)

// notificationsCmd represents the notifications command
var notificationsCmd = &cobra.Command{
	Use:     "notifications",
	Aliases: []string{"n"},
	Short:   "Show notifications",
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := git.GetProvider(cmd.Context())
		if err != nil {
			return err
		}
		prs, err := p.ListPRs(cmd.Context())
		if err != nil {
			return err
		}
		spew.Dump(prs)
		os.Exit(1)
		// iterate over prs and notify if there are any new notifications
		// maybe add something to show all vs only new comments

		return nil
	},
}

func init() {
	rootCmd.AddCommand(notificationsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// notificationsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// notificationsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
