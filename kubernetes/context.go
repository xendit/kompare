package kubernetes

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// SwitchContext switches the context of the Kubernetes client.
func (c *Client) SwitchContext(contextName string) error {
	// using `contextName` context in kubeConfig
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: c.KubeconfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: contextName,
		}).ClientConfig()

	if err != nil {
		return err
	}

	// recreate the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	c.Clientset = clientset
	return nil
}
