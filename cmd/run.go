package cmd

import (
	"encoding/json"
	"errors"
	"github.com/spf13/cobra"
	"github.com/wish/ctl/cmd/util/config"
	"github.com/wish/ctl/cmd/util/parsing"
	"github.com/wish/ctl/pkg/client"
	"github.com/wish/ctl/pkg/client/clusterext"
	"github.com/wish/ctl/pkg/client/filter"
	"os"
	"os/exec"
	"strings"
)

type runDetails struct {
	App       string   `json:"app"`
	Name      string   `json:"name"`
	ImageTag  string   `json:"image_tag"`
	Flags     []string `json:"flags"`
	AfterArgs []string `json:"afterargs"`
}

func runCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run TAG [flags]",
		Short: "Shortcut tool to using kubectl run",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			labelMatch, _ := parsing.LabelMatchFromCmd(cmd)
			overrideTag, _ := cmd.Flags().GetString("tag")

			m, err := config.GetCtlExt()
			if err != nil {
				return err
			}
			e := clusterext.Extension{m}
			ctxs := e.GetFilteredContexts(labelMatch)

			for _, ctx := range ctxs {
				// Find matching command through all runs
				if rawruns, ok := m[ctx]["_run"]; ok {
					runs := make(map[string]runDetails)
					err := json.Unmarshal([]byte(rawruns), &runs)
					if err != nil { // bad
						continue
					}
					// Parse run command specified
					if run, ok := runs[args[0]]; ok {
						combinedargs := []string{"run", "--context", ctx}
						// Find run image
						lm := &filter.LabelMatchEq{"app", run.App}
						l, err := c.ListDeployments(ctx, "", client.ListOptions{LabelMatch: lm})
						if err != nil {
							return err
						}
						for _, d := range l {
							if len(d.Spec.Template.Spec.Containers) > 0 { // Use first image
								image := d.Spec.Template.Spec.Containers[0].Image
								// Change tag on image
								i := strings.LastIndex(image, ":")
								if len(overrideTag) == 0 {
									overrideTag = run.ImageTag
								}
								if i != -1 && len(overrideTag) > 0 {
									image = image[:i+1] + overrideTag
								}
								combinedargs = append(combinedargs, "--image", image)
								break
							}
						}

						if len(run.Flags) > 0 {
							combinedargs = append(combinedargs, run.Flags...)
						}

						combinedargs = append(combinedargs, os.ExpandEnv(run.Name))

						if len(run.AfterArgs) > 0 {
							combinedargs = append(combinedargs, "--")
							combinedargs = append(combinedargs, run.AfterArgs...)
						}

						command := exec.Command("kubectl", combinedargs...)
						command.Stdout = os.Stdout
						command.Stderr = os.Stderr
						command.Stdin = os.Stdin

						cmd.Println("kubectl " + strings.Join(combinedargs, " "))

						return command.Run()
					}
				}
			}
			return errors.New("no command found to run")
		},
	}

	cmd.Flags().String("tag", "", "Change the image tag.")

	return cmd
}
