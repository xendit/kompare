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

// ListNameSpaces retrieves a list of Kubernetes namespaces.
// Parameters:
// - clientset: The Kubernetes clientset used to make the API call.
// Returns:
// - (*Corev1.NamespaceList): A list of Kubernetes namespaces.
// - (error): An error if any occurred during the API call.
func ListNameSpaces(clientset *kubernetes.Clientset) (*Corev1.NamespaceList, error) {
	nsList, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return nsList, nil
}

// ListK8sDeployments retrieves a list of Kubernetes deployments in the specified namespace.
// Parameters:
// - clientset: The Kubernetes clientset used to make the API call.
// - nameSpace: The namespace in which to list the deployments. If empty, uses the "default" namespace.
// Returns:
// - (*v1.DeploymentList): A list of Kubernetes deployments.
// - (error): An error if any occurred during the API call.
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

// ListHPAs retrieves a list of HorizontalPodAutoscalers in the specified namespace.
// Parameters:
// - clientset: The Kubernetes clientset used to make the API call.
// - nameSpace: The namespace in which to list the HPAs. If empty, uses the "default" namespace.
// Returns:
// - (*autoscalingv1.HorizontalPodAutoscalerList): A list of HorizontalPodAutoscalers.
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

// ListCronJobs retrieves a list of CronJobs in the specified namespace.
// Parameters:
// - clientset: The Kubernetes clientset used to make the API call.
// - nameSpace: The namespace in which to list the CronJobs.
// Returns:
// - (*batchv1.CronJobList): A list of CronJobs.
// - (error): An error if any occurred during the API call.
func ListCronJobs(clientset *kubernetes.Clientset, nameSpace string) (*batchv1.CronJobList, error) {
	listCronJobs, err := clientset.BatchV1().CronJobs(nameSpace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return listCronJobs, nil
}

// ListIngressRoutes retrieves a list of Ingresses in the specified namespace.
// Parameters:
// - clientset: The Kubernetes clientset used to make the API call.
// - nameSpace: The namespace in which to list the Ingresses.
// Returns:
// - (*networkingv1.IngressList): A list of Ingresses.
// - (error): An error if any occurred during the API call.
func ListIngressRoutes(clientset *kubernetes.Clientset, nameSpace string) (*networkingv1.IngressList, error) {
	listIngress, err := clientset.NetworkingV1().Ingresses(nameSpace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return listIngress, nil
}

// ListServices retrieves a list of Services in the specified namespace.
// Parameters:
// - clientset: The Kubernetes clientset used to make the API call.
// - nameSpace: The namespace in which to list the Services.
// Returns:
// - (*Corev1.ServiceList): A list of Services.
// - (error): An error if any occurred during the API call.
func ListServices(clientset *kubernetes.Clientset, nameSpace string) (*Corev1.ServiceList, error) {
	listServices, err := clientset.CoreV1().Services(nameSpace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return listServices, nil
}
