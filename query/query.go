package query

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	Corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type KCMetadata struct {
	namespaceName string
}

func ListK8sDeployments(clientset *kubernetes.Clientset, nameSpace string) *v1.DeploymentList {

	if nameSpace == "" {
		nameSpace = "default"
	}
	deployments_list, err := clientset.AppsV1().Deployments(nameSpace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	return deployments_list
}

func ListNameSpaces(clientset *kubernetes.Clientset) *Corev1.NamespaceList {
	nsList, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	return nsList
}
