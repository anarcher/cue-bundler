package main

import (
	"fmt"
	"os"

	"github.com/anarcher/cue-bundler/pkg/cb"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func initCmd(modDir string) *cobra.Command {
	c := &cobra.Command{
		Use:   "init",
		Short: " Initialize a new empty config cue file",
	}
	c.Run = func(cmd *cobra.Command, args []string) {
		if modDir == "" {
			fmt.Println("cue.mod not found")
			os.Exit(1)
		}

		if err := cb.Init(modDir); err != nil {
			fmt.Println(color.RedString("error:"), err)
			os.Exit(1)
		}
		fmt.Println(color.GreenString("success:"), "init config cue files")
	}
	return c
}
