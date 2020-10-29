package main

import (
	"fmt"
	"os"

	"github.com/anarcher/cue-bundler/pkg/cb"
	"github.com/spf13/cobra"
)

func updateCmd(modDir string) *cobra.Command {
	c := &cobra.Command{
		Use:     "update [<uris>...]",
		Aliases: []string{"u"},
		Short:   "Update all or specific dependencies",
	}

	c.Run = func(cmd *cobra.Command, args []string) {
		if modDir == "" {
			fmt.Println("cue.mod not found")
			os.Exit(1)
		}
		uris := args
		if err := cb.Update(modDir, uris); err != nil {
			fmt.Printf("error: %s\n", err)
			os.Exit(1)
		}
	}

	return c
}
