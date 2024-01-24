package query

import (
	"context"
	"kompare/kubernetes"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListDeployments lists deployments in a namespace.
func ListDeployments(client *kubernetes.Client, namespace string) (*v1.DeploymentList, error) {
	if namespace == "" {
		namespace = "default"
	}

	result, err := client.Clientset.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// deployments := make([]Deployment, len(result.Items))
	// for i, item := range result.Items {
	// 	deployments[i] = Deployment{}
	// }

	return result, nil
}
