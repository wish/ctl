package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
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

const (
	// DefaultDeadline - default amount of time to keep the ad hoc pods running
	DefaultDeadline string = "120"
	// MaxDeadline sets the max deadline for a job (defaulted to 1 day)
	MaxDeadline int = 60 * 60 * 24
)

func runCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run APPNAME [flags]",
		Short: "Shortcut tool to using kubectl run",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			labelMatch, _ := parsing.LabelMatchFromCmd(cmd)
			deadline, _ := cmd.Flags().GetString("deadline")

			// Check for valid input for deadline and set default if needed
			if deadlineInt, err := strconv.Atoi(deadline); err != nil || deadlineInt < 1 || deadlineInt > MaxDeadline {
				fmt.Printf("Setting default deadline of %v seconds \n", DefaultDeadline)
				deadline = DefaultDeadline
			}

			appName := args[0]

			// Get all kubernetes contexts from config file
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

					// Get a list of available app names
					availableAppNames := make([]string, len(runs))
					i := 0
					for k := range runs {
						availableAppNames[i] = k
						i++
					}

					// Check if the app name exists in the raw runs
					if run, ok := runs[appName]; ok {
						// Check if the app is active
						if run.Active {
							// Get hostname to use in job name
							user, err := os.Hostname()
							if err != nil {
								return errors.New("Unable to get hostname of machine")
							}
							// Template out hostname into job name in manifest
							manifest := regexp.MustCompile(`({USER})`).ReplaceAllString(run.Manifest, user)
							// Template active deadline seconds into manifest
							manifest = regexp.MustCompile(`("{ACTIVE_DEADLINE_SECONDS}")`).ReplaceAllString(manifest, deadline)

							// Pass the manifest into a reader for stdin
							r := strings.NewReader(manifest)
							command := exec.Command("kubectl", "apply", "-f", "-")
							command.Stdout = os.Stdout
							command.Stderr = os.Stderr
							command.Stdin = r

							fmt.Printf("Running %v in %v with a deadline of %vs\n", appName, ctx, deadline)

							return command.Run()
						}
						// If the app name is not active, let's just print its manifest
						cmd.Printf("WARN: App name not active in %v. Manifest file: %v\n", ctx, run.Manifest)
					} else {
						// App name not found, list out available app names to the user
						fmt.Printf("Available app names: %v\n", availableAppNames)
						return errors.New("App name not found")
					}
				}
			}
			return errors.New("no command found to run or not authorized to access cluster")
		},
	}

	cmd.Flags().String("deadline", "", "Time pod will stay alive in seconds")

	return cmd
}
