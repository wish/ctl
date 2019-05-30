package kron

import (
	"fmt"
	"github.com/spf13/cobra"
)

// kron/listCmd represents the kron/list command
var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "Uses list to select a job to operate on",
	Long:  "Uses list to select a job on which other commands can conveniently operate on.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO")
	},
}

func init() {
	KronCmd.AddCommand(selectCmd)
}
