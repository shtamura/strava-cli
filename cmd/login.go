/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log/slog"

	"github.com/shtamura/strava-cli/cmd/strava"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Strava.",
	Long:  `Login to Strava.`,
	Run: func(cmd *cobra.Command, args []string) {
		// get environment variable
		config := strava.NewConfig()
		err := strava.Authorize(config)
		if err != nil {
			slog.Error("Failed to authorize: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
