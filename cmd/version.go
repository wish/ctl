package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wish/ctl/pkg/client"
)

// Version set default value
var Version string = "v12.0.0"
var Commit string

func versionCmd(*client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show ctl version",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(cmd.Root().Use, "version", Version)
			cmd.Println("commit sha", Commit)
		},
	}
}
