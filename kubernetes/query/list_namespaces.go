package query

import (
	"context"
	"kompare/kubernetes"
	"kompare/kubernetes/dao"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListNamespaces lists namespaces inside a cluster.
func ListNamespaces(client *kubernetes.Client) ([]dao.Namespace, error) {
	result, err := client.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	namespaces := make([]dao.Namespace, len(result.Items))
	for i, item := range result.Items {
		namespaces[i] = dao.Namespace{
			Name: item.Name,
		}
	}

	return namespaces, nil
}
