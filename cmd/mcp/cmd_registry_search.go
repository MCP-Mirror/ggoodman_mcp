package main

import (
	"mcp/internal/registry"

	"github.com/spf13/cobra"
)

var (
	cmdRegistrySearch = &cobra.Command{
		Use:     "search [query...]",
		Short:   "Search the registry for a package.",
		Aliases: []string{"s"},
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			rc := &registry.FakeClient{}

			pkgs, err := rc.SearchIntegrations(ctx, args...)
			cobra.CheckErr(err)

			cmd.PrintErrf("Found %d packages\n", len(pkgs))

			for _, pkg := range pkgs {
				cmd.PrintErrf("%s %s\n", pkg.Name, pkg.Version)
			}
		},
	}
)

func init() {
}
