/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"net/http"

	"github.com/shtamura/strava-cli/cmd/strava"
	"github.com/spf13/cobra"
)

// athleteCmd represents the athlete command
var athleteCmd = &cobra.Command{
	Use:   "athlete",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// http GET "https://www.strava.com/api/v3/athlete" "Authorization: Bearer [[token]]"
		// request to strava api
		fmt.Println("athlete called")
		credential, err := strava.GetCredential()
		if err != nil {
			// TODO: with tired runner pictgram
			fmt.Println("Failed to get credential: %v", err)
			return
		}
		// http request with token header
		req, _ := http.NewRequest("GET", "https://www.strava.com/api/v3/athlete", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", credential.AccessToken))
		client := new(http.Client)
		res, err := client.Do(req)
		if err != nil {
			fmt.Println("Failed to request: %v", err)
			return
		}
		defer res.Body.Close()
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Failed to read body: %v", err)
			return
		}
		fmt.Println(string(body))
	},
}

func init() {
	rootCmd.AddCommand(athleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// athleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// athleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
