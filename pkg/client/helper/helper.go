package helper

import (
	"os"
	"path/filepath"
	"k8s.io/client-go/tools/clientcmd"
	"strings"
)

// GetKubeConfigPath returns the default location of a kubeconfig file
func GetKubeConfigPath() []string {
	if fl := os.Getenv("KUBECONFIG"); fl != "" {
		// split KUBECONFIG string to handle multiple kube config files
		kubeConfigs := strings.Split(fl, ":")
		return kubeConfigs
	}
	home, err := os.UserHomeDir()
	if err != nil { // Can't find home dir
		panic(err.Error())
	}

	return []string {filepath.Join(home, ".kube", "config")}
}

// GetContexts returns a list of clusters from a config file
func GetContexts(configpath string) []string {
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: configpath},
		&clientcmd.ConfigOverrides{}).RawConfig()

	if err != nil {
		panic(err.Error())
	}

	ctxs := make([]string, 0, len(config.Contexts))
	for k := range config.Contexts {
		ctxs = append(ctxs, k)
	}

	return ctxs
}
