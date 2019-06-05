package kron

import (
	"github.com/spf13/cobra"
)

func init() {
	// Nothing
}

var KronCmd = &cobra.Command{
	Use:   "kron",
	Short: "A tool for cron on kubernetes",
	Long:  "A subcommand for managing and reviewing cron jobs on kubernetes.",
}
