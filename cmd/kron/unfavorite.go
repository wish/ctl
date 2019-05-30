package kron

import (
	"fmt"
	"github.com/spf13/cobra"
)

var unfavoriteCmd = &cobra.Command{
	Use:   "unfavorite",
	Short: "Removes job(s) from favorite list",
	Long:  "Removes job(s) from favorite list. If no jobs are specified, removes selected job. If job is selected, opens a list to choose to remove from.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO")
	},
}

func init() {
	KronCmd.AddCommand(unfavoriteCmd)
}
