package kron

import (
	"fmt"
	"github.com/spf13/cobra"
)

// kron/listCmd represents the kron/list command
var summarizeCmd = &cobra.Command{
	Use:   "summarize",
	Short: "Get aggregate info about a group of jobs",
	Long:  "Get a summary about a group of jobs (still need to determine groups)", // TODO be more descriptive about what a summary is.
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO")
	},
}

func init() {
	KronCmd.AddCommand(summarizeCmd)
}
