package kron

import (
	"flag"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

// Kron Client object for all operations
type Client struct {
	clientset *kubernetes.Clientset
}

func GetKubeConfigPath() string {
	// For multiple calls
	fl := flag.Lookup("kubeconfig")
	if fl != nil {
		return fl.Value.String()
	}
	// Set kubeconfig value
	var kubeconfig *string
	var home string
	if home = os.Getenv("HOME"); home == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	return *kubeconfig
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
		config, err := clientcmd.BuildConfigFromFlags("", GetKubeConfigPath())
		return config, err
	})
	return v, err
}

func GetContextClient(context string) (*Client, error) {
	v, err := clientHelper(func() (*restclient.Config, error) {
		config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: GetKubeConfigPath()},
			&clientcmd.ConfigOverrides{CurrentContext: context}).ClientConfig()
		return config, err
	})
	return v, err
}
