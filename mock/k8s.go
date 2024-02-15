package mock

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"os"

	"github.com/gorilla/mux"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	RbacV1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func StartMockCluster() (string, *mux.Router, error) {
	r := mux.NewRouter()

	// Routes for different Kubernetes resources
	r.HandleFunc("/apis/apps/v1/namespaces/{namespace}/deployments", GetDeployments).Methods("GET")
	r.HandleFunc("/apis/networking.k8s.io/v1/namespaces/{namespace}/ingresses", GetDeployments).Methods("GET")
	r.HandleFunc("/apis/apps/v1/namespaces/{namespace}/secrets", GetSecrets).Methods("GET")
	r.HandleFunc("/api/v1/namespaces/{namespace}/configmaps", GetConfigMaps).Methods("GET")
	r.HandleFunc("/api/v1/namespaces/{namespace}/services", GetServices).Methods("GET")
	r.HandleFunc("/apis/autoscaling/v1/namespaces/{namespace}/horizontalpodautoscalers", GetHPAs).Methods("GET")
	r.HandleFunc("/apis/batch/v1/namespaces/{namespace}/cronjobs", GetCronJobs).Methods("GET")
	r.HandleFunc("/api/v1/namespaces", NamespacesHandler).Methods("GET")
	r.HandleFunc("/apis/rbac.authorization.k8s.io/v1/namespaces/{namespace}/roles", GetRoles).Methods("GET")
	r.HandleFunc("/apis/rbac.authorization.k8s.io/v1/namespaces/{namespace}/rolebindings", GetRoleBindings).Methods("GET")
	r.HandleFunc("/apis/apiextensions.k8s.io/v1/customresourcedefinitions", CustomResourceDefinitionsHandler).Methods("GET")
	r.HandleFunc("/apis/rbac.authorization.k8s.io/v1/clusterroles", ClusterRolesHandler).Methods("GET")
	r.HandleFunc("/apis/rbac.authorization.k8s.io/v1/clusterrolebindings", ClusterRoleBindingsHandler).Methods("GET")

	// Create a HTTP server instance
	server := &http.Server{
		Addr:         ":0", // Use port 0 for dynamic allocation
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Create a TCP listener
	listener, err := net.Listen("tcp", ":0") // Use port 0 for dynamic allocation
	if err != nil {
		return "", nil, fmt.Errorf("error creating listener: %v", err)
	}

	// Get the port
	port := listener.Addr().(*net.TCPAddr).Port
	clusterURL := fmt.Sprintf("http://localhost:%d", port)

	// Start serving using the HTTP server
	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error serving mock cluster: %v\n", err)
			os.Exit(1)
		}
	}()

	fmt.Println("Server is running...")
	fmt.Println("Cluster URL:", clusterURL)

	return clusterURL, r, nil
}

// NamespacesHandler handles requests to /api/v1/namespaces
func NamespacesHandler(w http.ResponseWriter, r *http.Request) {
	// Define a sample list of namespaces
	namespaces := &v1.NamespaceList{
		Items: []v1.Namespace{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "namespace1",
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "namespace2",
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "namespace3",
				},
			},
		},
	}

	// Convert the NamespaceList to JSON
	responseJSON, err := json.Marshal(namespaces)
	if err != nil {
		http.Error(w, "Error marshalling JSON response", http.StatusInternalServerError)
		return
	}

	// Set the response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Write the JSON response to the client
	_, err = w.Write(responseJSON)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error writing the JSON response: %v", err), http.StatusInternalServerError)
		return
	}
}

// CustomResourceDefinitionsHandler handles requests to /apis/apiextensions.k8s.io/v1/customresourcedefinitions
func CustomResourceDefinitionsHandler(w http.ResponseWriter, r *http.Request) {
	// Return a sample CRD list or an empty list if needed
	w.Header().Set("Content-Type", "application/json") // Set content type to JSON
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"kind":"CustomResourceDefinitionList","items":[]}`))
}

// ClusterRolesHandler handles requests to /apis/rbac.authorization.k8s.io/v1/clusterroles
func ClusterRolesHandler(w http.ResponseWriter, r *http.Request) {
	// Return a sample ClusterRole list or an empty list if needed
	w.Header().Set("Content-Type", "application/json") // Set content type to JSON
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"kind":"ClusterRoleList","items":[]}`))
}

