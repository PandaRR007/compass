package sync_cmd

import (
	"github.com/spf13/cobra"
)

func Run() {
	var rootCmd = &cobra.Command{Use: "map_rly"}

	rootCmd.AddCommand(cmdConfigFunc())
	err := rootCmd.Execute()
	if err != nil {
		return
	}
}
