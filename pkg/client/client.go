package client

import (
	"github.com/wish/ctl/pkg/client/clusterext"
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
	clusterext.Extension
}

// GetPlaceholderClient returns an empty client
func GetPlaceholderClient() *Client {
	return &Client{}
}

// AttachLabelForger creates and adds an Extension to the client
func (c *Client) AttachLabelForger(m map[string]map[string]string) {
	c.Extension = clusterext.Extension{m}
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
	return GetConfigClient(helper.GetKubeConfigPath())
}

// GetConfigClient returns a client with a specific kubeconfig path
func GetConfigClient(path string) *Client {
	contexts := helper.GetContexts(path)
	return &Client{
		clientsetGetter: &configClientsetGetter{
			clientsets: make(map[string]clusterFunctionality),
			config:     path,
		},
		contextsGetter: StaticContextsGetter{
			contexts: contexts,
		},
		Extension: clusterext.EmptyExtension(contexts),
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
		clientsetGetter: &fakeClientsetGetter{
			clientsets: clientsets,
		},
		contextsGetter: StaticContextsGetter{contexts: contexts},
		Extension:      clusterext.EmptyExtension(contexts),
	}
}
