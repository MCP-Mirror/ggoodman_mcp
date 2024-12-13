package main

import "github.com/spf13/cobra"

var (
	cmdPackageInstall = &cobra.Command{
		Use:     "install <package[@version]>",
		Short:   "Install a package from the registry.",
		Aliases: []string{"i"},
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
)
