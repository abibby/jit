/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/abibby/jit/jitlog"
	"github.com/abibby/jit/pm"
	"github.com/spf13/cobra"
)

// logTicketsCmd represents the logTickets command
var logTicketsCmd = &cobra.Command{
	Use:   "logTickets",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {

		// logger := slog.New(slog.NewJSONHandler(f, nil))

		// msg, err := os.ReadFile(msgFile)
		// if err != nil {
		// 	return fmt.Errorf("failed to open log: %v", err)
		// }
		// msg = bytes.TrimSpace(msg)

		// logger.Info(string(msg), "repo", parts.String(), "branch", branch)

		c, err := pm.GetProvider()
		if err != nil {
			return err
		}

		myIssues, err := c.GetMyIssues()
		if err != nil {
			return err
		}

		logger, err := jitlog.Logger("issue")
		if err != nil {
			return fmt.Errorf("failed to open logger: %v", err)
		}

		for _, i := range myIssues {
			logger.Info(i.Title, "status", i.Status, "id", i.ID)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(logTicketsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logTicketsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logTicketsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
