package helper

import (
	"os"
	"path/filepath"
	"strings"

	"k8s.io/client-go/tools/clientcmd"
)

// GetKubeConfigPath returns the default location of a kubeconfig file
func GetKubeConfigPath() string {
	// TODO: handle multiple paths in KUBECONFIG
	if fl := os.Getenv("KUBECONFIG"); fl != "" {
		return fl
	}
	home, err := os.UserHomeDir()
	if err != nil { // Can't find home dir
		panic(err.Error())
	}

	return filepath.Join(home, ".kube", "config")
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
	for k := range config.Contexts { // Currently ignoring mappings
		// Hardcode ignore test clusters
		if !strings.Contains(k, "test") { // REVIEW: Remove this when possible
			ctxs = append(ctxs, k)
		}
	}

	return ctxs
}
