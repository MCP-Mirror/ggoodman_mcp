package main

import "github.com/spf13/cobra"

var (
	cmdRegistry = &cobra.Command{
		Use:     "registry",
		Short:   "Commands to interact with the registry.",
		Aliases: []string{"r", "reg"},
	}
)

func init() {
	cmdRegistry.AddCommand(cmdRegistrySearch)
}
