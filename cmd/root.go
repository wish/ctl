package cmd

import (
	"os"

	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	configcmd "github.com/wish/ctl/cmd/config"
	"github.com/wish/ctl/cmd/cron"
	"github.com/wish/ctl/cmd/util/config"
	"github.com/wish/ctl/pkg/client"
)

func cmd() *cobra.Command {
	viper.SetDefault("deadline", 18000)

	viper.SetConfigName("config")
	var conf string
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
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if err = viper.ReadInConfig(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

	var c *client.Client
	if k := viper.GetString("kubeconfig"); len(k) > 0 {
		c = client.GetConfigClient(k)
		// konf = true
	} else {
		c = client.GetDefaultConfigClient()
	}

	m, err := config.GetCtlExt()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if m == nil { // Read map from contexts
		m = c.GetCtlExt()
		config.WriteCtlExt(m)
	}

	c.AttachLabelForger(m)

	if len(m) == 0 {
		fmt.Printf("Config is empty and there are no clusters. "+
			"Please check that the config file at %s is correctly loaded "+
			"and that your kube config is up to date and valid.\n", conf)
	}

	cmd := &cobra.Command{
		Use:          "ctl",
		Short:        "A CLI tool for discovering k8s pods/logs across multiple clusters",
		Long:         "ctl is a CLI tool for easily getting/exec pods/logs across multiple clusters/namespaces.",
		SilenceUsage: true,
	}

	cmd.AddCommand(describeCmd(c))
	cmd.AddCommand(getCmd(c))
	cmd.AddCommand(logsCmd(c))
	cmd.AddCommand(shCmd(c))
	cmd.AddCommand(versionCmd(c))
	cmd.AddCommand(runCmd(c))
	cmd.AddCommand(deleteCmd(c))
	cmd.AddCommand(cpCmd(c))
	cmd.AddCommand(cron.Cmd(c))
	cmd.AddCommand(configcmd.Cmd(c))

	cmd.PersistentFlags().StringSliceP("context", "x", nil, "Specify the context(s) to operate in. Defaults to all contexts.")
	cmd.PersistentFlags().StringP("namespace", "n", "", "Specify the namespace within all the contexts specified. Defaults to all namespaces.")
	cmd.PersistentFlags().StringArrayP("label", "l", nil, "Filter objects by label")
	for _, label := range viper.GetStringSlice("label_flags") {
		cmd.PersistentFlags().String(label, "", "Cluster level label flag \""+label+"\"")
	}

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
