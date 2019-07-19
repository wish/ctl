package config

import (
	"github.com/spf13/viper"
)

// GetCtlExt reads the cluster extensions from the config file set in viper
func GetCtlExt() (map[string]map[string]string, error) {
	var m map[string]map[string]string

	err := viper.UnmarshalKey("cluster-ext", &m)

	return m, err
}
