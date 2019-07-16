package cmd

import (
	"github.com/ContextLogic/ctl/cmd/kron"
	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.PersistentFlags().StringSliceP("context", "x", nil, "Specify the context(s) to operate in")
	rootCmd.PersistentFlags().StringP("namespace", "n", "", "Specify the namespace within all the contexts specified")
	rootCmd.PersistentFlags().StringArrayP("label", "l", nil, "Filter objects by label")

	// Commands'
	c := client.GetDefaultConfigClient()
	rootCmd.AddCommand(GetDescribeCmd(c))
	rootCmd.AddCommand(GetGetCmd(c))
	rootCmd.AddCommand(GetLogsCmd(c))
	rootCmd.AddCommand(GetShCmd(c))
	rootCmd.AddCommand(kron.GetKronCmd(c))
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "ctl",
	Short:        "A CLI tool for discovering k8s pods/logs across multiple clusters",
	Long:         `ctl is a CLI tool for easily getting/exec pods/logs across multiple clusters/namespaces. If you have any questions, problems, or requests please ask #automation.`,
	SilenceUsage: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// No printing of err needed because it already errors??
		os.Exit(1)
	}
}
