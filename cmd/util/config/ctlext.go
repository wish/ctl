package config

import (
	"github.com/spf13/viper"
	"strings"
)

// GetCtlExt reads the cluster extensions from the config file set in viper
func GetCtlExt() (map[string]map[string]string, error) {
	var m map[string]map[string]string

	err := viper.UnmarshalKey("cluster-ext", &m)

	return m, err
}

// WriteCtlExt processes the ctl extension and writes it to the config file
func WriteCtlExt(m map[string]map[string]string) {
	defaultColumnsOcc := make(map[string]int)
	for _, n := range m {
		if defaultColumns, ok := n["default_columns"]; ok {
			defaultColumnsOcc[defaultColumns]++
			delete(n, "default_columns")
		}
	}
	max := 0
	for col, occ := range defaultColumnsOcc {
		if occ > max {
			max = occ
			viper.Set("default_columns", strings.Split(col, ","))
		}
	}
	viper.Set("cluster-ext", m)
	viper.WriteConfig()
}
