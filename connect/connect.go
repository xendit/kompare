package connect

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func CreateConfig(a_config *string) *rest.Config {
	config_built, err := clientcmd.BuildConfigFromFlags("", *a_config)
	if err != nil {
		panic(err.Error())
	}
	return config_built
}

func NewK8sConnectionConfig(a_config_built *rest.Config) *kubernetes.Clientset {
	the_clientset, err := kubernetes.NewForConfig(a_config_built)
	if err != nil {
		panic(err.Error())
	}
	return the_clientset
}

func ConnectNow(a_config *string) *kubernetes.Clientset {
	return NewK8sConnectionConfig(CreateConfig(a_config))
}

func ContextSwitxh(contextName string, kubeconfig *string) *kubernetes.Clientset {
	// using `contextName` context in kubeConfig
	config, err := buildConfigWithContextFromFlags(contextName, *kubeconfig)
	if err != nil {
		panic(err)
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func buildConfigWithContextFromFlags(context string, kubeconfigPath string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}
