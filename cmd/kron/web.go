package kron

import (
	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/ContextLogic/ctl/pkg/web"
	"github.com/spf13/cobra"
)

func GetWebCmd(c *client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "web [port]",
		Short: "Serves a web ui of kron features",
		Long: `Runs a web server on the specified localhost port.
If no port is specified, runs on port 5766.`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				web.Serve(":5766")
			} else {
				web.Serve(":" + args[0])
			}
		},
	}
}
