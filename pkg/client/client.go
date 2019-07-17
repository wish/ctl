package client

import (
	"github.com/wish/ctl/pkg/client/helper"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc" // for "oidc" auth provider
	restclient "k8s.io/client-go/rest"
)

// Client object for all operations
type Client struct {
	// Add more functionality here...?
	clientsetGetter
	contextsGetter
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

// GetDefaultConfigClient returns a functioning client from the default kubeconfig path
func GetDefaultConfigClient() *Client {
	return &Client{
		&configClientsetGetter{
			clientsets: make(map[string]kubernetes.Interface),
			config:     helper.GetKubeConfigPath(),
		},
		StaticContextsGetter{
			contexts: helper.GetContexts(helper.GetKubeConfigPath()),
		},
	}
}

// GetFakeConfigClient returns a fake client with the objects in the clusters specified
func GetFakeConfigClient(clusters map[string][]runtime.Object) *Client {
	clientsets := make(map[string]kubernetes.Interface)
	var contexts []string
	for context, objs := range clusters {
		clientsets[context] = fake.NewSimpleClientset(objs...)
		contexts = append(contexts, context)
	}
	return &Client{
		&fakeClientsetGetter{
			clientsets: clientsets,
		},
		StaticContextsGetter{contexts: contexts},
	}
}
