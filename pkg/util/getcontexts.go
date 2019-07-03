package util

import (
	"github.com/spf13/cobra"
)

func GetContexts(cmd *cobra.Command) ([]string, error) {
	return cmd.Flags().GetStringSlice("context")
}
