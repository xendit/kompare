package query

import (
	"context"
	"fmt"
	"kompare/connect"

	v1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	batchv1 "k8s.io/api/batch/v1"

	networkingv1 "k8s.io/api/networking/v1"
	// traefikv1alpha1 "github.com/traefik/traefik/v3/pkg/provider/kubernetes/crd/traefikio/v1alpha1"
	apiextensionv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"

	Corev1 "k8s.io/api/core/v1"
	RbacV1 "k8s.io/api/rbac/v1"
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
		return nil, fmt.Errorf("failed to query the NameSpace List: %w", err)
	}
	return nsList, nil
}

func GetNamespace(clientset *kubernetes.Clientset, name string) (*Corev1.Namespace, error) {
	ns, err := clientset.CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to query the NameSpace Get: %w", err)
	}
	return ns, nil
}

// ListK8sDeployments retrieves a list of Kubernetes deployments in the specified namespace.
// Parameters:
// - clientset: The Kubernetes clientset used to make the API call.
// - nameSpace: The namespace in which to list the deployments. If empty, uses the "default" namespace.
// Returns:
// - (*v1.DeploymentList): A list of Kubernetes deployments.
// - (error): An error if any occurred during the API call.
func ListDeployments(clientset *kubernetes.Clientset, nameSpace string) (*v1.DeploymentList, error) {
	if nameSpace == "" {
		nameSpace = "default"
	}
	deployments_list, err := clientset.AppsV1().Deployments(nameSpace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to query the Deployment List: %w", err)
	}
	return deployments_list, nil
}

// ListHPAs retrieves a list of HorizontalPodAutoscalers in the specified namespace.
// Parameters:
// - clientset: The Kubernetes clientset used to make the API call.
// - nameSpace: The namespace in which to list the HPAs. If empty, uses the "default" namespace.
// Returns:
// - (*autoscalingv1.HorizontalPodAutoscalerList): A list of HorizontalPodAutoscalers.
func ListHPAs(clientset *kubernetes.Clientset, nameSpace string) (*autoscalingv1.HorizontalPodAutoscalerList, error) {
	if nameSpace == "" {
		nameSpace = "default"
	}
	listHPA, err := clientset.AutoscalingV1().HorizontalPodAutoscalers(nameSpace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to query the HPA List: %w", err)
	}
	return listHPA, nil
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
		return nil, fmt.Errorf("failed to query the Cron Job List: %w", err)
	}
	return listCronJobs, nil
}

// ListCRDs retrieves a list of CRDs in the specified namespace.
// Parameters:
// - kubeconfig: The Kubernetes client config used to make the API call.
// - nameSpace: The namespace in which to list the Ingresses.
// Returns:
// - (*apiextensionv1.CustomResourceDefinitionList): A list of CRDs.
// - (error): An error if any occurred during the API call.
func ListCRDs(ctx, kubeconfig string) (*apiextensionv1.CustomResourceDefinitionList, error) {
	config, err := connect.BuildConfigWithContextFromFlags(ctx, kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create config for quering CRDs: %w", err)
	}
	kubeClient, err := apiextension.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to get clientset for quering CRDs: %w", err)
	}

	listCRDs, err := kubeClient.ApiextensionsV1().CustomResourceDefinitions().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to query the CRD List: %w", err)
	}
	return listCRDs, nil
}

// ListIngresses retrieves a list of Ingresses in the specified namespace.
// Parameters:
// - clientset: The Kubernetes clientset used to make the API call.
// - nameSpace: The namespace in which to list the Services.
// Returns:
// - (*networkingv1.IngressList): A list of Services.
// - (error): An error if any occurred during the API call.
func ListIngresses(clientset *kubernetes.Clientset, nameSpace string) (*networkingv1.IngressList, error) {
	listIngress, err := clientset.NetworkingV1().Ingresses(nameSpace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to query the Ingress List: %w", err)
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
		return nil, fmt.Errorf("failed to query the Kubernetes Services List: %w", err)
	}
	return listServices, nil
}

// Get Config Maps list.
func ListConfigMaps(clientset *kubernetes.Clientset, nameSpace string) (*Corev1.ConfigMapList, error) {
	ListConfigMaps, err := clientset.CoreV1().ConfigMaps(nameSpace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to query the Config Map List: %w", err)
	}
	return ListConfigMaps, nil
}

// Get Secrets list.
func ListSecrets(clientset *kubernetes.Clientset, nameSpace string) (*Corev1.SecretList, error) {
	ListSercrets, err := clientset.CoreV1().Secrets(nameSpace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to query the Secrets List: %w", err)
	}
	return ListSercrets, nil
}

// Get Service Accounts list.
func ListServiceAccounts(clientset *kubernetes.Clientset, nameSpace string) (*Corev1.ServiceAccountList, error) {
	ListServiceAccounts, err := clientset.CoreV1().ServiceAccounts(nameSpace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to query the Service Accounts List: %w", err)
	}
	return ListServiceAccounts, nil
}

// roles
func ListRoles(clientset *kubernetes.Clientset, nameSpace string) (*RbacV1.RoleList, error) {
	Listroles, err := clientset.RbacV1().Roles(nameSpace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to query the Roles List: %w", err)
	}
	return Listroles, nil
}

// rolebindings
func ListRoleBindings(clientset *kubernetes.Clientset, nameSpace string) (*RbacV1.RoleBindingList, error) {
	ListRoleBindings, err := clientset.RbacV1().RoleBindings(nameSpace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to query the Role bindings List: %w", err)
	}
	return ListRoleBindings, nil
}

// clusterroles
func ListClusterRoles(clientset *kubernetes.Clientset) (*RbacV1.ClusterRoleList, error) {
	ListClusterRoles, err := clientset.RbacV1().ClusterRoles().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to query the Cluster Roles List: %w", err)
	}
	return ListClusterRoles, nil
}

// clusterrolebindings
func ListClusterRoleBindings(clientset *kubernetes.Clientset) (*RbacV1.ClusterRoleBindingList, error) {
	ListClusterRoleBindings, err := clientset.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to query the Cluster Role Bindings List: %w", err)
	}
	return ListClusterRoleBindings, nil
}

// ListNetworkPolicy retrieves a list of NetworkPolicies in the specified namespace.
func ListNetworkPolicies(clientset *kubernetes.Clientset, nameSpace string) (*networkingv1.NetworkPolicyList, error) {
	listNetworkPolicy, err := clientset.NetworkingV1().NetworkPolicies(nameSpace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to query the Network Policy List: %w", err)
	}
	return listNetworkPolicy, nil
}
