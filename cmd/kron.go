package cmd

import (
	"github.com/ContextLogic/ctl/cmd/kron"
)

func init() {
	rootCmd.AddCommand(kron.KronCmd)
}
