package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wish/ctl/pkg/client"
)

func fetchCmd(c *client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "fetch",
		Short: "Update extensions",
		Run: func(cmd *cobra.Command, args []string) {
			m := c.GetCtlExt()
			viper.Set("cluster-ext", m)
			viper.WriteConfig()
		},
	}
}
