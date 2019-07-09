package client

import (
	"errors"
	"github.com/ContextLogic/ctl/pkg/client/helper"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sync"
)

type clientsetGetter interface {
	getContextInterface(string) (kubernetes.Interface, error)
}

type configClientsetGetter struct {
	clientsets map[string]kubernetes.Interface
	config     string
	cslock     sync.RWMutex
}

func (d *configClientsetGetter) getContextInterface(context string) (kubernetes.Interface, error) {
	d.cslock.RLock()
	if cs, ok := d.clientsets[context]; ok {
		d.cslock.RUnlock()
		return cs, nil
	}
	d.cslock.RUnlock()
	v, err := clientsetHelper(func() (*restclient.Config, error) {
		config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: d.config},
			&clientcmd.ConfigOverrides{CurrentContext: context}).ClientConfig()
		return config, err
	})
	if err != nil {
		return nil, err
	}
	d.cslock.Lock()
	d.clientsets[context] = v
	d.cslock.Unlock()
	return v, nil
}

func GetDefaultConfigClient() *Client {
	return &Client{
		&configClientsetGetter{
			clientsets: make(map[string]kubernetes.Interface),
			config:     helper.GetKubeConfigPath(),
		},
	}
}

type fakeClientsetGetter struct {
	clientsets map[string]kubernetes.Interface
	cslock     sync.Mutex
}

func (f *fakeClientsetGetter) getContextInterface(context string) (kubernetes.Interface, error) {
	f.cslock.Lock()
	if cs, ok := f.clientsets[context]; ok {
		f.cslock.Unlock()
		return cs, nil
	}
	return nil, errors.New("The context specified does not exist")
}

func GetFakeConfigClient(clusters map[string][]runtime.Object) *Client {
	clientsets := make(map[string]kubernetes.Interface)
	for context, objs := range clusters {
		clientsets[context] = fake.NewSimpleClientset(objs...)
	}
	return &Client{
		&fakeClientsetGetter{
			clientsets: clientsets,
		},
	}
}
