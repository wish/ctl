package kron

import (
	"github.com/ContextLogic/ctl/cmd/kron/runs"
)

func init() {
	KronCmd.AddCommand(runs.RunsCmd)
}
