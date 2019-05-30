package kron

import (
	// "fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// kron/listCmd represents the kron/list command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get info about a job",
	Long:  "Get info about a specific job, or the selected job if none is specified.",
	Run: func(cmd *cobra.Command, args []string) {
		color.Blue("TODO")
	},
}

func init() {
	KronCmd.AddCommand(infoCmd)
}
