/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type Commit struct {
	Repo    string
	Branch  string
	Summary string
	Author  string
	Time    time.Time
}

// trackTimeCmd represents the trackTime command
var trackTimeCmd = &cobra.Command{
	Use:     "track-time",
	Aliases: []string{"t"},
	Short:   "",
	Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		o, e, err := gitOutput("log", "--pretty=format:%s%n%ct")
		if err != nil {
			return fmt.Errorf("%s: %w", strings.TrimSpace(e), err)
		}

		fmt.Print(o)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(trackTimeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// trackTimeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// trackTimeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
