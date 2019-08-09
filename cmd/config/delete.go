package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wish/ctl/pkg/client"
	"os"
)

func deleteCmd(c *client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "Update extensions",
		RunE: func(cmd *cobra.Command, args []string) error {
			return os.Remove(viper.ConfigFileUsed())
		},
	}
}
