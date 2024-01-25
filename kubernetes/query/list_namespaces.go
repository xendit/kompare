package query

import (
	"context"
	"kompare/kubernetes"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListNamespaces lists namespaces inside a cluster.
func ListNamespaces(client *kubernetes.Client) (*v1.NamespaceList, error) {
	result, err := client.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return result, nil
}
