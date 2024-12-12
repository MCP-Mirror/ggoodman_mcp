package main

import (
	"github.com/spf13/cobra"
)

func main() {
	err := cmdRoot.Execute()
	cobra.CheckErr(err)
}
