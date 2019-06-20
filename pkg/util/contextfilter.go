package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/ContextLogic/ctl/pkg/client/helper"
	"github.com/spf13/viper"
)

type ContextFilter struct {
	Az     []string
	Region []string
	Env    []string
}

// Cluster Struct used to unmarshal yaml config
type Cluster struct {
	Name        string
	Region      string
	Environment string
	Az          string
	Hidden      bool
}

//Config Struct used to unmarshal yaml config
type config struct {
	Clusters []Cluster
}

var clusterInfo map[string]Cluster

func init() {
	configpath := os.Getenv("CTL_CONFIG")
	if configpath == "" {
		configpath, _ = os.Getwd()
		configpath = fmt.Sprintf(configpath + "/config/CTL.yml")
	}

	conf := getConf(configpath)
	clusterInfo = make(map[string]Cluster)
	for _, c := range conf.Clusters {
		clusterInfo[c.Name] = c
	}
}

func filterContexts(contextList []string, region, az, env map[string]bool) []string {
	clusters := make([]string, 0)
	for _, c := range contextList {
		if clusterInfo, ok := clusterInfo[c]; ok {
			if (!clusterInfo.Hidden) &&
				(len(region) == 0 || region[strings.Trim(clusterInfo.Region, " \r\n")]) &&
				(len(az) == 0 || az[strings.Trim(clusterInfo.Az, " \r\n")]) &&
				(len(env) == 0 || env[strings.Trim(clusterInfo.Environment, " \r\n")]) {
				clusters = append(clusters, c)
			}
		} else {
			fmt.Printf("WARNING: the cluster ", c, " is not defined in the configuration",
				" pods in this cluster are not included in the results \n")
		}
	}
	return clusters
}

//GetClusterClusterInfo get AZ, region and env from cluster
func GetClusterClusterInfo(context string) Cluster {
	return clusterInfo[context]
}

//GetFilteredClusters Gets a filterned list of clusters given region, environment and AZ
func GetFilteredClusters(filter ContextFilter) ([]string, error) {
	clusterList := helper.GetContexts()
	regionMap := make(map[string]bool)
	for _, r := range filter.Region {
		regionMap[strings.Trim(r, " \r\n")] = true
	}
	azMap := make(map[string]bool)
	for _, a := range filter.Az {
		azMap[strings.Trim(a, " \r\n")] = true
	}
	envMap := make(map[string]bool)
	for _, e := range filter.Env {
		envMap[strings.Trim(e, " \r\n")] = true
	}

	clusters := filterContexts(clusterList, regionMap, azMap, envMap)

	return clusters, nil
}

// Unmarshal config file
func getConf(configpath string) *config {
	viper.SetConfigFile(configpath)
	err := viper.ReadInConfig()

	if err != nil {
		fmt.Printf("%v", err)
	}

	conf := &config{}
	err = viper.Unmarshal(conf)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}

	return conf
}

func GetEnvs() []string {
	clusterList := helper.GetContexts()

	set := make(map[string]struct{})
	var ret []string

	for _, ctx := range clusterList {
		env := GetClusterClusterInfo(ctx).Environment
		if _, ok := set[env]; !ok {
			set[env] = struct{}{}
			ret = append(ret, env)
		}
	}

	return ret
}