// ClusterRoleBindingsHandler handles requests to /apis/rbac.authorization.k8s.io/v1/clusterrolebindings
func ClusterRoleBindingsHandler(w http.ResponseWriter, r *http.Request) {
	// Return a sample ClusterRoleBinding list or an empty list if needed
	w.Header().Set("Content-Type", "application/json") // Set content type to JSON
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"kind":"ClusterRoleBindingList","items":[]}`))
}

func GetDeployments(w http.ResponseWriter, r *http.Request) {
	// Create an empty DeploymentList object
	deployments := &appsv1.DeploymentList{
		ListMeta: metav1.ListMeta{
			ResourceVersion: "320850103", // Set a sample resource version
		},
		Items: []appsv1.Deployment{}, // Empty list of deployments
	}

	// Convert the DeploymentList object to JSON
	jsonResponse, err := json.Marshal(deployments)
	if err != nil {
		http.Error(w, "Error marshalling JSON response", http.StatusInternalServerError)
		return
	}

	// Set the response headers and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
func GetIngresses(w http.ResponseWriter, r *http.Request) {
	// Create an empty IngressList object
	ingresses := &networkingv1.IngressList{
		ListMeta: metav1.ListMeta{
			ResourceVersion: "320850103", // Set a sample resource version
		},
		Items: []networkingv1.Ingress{}, // Empty list of ingresses
	}

	// Convert the IngressList object to JSON
	jsonResponse, err := json.Marshal(ingresses)
	if err != nil {
		http.Error(w, "Error marshalling JSON response", http.StatusInternalServerError)
		return
	}

	// Set the response headers and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// Sample handler for Secrets
func GetSecrets(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]

	// Sample response
	response := map[string]interface{}{
		"kind":      "SecretList",
		"namespace": namespace,
		"items":     []string{}, // Add secret items if needed
	}

	// Convert response to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error marshalling JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// Sample handler for ConfigMaps
func GetConfigMaps(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]

	// Sample response
	response := map[string]interface{}{
		"kind":      "ConfigMapList",
		"namespace": namespace,
		"items":     []string{}, // Add config map items if needed
	}

	// Convert response to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error marshalling JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// Sample handler for Services
func GetServices(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]

	// Sample response
	response := map[string]interface{}{
		"kind":      "ServiceList",
		"namespace": namespace,
		"items":     []string{}, // Add service items if needed
	}

	// Convert response to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error marshalling JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// GetRoles handles HTTP requests to retrieve Role resources.
func GetRoles(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// namespace := vars["namespace"]

	// Sample response for the RoleList
	roles := &RbacV1.RoleList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "RoleList",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		Items: []RbacV1.Role{}, // Placeholder for Role items
	}

	// Convert the RoleList object to JSON
	jsonResponse, err := json.Marshal(roles)
	if err != nil {
		http.Error(w, "Error marshalling JSON response", http.StatusInternalServerError)
		return
	}

	// Set the response headers and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// GetRoleBindings handles HTTP requests to retrieve RoleBinding resources.
func GetRoleBindings(w http.ResponseWriter, r *http.Request) {

	// Sample response for the RoleBindingList
	roleBindings := &RbacV1.RoleBindingList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "RoleBindingList",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		Items: []RbacV1.RoleBinding{}, // Placeholder for RoleBinding items
	}

	// Convert the RoleBindingList object to JSON
	jsonResponse, err := json.Marshal(roleBindings)
	if err != nil {
		http.Error(w, "Error marshalling JSON response", http.StatusInternalServerError)
		return
	}

	// Set the response headers and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// GetHPAs handles HTTP requests to retrieve HorizontalPodAutoscaler resources.
func GetHPAs(w http.ResponseWriter, r *http.Request) {

	// Sample response for the HorizontalPodAutoscalerList
	hpas := &autoscalingv1.HorizontalPodAutoscalerList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HorizontalPodAutoscalerList",
			APIVersion: "autoscaling/v1",
		},
		Items: []autoscalingv1.HorizontalPodAutoscaler{}, // Placeholder for HPA items
	}

	// Convert the HorizontalPodAutoscalerList object to JSON
	jsonResponse, err := json.Marshal(hpas)
	if err != nil {
		http.Error(w, "Error marshalling JSON response", http.StatusInternalServerError)
		return
	}

	// Set the response headers and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// GetCronJobs handles HTTP requests to retrieve CronJob resources.
func GetCronJobs(w http.ResponseWriter, r *http.Request) {
	// Sample response for the CronJobList
	cronJobs := &batchv1.CronJobList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CronJobList",
			APIVersion: "batch/v1beta1",
		},
		Items: []batchv1.CronJob{}, // Placeholder for CronJob items
	}

	// Convert the CronJobList object to JSON
	jsonResponse, err := json.Marshal(cronJobs)
	if err != nil {
		http.Error(w, "Error marshalling JSON response", http.StatusInternalServerError)
		return
	}

	// Set the response headers and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func SetupTestEnvironment() (string, string, *os.File) {
	// Set up temporary kubeconfig file with mock cluster URLs
	sourceClusterURL, _, _ := StartMockCluster()
	targetClusterURL, _, _ := StartMockCluster()

	kubeconfigData := []byte(fmt.Sprintf(`
apiVersion: v1
clusters:
- cluster:
    server: %s
  name: source-context
- cluster:
    server: %s
  name: target-context
contexts:
- context:
    cluster: source-context
    user: ""
  name: source-context
- context:
    cluster: target-context
    user: ""
  name: target-context
current-context: source-context
`, sourceClusterURL, targetClusterURL))

	tempKubeconfig, err := os.CreateTemp("", "kubeconfig")
	if err != nil {
		panic(fmt.Sprintf("Error creating temporary kubeconfig: %v", err))
	}
	// defer os.Remove(tempKubeconfig.Name())

	_, err = tempKubeconfig.Write(kubeconfigData)
	if err != nil {
		panic(fmt.Sprintf("Error writing kubeconfig data: %v", err))
	}

	return sourceClusterURL, targetClusterURL, tempKubeconfig
}
