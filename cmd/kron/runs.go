package kron

import (
	"github.com/spf13/cobra"
)

func init() {
	KronCmd.AddCommand(runsCmd)
}

var runsCmd = &cobra.Command{
	Use:   "runs",
	Short: "Subcommand for operating on runs of a command",
	Long:  "", // TODO
}
