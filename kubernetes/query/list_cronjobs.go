package query

import (
	"context"
	"kompare/kubernetes"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ListCronjobs(client *kubernetes.Client, namespace string) (*batchv1.CronJobList, error) {
	if namespace == "" {
		namespace = "default"
	}

	result, err := client.Clientset.BatchV1().CronJobs(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return result, nil
}
