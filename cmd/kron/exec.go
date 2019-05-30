package kron

import (
	"fmt"
	"github.com/spf13/cobra"
)

// kron/listCmd represents the kron/list command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Executes a job now",
	Long:  "Executes the specified job or the selected job if none.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO")
	},
}

func init() {
	KronCmd.AddCommand(execCmd)
}
