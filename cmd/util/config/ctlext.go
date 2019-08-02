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
	// Get all labels
	allM := make(map[string]struct{})
	var all []string
	// Most common default_columns
	defaultColumnsOcc := make(map[string]int)
	for _, n := range m {
		if defaultColumns, ok := n["_default_columns"]; ok {
			defaultColumnsOcc[defaultColumns]++
			delete(n, "_default_columns")
		}
		for k := range n {
			if _, ok := allM[k]; !strings.HasPrefix(k, "_") && !ok {
				allM[k] = struct{}{}
				all = append(all, k)
			}
		}
	}

	viper.Set("label_flags", all)

	// Default columns
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
