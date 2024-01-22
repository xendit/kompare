package query

import (
	"context"
	"kompare/kubernetes"

	autoscalingv1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TODO: to decide which fields to take
type HorizontalPodAutoscaler struct {
}

func ListHorizontalPodAutoscalers(client *kubernetes.Client, namespace string) (*autoscalingv1.HorizontalPodAutoscalerList, error) {
	if namespace == "" {
		namespace = "default"
	}

	result, err := client.Clientset.AutoscalingV1().HorizontalPodAutoscalers(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// deployments := make([]Deployment, len(result.Items))
	// for i, item := range result.Items {
	// 	deployments[i] = Deployment{}
	// }

	return result, nil
}
