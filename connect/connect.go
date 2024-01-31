package connect

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// CreateConfig creates a Kubernetes configuration based on the provided config file path.
// Parameters:
// - a_config: Pointer to a string containing the path to the Kubernetes config file.
// Returns:
// - (*rest.Config): The built Kubernetes config.
// - (error): An error if any occurred during the config building process.
func CreateConfig(a_config *string) (*rest.Config, error) {
	config_built, err := clientcmd.BuildConfigFromFlags("", *a_config)
	if err != nil {
		return nil, err
	}
	return config_built, nil
}

// NewK8sConnectionConfig creates a Kubernetes clientset using the provided config.
// Parameters:
// - a_config_built: Pointer to a rest.Config containing the Kubernetes config.
// Returns:
// - (*kubernetes.Clientset): The created Kubernetes clientset.
// - (error): An error if any occurred during the clientset creation process.
func NewK8sConnectionConfig(a_config_built *rest.Config) (*kubernetes.Clientset, error) {
	the_clientset, err := kubernetes.NewForConfig(a_config_built)
	if err != nil {
		return nil, err
	}
	return the_clientset, nil
}

// ConnectNow creates a Kubernetes clientset using the provided config file path.
// Parameters:
// - a_config: Pointer to a string containing the path to the Kubernetes config file.
// Returns:
// - (*kubernetes.Clientset): The created Kubernetes clientset.
// - (error): An error if any occurred during the clientset creation process.
func ConnectNow(a_config *string) (*kubernetes.Clientset, error) {
	config_built, err := CreateConfig(a_config)
	if err != nil {
		return nil, err
	}
	return NewK8sConnectionConfig(config_built)
}

// ContextSwitch creates a Kubernetes clientset by building the config with the specified context from a config file.
// Parameters:
// - contextName: The name of the Kubernetes context to switch to.
// - kubeconfig: Pointer to a string containing the path to the Kubernetes config file.
// Returns:
// - (*kubernetes.Clientset): The created Kubernetes clientset.
// - (error): An error if any occurred during the clientset creation process.
func ContextSwitch(contextName string, kubeconfig *string) (*kubernetes.Clientset, error) {
	config, err := BuildConfigWithContextFromFlags(contextName, *kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

// buildConfigWithContextFromFlags creates a Kubernetes configuration by using the provided context and config file path.
// Parameters:
// - context: The name of the Kubernetes context to switch to.
// - kubeconfigPath: The path to the Kubernetes config file.
// Returns:
// - (*rest.Config): The built Kubernetes config.
// - (error): An error if any occurred during the config building process.
func BuildConfigWithContextFromFlags(context string, kubeconfigPath string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}

func ConnectToSource(strSourceClusterContext string, configFile *string) (*kubernetes.Clientset, error) {
	var clientsetToSource *kubernetes.Clientset
	var err error
	if strSourceClusterContext != "" {
		clientsetToSource, err = ContextSwitch(strSourceClusterContext, configFile)
	} else {
		clientsetToSource, err = ConnectNow(configFile)
	}
	if err != nil {
		return nil, err
	}
	return clientsetToSource, nil
}
