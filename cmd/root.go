package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	configcmd "github.com/wish/ctl/cmd/config"
	"github.com/wish/ctl/cmd/kron"
	"github.com/wish/ctl/cmd/util/config"
	"github.com/wish/ctl/pkg/client"
)

func cmd() *cobra.Command {
	// Placeholder client
	c := client.GetPlaceholderClient()

	cmd := &cobra.Command{
		Use:          "ctl",
		Short:        "A CLI tool for discovering k8s pods/logs across multiple clusters",
		Long:         "ctl is a CLI tool for easily getting/exec pods/logs across multiple clusters/namespaces.",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// REVIEW: this is quite sketchy. Should use another method
			// konf := false
			if k, _ := cmd.Flags().GetString("kubeconfig"); len(k) > 0 {
				*c = *client.GetConfigClient(k)
				// konf = true
			} else {
				*c = *client.GetDefaultConfigClient()
			}

			viper.SetConfigName("config")
			conf, _ := cmd.Flags().GetString("config")
			if len(conf) == 0 {
				if v, ok := os.LookupEnv("XDG_CONFIG_DIR"); ok {
					conf = v + "/ctl/config.yml"
				} else {
					conf = os.Getenv("HOME") + "/.config/ctl/config.yml"
				}
			}
			viper.SetConfigFile(conf)
			if err := viper.ReadInConfig(); err != nil {
				err = config.Create(conf)
				if err != nil {
					return err
				}
				if err = viper.ReadInConfig(); err != nil {
					return err
				}
			}

			m, err := config.GetCtlExt()
			if err != nil {
				return err
			}
			if m == nil { // Read map from contexts
				m = c.GetCtlExt()
				viper.Set("cluster-ext", m)
				viper.WriteConfig()
			}

			c.AttachLabelForger(m)

			return nil
		},
	}

	cmd.AddCommand(describeCmd(c))
	cmd.AddCommand(getCmd(c))
	cmd.AddCommand(logsCmd(c))
	cmd.AddCommand(shCmd(c))
	cmd.AddCommand(versionCmd(c))
	cmd.AddCommand(kron.Cmd(c))
	cmd.AddCommand(configcmd.Cmd(c))

	cmd.PersistentFlags().StringSliceP("context", "x", nil, "Specify the context(s) to operate in")
	cmd.PersistentFlags().StringP("namespace", "n", "", "Specify the namespace within all the contexts specified")
	cmd.PersistentFlags().StringArrayP("label", "l", nil, "Filter objects by label")
	cmd.PersistentFlags().String("kubeconfig", "", "Custom kubeconfig file")
	cmd.PersistentFlags().String("config", "", "Choose a different directory for ctl config")

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
