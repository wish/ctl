package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.PersistentFlags().StringSliceP("context", "c", nil, "Specify the context(s) to operate in")
	rootCmd.PersistentFlags().StringP("namespace", "n", "", "Specify the namespace within all the contexts specified")
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ctl",
	Short: "A CLI tool for discovering k8s pods/logs across multiple clusters",
	Long:  `ctl is a CLI tool for easily getting/exec pods/logs across multiple clusters/namespaces. If you have any questions, problems, or requests please ask #automation.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
