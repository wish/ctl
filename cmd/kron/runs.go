package kron

import (
	//"fmt"
	"github.com/spf13/cobra"
)

var runsCmd = &cobra.Command{
	Use:   "runs",
	Short: "Subcommand for operating on runs of a command",
	Long:  "", // TODO
}

func init() {
	KronCmd.AddCommand(runsCmd)
}
