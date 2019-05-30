package kron

import (
	"fmt"
	"github.com/spf13/cobra"
)

var favoriteCmd = &cobra.Command{
	Use:   "favorite",
	Short: "Adds a job to favorite list",
	Long:  "Adds specified job(s) to the favorite list. If no job was specified the selected job is added.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO")
	},
}

func init() {
	KronCmd.AddCommand(favoriteCmd)
}
