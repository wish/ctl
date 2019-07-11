package cmd

import (
	"github.com/spf13/cobra"
)

// Version set default value
var Version = "unset"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show ctl version",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println(rootCmd.Use + " version: " + Version)
	},
}
