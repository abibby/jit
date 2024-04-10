/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// postCommitCmd represents the postCommit command
var postCommitCmd = &cobra.Command{
	Use:    "postCommit",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return os.WriteFile("./commit", []byte(fmt.Sprintf("%#v", args)), 0o644)
	},
}

func init() {
	rootCmd.AddCommand(postCommitCmd)
}
