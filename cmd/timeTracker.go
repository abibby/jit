/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// timeTrackerCmd represents the timeTracker command
var timeTrackerCmd = &cobra.Command{
	Use:   "time-tracker",
	Short: "Show how long you worked on each task",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("timeTracker called")
	},
}

func init() {
	rootCmd.AddCommand(timeTrackerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// timeTrackerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// timeTrackerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
