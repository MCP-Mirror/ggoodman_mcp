package main

import "github.com/spf13/cobra"

var (
	cmdPackage = &cobra.Command{
		Use:     "package",
		Aliases: []string{"p", "pkg"},
		Short:   "Manage installed packages.",
	}
)

func init() {
	cmdPackage.AddCommand(cmdPackageInstall)
	cmdPackage.AddCommand(cmdPackageList)
}
