package runs

import (
	"fmt"
	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/spf13/cobra"
	"os"
)

func GetLogsCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "logs pod [flags]",
		Aliases: []string{"log"},
		Short:   "Get log of a container in a pod",
		Long: `Print logs from the pods belonging to a cron job run.
	If namespace not specified, it will get all the pods across all the namespaces.
	If context(s) not specified, it will search through all contexts.`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctxs, err := cmd.Flags().GetStringSlice("context")
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			namespace, _ := cmd.Flags().GetString("namespace")
			container, _ := cmd.Flags().GetString("container")

			pods, err := c.ListPodsOfRun(ctxs, namespace, args[0], client.ListOptions{})
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			for _, pod := range pods {
				res, err := c.LogPod(pod.Context, pod.Namespace, pod.Name, container, client.LogOptions{})

				raw, err := res.Raw()
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
				// REVIEW: Format??
				fmt.Printf("Logs from %s:\n", pod.Name)
				fmt.Print(string(raw))
				fmt.Println("------")
			}
		},
	}

	cmd.Flags().StringP("container", "c", "", "Specify the container")

	return cmd
}
