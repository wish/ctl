package runs

import (
	"github.com/spf13/cobra"
	"github.com/wish/ctl/cmd/util/parsing"
	"github.com/wish/ctl/pkg/client"
)

func logsCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "logs run [flags]",
		Aliases: []string{"log"},
		Short:   "Get log of a container in a pod",
		Long:    `Print logs from the pods belonging to a cron job run.`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctxs, err := cmd.Flags().GetStringSlice("context")
			if err != nil {
				return err
			}
			namespace, _ := cmd.Flags().GetString("namespace")
			container, _ := cmd.Flags().GetString("container")
			lm, err := parsing.LabelMatchFromCmd(cmd)

			pods, err := c.ListPodsOfJob(ctxs, namespace, args[0], client.ListOptions{lm, nil})
			if err != nil {
				return err
			}

			for _, pod := range pods {
				req, err := c.LogPod(pod.Context, pod.Namespace, pod.Name, container, client.LogOptions{LabelMatch: lm})
				if err != nil {
					return err
				}

				res := req.Do()
				raw, err := res.Raw()
				if err != nil {
					return err
				}
				cmd.Printf("Logs from %s:\n", pod.Name)
				cmd.Print(string(raw))
			}
			return nil
		},
	}

	cmd.Flags().StringP("container", "c", "", "Specify the container")

	return cmd
}
