package kron

import (
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Kron Client object for all operations
type Client struct {
	clientset *kubernetes.Clientset
}

func clientHelper(getConfig func() (*restclient.Config, error)) (*Client, error) {
	var cl Client // Return client

	config, err := getConfig()

	if err != nil {
		return nil, err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	cl.clientset = clientset
	if err != nil {
		return &cl, err
	}
	return &cl, nil
}

// Creates a kron client from the kubeconfig file
func GetDefaultClient() (*Client, error) {
	v, err := clientHelper(func() (*restclient.Config, error) {
		config, err := clientcmd.BuildConfigFromFlags("", getKubeConfigPath())
		return config, err
	})
	return v, err
}

func GetContextClient(context string) (*Client, error) {
	v, err := clientHelper(func() (*restclient.Config, error) {
		config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: getKubeConfigPath()},
			&clientcmd.ConfigOverrides{CurrentContext: context}).ClientConfig()
		return config, err
	})
	return v, err
}
