package main

import (
	"fmt"
	"os"

	"github.com/anarcher/cue-bundler/pkg/cb"
	"github.com/spf13/cobra"
)

func installCmd(cueModDir string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "install [<uris>...]",
		Aliases: []string{"i"},
		Short:   "Install new dependencies. Existing ones are silently skipped",
	}

	cmd.Run = func(cmd *cobra.Command, args []string) {
		if cueModDir == "" {
			fmt.Println("cue.mod not found")
			os.Exit(1)
		}

		uris := args
		if err := cb.Install(cueModDir, uris); err != nil {
			fmt.Printf("error: %s\n", err)
			os.Exit(1)
		}
	}

	return cmd
}
