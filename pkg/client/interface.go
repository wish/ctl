package client

import (
	"errors"
	"k8s.io/client-go/kubernetes"
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

type fakeClientsetGetter struct {
	clientsets map[string]kubernetes.Interface
}

func (f *fakeClientsetGetter) getContextInterface(context string) (kubernetes.Interface, error) {
	if cs, ok := f.clientsets[context]; ok {
		return cs, nil
	}
	return nil, errors.New("The context specified does not exist")
}
