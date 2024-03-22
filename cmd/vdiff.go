/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/abibby/jit/git"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

// fixCmd represents the fix command
var vdiff = &cobra.Command{
	Use:   "vdiff",
	Short: "Create a fix branch",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := git.GetProvider(cmd.Context())
		if err != nil {
			return err
		}
		u, err := p.DiffURL(cmd.Context())
		if err != nil {
			return err
		}
		return open.Start(u)
	},
}

func init() {
	rootCmd.AddCommand(vdiff)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fixCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fixCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
