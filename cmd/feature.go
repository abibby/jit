/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// fixCmd represents the fix command
var featureCmd = &cobra.Command{
	Use:     "feature",
	Aliases: []string{"fe"},
	Short:   "Create a feature branch",
	Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return creteBranch(cmd.Context(), "feature", args)
	},
}

func init() {
	rootCmd.AddCommand(featureCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fixCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fixCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
