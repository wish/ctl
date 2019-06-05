package client

import (
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"github.com/ContextLogic/ctl/pkg/client/helper"
)

// Client object for all operations
type Client struct {
	config string	// Config file location
	clientsets map[string]*kubernetes.Clientset // maps from context name to client
}

func clientsetHelper(getConfig func() (*restclient.Config, error)) (*kubernetes.Clientset, error) {
	config, err := getConfig()

	if err != nil {
		return nil, err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	return clientset, err
}

// TODO: Add more constructors??
// Creates a client from the kubeconfig file
func GetDefaultConfigClient() *Client {
	return &Client{helper.GetKubeConfigPath(), make(map[string]*kubernetes.Clientset)}
}

func (c *Client) getContextClientset(context string) (*kubernetes.Clientset, error) {
	if cs, ok := c.clientsets[context]; ok {
		return cs, nil
	}
	v, err := clientsetHelper(func() (*restclient.Config, error) {
		config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: helper.GetKubeConfigPath()},
			&clientcmd.ConfigOverrides{CurrentContext: context}).ClientConfig()
		return config, err
	})
	if err != nil {
		return nil, err
	}
	c.clientsets[context] = v
	return v, nil
}
