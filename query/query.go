package query

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	batchv1 "k8s.io/api/batch/v1"
	networkingv1 "k8s.io/api/networking/v1"

	Corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ListNameSpaces(clientset *kubernetes.Clientset) (*Corev1.NamespaceList, error) {
	nsList, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return nsList, nil
}

func ListK8sDeployments(clientset *kubernetes.Clientset, nameSpace string) (*v1.DeploymentList, error) {
	if nameSpace == "" {
		nameSpace = "default"
	}
	deployments_list, err := clientset.AppsV1().Deployments(nameSpace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return deployments_list, nil
}

func ListHPAs(clientset *kubernetes.Clientset, nameSpace string) *autoscalingv1.HorizontalPodAutoscalerList {
	if nameSpace == "" {
		nameSpace = "default"
	}
	listHPA, err := clientset.AutoscalingV1().HorizontalPodAutoscalers(nameSpace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	return listHPA
}

func ListCronJobs(clientset *kubernetes.Clientset, nameSpace string) (*batchv1.CronJobList, error) {
	listCronJobs, err := clientset.BatchV1().CronJobs(nameSpace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return listCronJobs, nil
}

func ListIngressRoutes(clientset *kubernetes.Clientset, nameSpace string) (*networkingv1.IngressList, error) {
	listIngress, err := clientset.NetworkingV1().Ingresses(nameSpace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return listIngress, nil
}

func ListServices(clientset *kubernetes.Clientset, nameSpace string) (*Corev1.ServiceList, error) {
	listServices, err := clientset.CoreV1().Services(nameSpace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return listServices, nil
}
