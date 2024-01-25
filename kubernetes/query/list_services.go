package query

import (
	"context"
	"kompare/kubernetes"

	Corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ListServices(client *kubernetes.Client, namespace string) (*Corev1.ServiceList, error) {
	if namespace == "" {
		namespace = "default"
	}

	result, err := client.Clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return result, nil
}
