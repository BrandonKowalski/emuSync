package cmd

import (
	"github.com/spf13/cobra"
)

var (
	emuSync = &cobra.Command{
		Use:   "emuSync",
		Short: "A CLI tool for managing Android-Based Emulation Handhelds",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},
	}
)

func Execute() error {
	return emuSync.Execute()
}

func init() {
}
