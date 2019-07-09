package client

import (
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	restclient "k8s.io/client-go/rest"
)

// Client object for all operations
type Client struct {
	// Add more functionality here...?
	clientsetGetter
}

func clientsetHelper(getConfig func() (*restclient.Config, error)) (kubernetes.Interface, error) {
	config, err := getConfig()

	if err != nil {
		return nil, err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	return clientset, err
}
