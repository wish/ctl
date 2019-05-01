package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wishctl",
	Short: "A CLI tool for discovering k8s pods/logs across multiple clusters",
	Long:  `wishctl is a CLI tool for easily getting/exec pods/logs across multiple clusters/namespaces. If you have any questions, problems, or requests please ask #automation.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		_, err := exec.LookPath("kubectl")
		if err != nil {
			fmt.Println("kubectl not installed. k8s commands are unavaible.")
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}