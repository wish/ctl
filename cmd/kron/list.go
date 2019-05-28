package kron

import (
	"fmt"
	"github.com/spf13/cobra"
)

// kron/listCmd represents the kron/list command
var listCmd = &cobra.Command{
	Use: "list",
  Short: "Get a list of cronjobs",
  Long: "Get a list of cronjobs based on specified search criteria.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO")
	},
}

func init() {
	KronCmd.AddCommand(listCmd)
}
