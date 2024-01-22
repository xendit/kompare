package connect

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func CreateConfig(a_config *string) (*rest.Config, error) {
	config_built, err := clientcmd.BuildConfigFromFlags("", *a_config)
	if err != nil {
		return nil, err
	}
	return config_built, nil
}

func NewK8sConnectionConfig(a_config_built *rest.Config) (*kubernetes.Clientset, error) {
	the_clientset, err := kubernetes.NewForConfig(a_config_built)
	if err != nil {
		return nil, err
	}
	return the_clientset, nil
}

func ConnectNow(a_config *string) (*kubernetes.Clientset, error) {
	config_built, err := CreateConfig(a_config)
	if err != nil {
		return nil, err
	}
	return NewK8sConnectionConfig(config_built)
}

func ContextSwitch(contextName string, kubeconfig *string) (*kubernetes.Clientset, error) {
	config, err := buildConfigWithContextFromFlags(contextName, *kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func buildConfigWithContextFromFlags(context string, kubeconfigPath string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}
