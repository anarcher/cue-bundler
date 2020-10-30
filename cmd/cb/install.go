package main

import (
	"fmt"
	"os"

	"github.com/anarcher/cue-bundler/pkg/cb"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func installCmd(modDir string) *cobra.Command {
	c := &cobra.Command{
		Use:     "install [<uris>...]",
		Aliases: []string{"i"},
		Short:   "Install new dependencies. Existing ones are silently skipped",
	}

	c.Run = func(cmd *cobra.Command, args []string) {
		if modDir == "" {
			fmt.Println("cue.mod not found")
			os.Exit(1)
		}

		uris := args
		if err := cb.Install(modDir, uris); err != nil {
			fmt.Println(color.RedString("error:"), err)
			os.Exit(1)
		}
		fmt.Println(color.GreenString("success:"), "installed packages")
	}

	return c
}
