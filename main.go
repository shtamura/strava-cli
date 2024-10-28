/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log/slog"

	"github.com/shtamura/strava-cli/cmd"
)

func main() {
	// slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	slog.SetDefault(slog.New(slog.Default().Handler()))
	slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.Info("hoge")

	cmd.Execute()
}
