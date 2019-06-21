package cmd

import (
	"fmt"
	"os"
	"reflect"

	"github.com/ContextLogic/ctl/pkg/util"
	"github.com/spf13/cobra"
)

var (
	ctxs      []string
	namespace string
	filter    util.ContextFilter
)

func init() {
	rootCmd.PersistentFlags().StringSliceVarP(&ctxs, "context", "x", nil, "Specify the context(s) to operate in")
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "Specify the namespace within all the contexts specified")
	rootCmd.PersistentFlags().StringSliceVar(&filter.Region, "region", nil, "Specify the context(s) to operate in")
	rootCmd.PersistentFlags().StringSliceVar(&filter.Az, "az", nil, "Specify the context(s) to operate in")
	rootCmd.PersistentFlags().StringSliceVar(&filter.Env, "env", nil, "Specify the context(s) to operate in")

}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ctl",
	Short: "A CLI tool for discovering k8s pods/logs across multiple clusters",
	Long:  `ctl is a CLI tool for easily getting/exec pods/logs across multiple clusters/namespaces. If you have any questions, problems, or requests please ask #automation.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !reflect.DeepEqual(util.ContextFilter{}, filter) {
			if ctxs != nil {
				fmt.Printf("Cannot specify context and {region, az, env} at the same time")
				os.Exit(1)
			}
			ctxs, _ = util.GetFilteredClusters(filter)
			if len(ctxs) == 0 {
				fmt.Printf("No context(s) with matching criteria found... Exiting")
				os.Exit(1)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
