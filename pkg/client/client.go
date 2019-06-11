package client

import (
	"github.com/ContextLogic/ctl/pkg/client/helper"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sync"
)

// Client object for all operations
type Client struct {
	config     string                           // Config file location
	clientsets map[string]*kubernetes.Clientset // maps from context name to client
	cslock     sync.RWMutex                     // For concurrent access of clientsets
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
	return &Client{config: helper.GetKubeConfigPath(), clientsets: make(map[string]*kubernetes.Clientset)}
}

func (c *Client) getContextClientset(context string) (*kubernetes.Clientset, error) {
	c.cslock.RLock()
	if cs, ok := c.clientsets[context]; ok {
		c.cslock.RUnlock()
		return cs, nil
	}
	c.cslock.RUnlock()
	v, err := clientsetHelper(func() (*restclient.Config, error) {
		config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: helper.GetKubeConfigPath()},
			&clientcmd.ConfigOverrides{CurrentContext: context}).ClientConfig()
		return config, err
	})
	if err != nil {
		return nil, err
	}
	c.cslock.Lock()
	c.clientsets[context] = v
	c.cslock.Unlock()
	return v, nil
}
