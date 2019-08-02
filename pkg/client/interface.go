package client

import (
	"errors"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sync"
	"k8s.io/apimachinery/pkg/runtime/schema"
	describeversioned "k8s.io/kubectl/pkg/describe/versioned"
	"k8s.io/kubectl/pkg/describe"
)

type clientsetGetter interface {
	getContextInterface(string) (kubernetes.Interface, error)
	getDescriber(string, schema.GroupKind) (describe.Describer, error)
}

type clusterFunctionality struct {
	kubernetes.Interface
	config *restclient.Config
}

type configClientsetGetter struct {
	clientsets map[string]clusterFunctionality
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
	// Get config
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: d.config},
		&clientcmd.ConfigOverrides{CurrentContext: context}).ClientConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	d.cslock.Lock()
	d.clientsets[context] = clusterFunctionality{clientset, config}
	d.cslock.Unlock()
	return clientset, nil
}

func (d *configClientsetGetter) getDescriber(context string, kind schema.GroupKind) (describe.Describer, error) {
	_, err := d.getContextInterface(context)
	if err != nil {
		return nil, err
	}
	config := d.clientsets[context].config
	describer, ok := describeversioned.DescriberFor(kind, config)
	if !ok {
		return nil, errors.New("could not retrieve describer for context " + context)
	}
	return describer, nil
}

type fakeClientsetGetter struct {
	clientsets map[string]kubernetes.Interface
}

func (f *fakeClientsetGetter) getContextInterface(context string) (kubernetes.Interface, error) {
	if cs, ok := f.clientsets[context]; ok {
		return cs, nil
	}
	return nil, errors.New("the context specified does not exist")
}

func (*fakeClientsetGetter) getDescriber(string, schema.GroupKind) (describe.Describer, error) {
	return nil, errors.New("fake client cannot describe")
}
