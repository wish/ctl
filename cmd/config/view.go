package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wish/ctl/pkg/client"
	"io/ioutil"
	"os"
)

func viewCmd(c *client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "view",
		Short: "print out the cached ctl config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			file := viper.ConfigFileUsed()
			cmd.Printf("Reading contents of %s:\n", file)
			config, err := os.Open(file)
			if err != nil {
				return err
			}
			defer config.Close()
			b, err := ioutil.ReadAll(config)
			if err != nil {
				return err
			}
			cmd.Print(string(b))
			return nil
		},
	}
}
