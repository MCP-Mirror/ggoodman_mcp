package main

import (
	"github.com/spf13/cobra"
)

var (
	cmdServe = &cobra.Command{
		Use:   "serve",
		Short: "Commands to start mcp as a server.",
	}
)

func init() {
	cmdServe.AddCommand(cmdServeStdio)
}
