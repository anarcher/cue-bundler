package main

import (
	"github.com/anarcher/cue-bundler/pkg/cueutil"

	"github.com/spf13/cobra"
)

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cb",
		Short: "A cue package manager",
	}

	cueModDir := cueutil.FindModDirPath()

	cmd.AddCommand(initCmd(cueModDir))
	cmd.AddCommand(installCmd(cueModDir))
	cmd.AddCommand(updateCmd(cueModDir))
	return cmd
}
