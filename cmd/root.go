package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wish/ctl/cmd/kron"
	"github.com/wish/ctl/pkg/client"
	"os"
)

func cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "ctl",
		Short:        "A CLI tool for discovering k8s pods/logs across multiple clusters",
		Long:         `ctl is a CLI tool for easily getting/exec pods/logs across multiple clusters/namespaces. If you have any questions, problems, or requests please ask #automation.`,
		SilenceUsage: true,
	}

	// If config flag is set here...
	c := client.GetDefaultConfigClient()
	cmd.AddCommand(describeCmd(c))
	cmd.AddCommand(getCmd(c))
	cmd.AddCommand(logsCmd(c))
	cmd.AddCommand(shCmd(c))
	cmd.AddCommand(versionCmd(c))
	cmd.AddCommand(kron.Cmd(c))

	cmd.PersistentFlags().StringSliceP("context", "x", nil, "Specify the context(s) to operate in")
	cmd.PersistentFlags().StringP("namespace", "n", "", "Specify the namespace within all the contexts specified")
	cmd.PersistentFlags().StringArrayP("label", "l", nil, "Filter objects by label")

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := cmd().Execute(); err != nil {
		// No printing of err needed because it already errors??
		os.Exit(1)
	}
}
