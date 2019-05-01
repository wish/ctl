package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version set default value
var Version = "unset"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show wishctl version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(rootCmd.Use + " version: " + Version)
	},
}
