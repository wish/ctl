package cmd

import (
	"encoding/json"
	"errors"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/wish/ctl/cmd/util/config"
	"github.com/wish/ctl/cmd/util/parsing"
	"github.com/wish/ctl/pkg/client"
	"github.com/wish/ctl/pkg/client/clusterext"
)

type runDetails struct {
	Active   bool   `json:"active"`
	Manifest string `json:"manifest"`
}

func runCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run APPNAME [flags]",
		Short: "Shortcut tool to using kubectl run",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			labelMatch, _ := parsing.LabelMatchFromCmd(cmd)
			activeDeadline, _ := cmd.Flags().GetInt("deadline")

			m, err := config.GetCtlExt()
			if err != nil {
				return err
			}
			e := clusterext.Extension{ClusterExt: m, K8Envs: nil}
			ctxs := e.GetFilteredContexts(labelMatch)

			// Random shuffle contexts
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(ctxs), func(i, j int) { ctxs[i], ctxs[j] = ctxs[j], ctxs[i] })

			for _, ctx := range ctxs {

				if rawruns, ok := m[ctx]["_run"]; ok {
					runs := make(map[string]runDetails)
					err := json.Unmarshal([]byte(rawruns), &runs)
					if err != nil { // bad
						continue
					}

					// Parse run command specified
					if run, ok := runs[args[0]]; ok {
						if run.Active {
							user, err := os.Hostname()
							if err != nil {
								errors.New("Unnable to get hostname of machine")
							}
							manifest := regexp.MustCompile(`({USER}\)`).ReplaceAllString(run.Manifest, user)
							manifest = regexp.MustCompile(`("{ACTIVE_DEADLINE_SECONDS}\")`).ReplaceAllString(manifest, strconv.Itoa(activeDeadline))

							r := strings.NewReader(manifest)
							command := exec.Command("kubectl", "apply", "-f", "-")
							command.Stdout = os.Stdout
							command.Stderr = os.Stderr
							command.Stdin = r

							return command.Run()
						}
						cmd.Printf("WARN: Command found in a cluster but is disabled with message: %s", run.Manifest)
					}
				}
			}
			return errors.New("no command found to run or not authorized to access cluster")
		},
	}

	cmd.Flags().String("deadline", "", "Time pod will stay alive in seconds")

	return cmd
}
