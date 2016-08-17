package main

import (
	"time"

	"github.com/fagongzi/netproxy/cmd/clicmd"
	"github.com/spf13/cobra"
)

const (
	cliName        = "cli"
	cliDescription = "A simple command line client for netproxy."

	defaultDialTimeout    = 2 * time.Second
	defaultCommandTimeOut = 5 * time.Second
)

var (
	rootCmd = &cobra.Command{
		Use:        cliName,
		Short:      cliDescription,
		SuggestFor: []string{"cli"},
	}
)

func main() {
	rootCmd.PersistentFlags().StringVar(&clicmd.Global.Endpoints, "endpoints", "127.0.0.1:8080", "netproxt api address")

	rootCmd.AddCommand(clicmd.NewListCommand(), clicmd.NewUpdateCommand())

	if err := rootCmd.Execute(); err != nil {
		return
	}
}
