package kron

import (
	"k8s.io/client-go/tools/clientcmd"
)

func getContexts() []string {
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{Precedence: []string{GetKubeConfigPath()}},
		&clientcmd.ConfigOverrides{}).RawConfig()

	if err != nil {
		panic(err.Error())
	}

	ctxs := make([]string, 0, len(config.Contexts))
	for k, _ := range config.Contexts { // Currently ignoring mappings
		ctxs = append(ctxs, k)
	}

	return ctxs
}
