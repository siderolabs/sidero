package client

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func NewClient(kubeconfig *string) (dynamic.Interface, error) {
	var (
		config *rest.Config
		err    error
	)

	if kubeconfig == nil {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	c, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return c, nil
}
