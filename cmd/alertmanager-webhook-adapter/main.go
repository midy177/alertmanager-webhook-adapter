package main

import (
	"fmt"
	"os"

	"alertmanager-webhook-adapter/cmd/alertmanager-webhook-adapter/app"
	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
}

func main() {
	command := app.NewRootCommand()
	if err := command.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
