package kubernetes

import (
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	Clientset      *kubernetes.Clientset
	Config         *rest.Config
	KubeconfigPath string
	ClusterContext string
}

func NewKubernetesClient(kubeconfigPath string) (*Client, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		Clientset:      clientset,
		Config:         config,
		KubeconfigPath: kubeconfigPath,
		ClusterContext: "",
	}, nil
}
