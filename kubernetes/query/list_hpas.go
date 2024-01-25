package query

import (
	"context"
	"kompare/kubernetes"

	autoscalingv1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListHorizontalPodAutoscalers lists horizontal pod autoscalers in a namespace.
func ListHorizontalPodAutoscalers(client *kubernetes.Client, namespace string) (*autoscalingv1.HorizontalPodAutoscalerList, error) {
	if namespace == "" {
		namespace = "default"
	}

	result, err := client.Clientset.AutoscalingV1().HorizontalPodAutoscalers(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return result, nil
}
