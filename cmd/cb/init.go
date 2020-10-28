package main

import "github.com/spf13/cobra"

func initCmd(modDir string) *cobra.Command {
	c := &cobra.Command{
		Use: "init",
	}

	return c
}
