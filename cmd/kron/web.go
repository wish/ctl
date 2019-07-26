package kron

import (
	"github.com/spf13/cobra"
	"github.com/wish/ctl/pkg/client"
	"github.com/wish/ctl/pkg/web"
)

func webCmd(c *client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "web ADDR",
		Short: "Serves a web ui of kron features",
		Long: `Runs a web server on the address.
If no address is specified, runs on localhost:5766.`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				web.Serve(c, ":5766")
			} else {
				web.Serve(c, args[0])
			}
		},
	}
}
