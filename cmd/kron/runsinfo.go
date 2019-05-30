package kron

import (
	"fmt"
	"github.com/spf13/cobra"
)

var runsInfoCmd = &cobra.Command{
	Use:   "info [job [run] | run]",
	Short: "Get info about a run",
	Long:  "Get information about a specific run of a cron job.", // TODO
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO")
	},
}

func init() {
	runsCmd.AddCommand(runsInfoCmd)
}
