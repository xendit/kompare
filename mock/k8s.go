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
	Corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	RbacV1 "k8s.io/api/rbac/v1"
	apiextensionv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func StartMockCluster() (string, *mux.Router, error) {
	r := mux.NewRouter()

	// Routes for different Kubernetes resources
	r.HandleFunc("/apis/apps/v1/namespaces/{namespace}/deployments", GetDeployments).Methods("GET")
	r.HandleFunc("/apis/networking.k8s.io/v1/namespaces/{namespace}/ingresses", GetIngresses).Methods("GET")
	r.HandleFunc("/api/v1/namespaces/{namespace}/secrets", GetSecrets).Methods("GET")
	r.HandleFunc("/api/v1/namespaces/{namespace}/configmaps", GetConfigMaps).Methods("GET")
	r.HandleFunc("/api/v1/namespaces/{namespace}/services", GetServices).Methods("GET")
	r.HandleFunc("/apis/autoscaling/v1/namespaces/{namespace}/horizontalpodautoscalers", GetHPAs).Methods("GET")
	r.HandleFunc("/apis/batch/v1/namespaces/{namespace}/cronjobs", GetCronJobs).Methods("GET")
	r.HandleFunc("/api/v1/namespaces", GetNamespaces).Methods("GET")
	r.HandleFunc("/api/v1/namespaces/namespace2", GetNamespace).Methods("GET")
	r.HandleFunc("/apis/rbac.authorization.k8s.io/v1/namespaces/{namespace}/roles", GetRoles).Methods("GET")
	r.HandleFunc("/apis/rbac.authorization.k8s.io/v1/namespaces/{namespace}/rolebindings", GetRoleBindings).Methods("GET")
	r.HandleFunc("/apis/apiextensions.k8s.io/v1/customresourcedefinitions", GetCustomResourceDefinitions).Methods("GET")
	r.HandleFunc("/apis/rbac.authorization.k8s.io/v1/clusterroles", GetClusterRoles).Methods("GET")
	r.HandleFunc("/apis/rbac.authorization.k8s.io/v1/clusterrolebindings", GetClusterRoleBindings).Methods("GET")
	r.HandleFunc("/api/v1/namespaces/{namespace}/serviceaccounts", GetServiceAccounts).Methods("GET")

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

// GetNamespace handles requests to /api/v1/namespaces/{name}
func GetNamespace(w http.ResponseWriter, r *http.Request) {
	// Define the namespace
	namespace := v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "namespace2",
		},
	}

	// Marshal namespace to JSON
	namespaceJSON, err := json.Marshal(namespace)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshaling namespace to JSON: %v", err), http.StatusInternalServerError)
		return
	}

	// Set Content-Type header
	w.Header().Set("Content-Type", "application/json")

	// Write JSON response
	_, err = w.Write(namespaceJSON)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error writing the JSON response: %v", err), http.StatusInternalServerError)
		return
	}
}

// NamespacesHandler handles requests to /api/v1/namespaces
func GetNamespaces(w http.ResponseWriter, r *http.Request) {
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

func GetDeployments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]

	// Check if the namespace is "namespace2"
	if namespace == "namespace2" {
		// Create three sample deployments for namespace2
		deployments := &appsv1.DeploymentList{
			ListMeta: metav1.ListMeta{
				ResourceVersion: "320850103", // Set a sample resource version
			},
			Items: []appsv1.Deployment{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "deployment1",
						Namespace: "namespace2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "deployment2",
						Namespace: "namespace2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "deployment3",
						Namespace: "namespace2",
					},
				},
			},
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
		_, err = w.Write(jsonResponse)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error writing the JSON response: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		// If the namespace is not "namespace2", return an empty list of deployments
		deployments := &appsv1.DeploymentList{
			ListMeta: metav1.ListMeta{
				ResourceVersion: "0",
			},
			Items: []appsv1.Deployment{},
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
		_, err = w.Write(jsonResponse)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error writing the JSON response: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

// GetHPAs handles HTTP requests to retrieve HorizontalPodAutoscaler resources.
func GetHPAs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]

	// Check if the namespace is "namespace2"
	if namespace == "namespace2" {
		// Create three sample HPAs for namespace2
		hpas := &autoscalingv1.HorizontalPodAutoscalerList{
			ListMeta: metav1.ListMeta{
				ResourceVersion: "320850103", // Set a sample resource version
			},
			Items: []autoscalingv1.HorizontalPodAutoscaler{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "hpa1",
						Namespace: "namespace2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "hpa2",
						Namespace: "namespace2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "hpa3",
						Namespace: "namespace2",
					},
				},
			},
		}

		// Convert the HPAList object to JSON
		jsonResponse, err := json.Marshal(hpas)
		if err != nil {
			http.Error(w, "Error marshalling JSON response", http.StatusInternalServerError)
			return
		}

		// Set the response headers and write the JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonResponse)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error writing the JSON response: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		// If the namespace is not "namespace2", return an empty list of HPAs
		hpas := &autoscalingv1.HorizontalPodAutoscalerList{
			ListMeta: metav1.ListMeta{
				ResourceVersion: "0",
			},
			Items: []autoscalingv1.HorizontalPodAutoscaler{},
		}

		// Convert the HPAList object to JSON
		jsonResponse, err := json.Marshal(hpas)
		if err != nil {
			http.Error(w, "Error marshalling JSON response", http.StatusInternalServerError)
			return
		}

		// Set the response headers and write the JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonResponse)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error writing the JSON response: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

// GetCronJobs handles HTTP requests to retrieve CronJob resources.
func GetCronJobs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]

	// Check if the namespace is "namespace2"
	if namespace == "namespace2" {
		// Create three sample CronJobs for namespace2
		cronJobs := &batchv1.CronJobList{
			TypeMeta: metav1.TypeMeta{
				Kind:       "CronJobList",
				APIVersion: "batch/v1",
			},
			Items: []batchv1.CronJob{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "cronjob1",
						Namespace: "namespace2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "cronjob2",
						Namespace: "namespace2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "cronjob3",
						Namespace: "namespace2",
					},
				},
			},
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
		_, err = w.Write(jsonResponse)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error writing the JSON response: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		// If the namespace is not "namespace2", return an empty list of CronJobs
		cronJobs := &batchv1.CronJobList{
			TypeMeta: metav1.TypeMeta{
				Kind:       "CronJobList",
				APIVersion: "batch/v1",
			},
			Items: []batchv1.CronJob{},
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
		_, err = w.Write(jsonResponse)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error writing the JSON response: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

// GetCustomResourceDefinitions handles HTTP requests to retrieve Custom Resource Definitions (CRDs).
func GetCustomResourceDefinitions(w http.ResponseWriter, r *http.Request) {
	// Sample response for the CRDList
	crdList := &apiextensionv1.CustomResourceDefinitionList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CustomResourceDefinitionList",
			APIVersion: "apiextensions.k8s.io/v1",
		},
		Items: []apiextensionv1.CustomResourceDefinition{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "crd1",
				},
				// Add more details for CRD1 as needed
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "crd2",
				},
				// Add more details for CRD2 as needed
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "crd3",
				},
				// Add more details for CRD3 as needed
			},
		},
	}

	// Convert the CRDList object to JSON
	jsonResponse, err := json.Marshal(crdList)
	if err != nil {
		http.Error(w, "Error marshalling JSON response", http.StatusInternalServerError)
		return
	}

	// Set the response headers and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error writing the JSON response: %v", err), http.StatusInternalServerError)
		return
	}
}

func GetIngresses(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]

	// Check if the namespace is "namespace2"
	if namespace == "namespace2" {
		// Create three sample ingresses for namespace2
		ingresses := &networkingv1.IngressList{
			ListMeta: metav1.ListMeta{
				ResourceVersion: "320850103", // Set a sample resource version
			},
			Items: []networkingv1.Ingress{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "ingress1",
						Namespace: "namespace2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "ingress2",
						Namespace: "namespace2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "ingress3",
						Namespace: "namespace2",
					},
				},
			},
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
		_, err = w.Write(jsonResponse)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error writing the JSON response: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		// If the namespace is not "namespace2", return an empty list of ingresses
		ingresses := &networkingv1.IngressList{
			ListMeta: metav1.ListMeta{
				ResourceVersion: "0",
			},
			Items: []networkingv1.Ingress{},
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
		_, err = w.Write(jsonResponse)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error writing the JSON response: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

// GetSecrets handles HTTP requests to retrieve Secrets.
func GetSecrets(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]

	// Check if the namespace is "namespace2"
	if namespace == "namespace2" {
		// Create three sample secrets for namespace2
		secrets := &Corev1.SecretList{
			ListMeta: metav1.ListMeta{
				ResourceVersion: "320850103", // Set a sample resource version
			},
			Items: []Corev1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "secret1",
						Namespace: "namespace2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "secret2",
						Namespace: "namespace2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "secret3",
						Namespace: "namespace2",
					},
				},
			},
		}

		// Convert the SecretList object to JSON
		jsonResponse, err := json.Marshal(secrets)
		if err != nil {
			http.Error(w, "Error marshalling JSON response", http.StatusInternalServerError)
			return
		}

		// Set the response headers and write the JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonResponse)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error writing the JSON response: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		// If the namespace is not "namespace2", return an empty list of secrets
		secrets := &Corev1.SecretList{
			ListMeta: metav1.ListMeta{
				ResourceVersion: "0",
			},
			Items: []Corev1.Secret{},
		}

		// Convert the SecretList object to JSON
		jsonResponse, err := json.Marshal(secrets)
		if err != nil {
			http.Error(w, "Error marshalling JSON response", http.StatusInternalServerError)
			return
		}

		// Set the response headers and write the JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonResponse)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error writing the JSON response: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

// GetConfigMaps handles HTTP requests to retrieve ConfigMaps.
func GetConfigMaps(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]

	// Check if the namespace is "namespace2"
	if namespace == "namespace2" {
		// Create three sample config maps for namespace2
		configMaps := &Corev1.ConfigMapList{
			ListMeta: metav1.ListMeta{
				ResourceVersion: "320850103", // Set a sample resource version
			},
			Items: []Corev1.ConfigMap{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "configmap1",
						Namespace: "namespace2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "configmap2",
						Namespace: "namespace2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "configmap3",
						Namespace: "namespace2",
					},
				},
			},
		}

		// Convert the ConfigMapList object to JSON
		jsonResponse, err := json.Marshal(configMaps)
		if err != nil {
			http.Error(w, "Error marshalling JSON response", http.StatusInternalServerError)
			return
		}

		// Set the response headers and write the JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonResponse)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error writing the JSON response: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		// If the namespace is not "namespace2", return an empty list of config maps
		configMaps := &Corev1.ConfigMapList{
			ListMeta: metav1.ListMeta{
				ResourceVersion: "0",
			},
			Items: []Corev1.ConfigMap{},
		}

		// Convert the ConfigMapList object to JSON
		jsonResponse, err := json.Marshal(configMaps)
		if err != nil {
			http.Error(w, "Error marshalling JSON response", http.StatusInternalServerError)
			return
		}

		// Set the response headers and write the JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonResponse)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error writing the JSON response: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

func GetServices(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]

	// Check if the namespace is "namespace2"
	if namespace == "namespace2" {
		// Create three sample services for namespace2
		services := &Corev1.ServiceList{
			ListMeta: metav1.ListMeta{
				ResourceVersion: "320850103", // Set a sample resource version
			},
			Items: []Corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service1",
						Namespace: "namespace2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service2",
						Namespace: "namespace2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service3",
						Namespace: "namespace2",
					},
				},
			},
		}

		// Convert the ServiceList object to JSON
		jsonResponse, err := json.Marshal(services)
		if err != nil {
			http.Error(w, "Error marshalling JSON response", http.StatusInternalServerError)
			return
		}

		// Set the response headers and write the JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonResponse)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error writing the JSON response: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		// If the namespace is not "namespace2", return an empty list of services
		services := &Corev1.ServiceList{
			ListMeta: metav1.ListMeta{
				ResourceVersion: "0",
			},
			Items: []Corev1.Service{},
		}

		// Convert the ServiceList object to JSON
		jsonResponse, err := json.Marshal(services)
		if err != nil {
			http.Error(w, "Error marshalling JSON response", http.StatusInternalServerError)
			return
		}

		// Set the response headers and write the JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonResponse)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error writing the JSON response: %v", err), http.StatusInternalServerError)
			return
		}
	}
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

// GetServiceAccounts handles HTTP requests to retrieve Service Accounts.
func GetServiceAccounts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]

	// Check if the namespace is "namespace2"
	if namespace == "namespace2" {
		// Create three sample service accounts for namespace2
		serviceAccounts := &Corev1.ServiceAccountList{
			ListMeta: metav1.ListMeta{
				ResourceVersion: "320850103", // Set a sample resource version
			},
			Items: []Corev1.ServiceAccount{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "serviceaccount1",
						Namespace: "namespace2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "serviceaccount2",
						Namespace: "namespace2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "serviceaccount3",
						Namespace: "namespace2",
					},
				},
			},
		}

		// Convert the ServiceAccountList object to JSON
		jsonResponse, err := json.Marshal(serviceAccounts)
		if err != nil {
			http.Error(w, "Error marshalling JSON response", http.StatusInternalServerError)
			return
		}

		// Set the response headers and write the JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonResponse)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error writing the JSON response: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		// If the namespace is not "namespace2", return an empty list of service accounts
		serviceAccounts := &Corev1.ServiceAccountList{
			ListMeta: metav1.ListMeta{
				ResourceVersion: "0",
			},
			Items: []Corev1.ServiceAccount{},
		}

		// Convert the ServiceAccountList object to JSON
		jsonResponse, err := json.Marshal(serviceAccounts)
		if err != nil {
			http.Error(w, "Error marshalling JSON response", http.StatusInternalServerError)
			return
		}

		// Set the response headers and write the JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonResponse)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error writing the JSON response: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

// GetRoles handles HTTP requests to retrieve Roles.
func GetRoles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]

	// Check if the namespace is "namespace2"
	if namespace == "namespace2" {
		// Create three sample roles for namespace2
		roles := &RbacV1.RoleList{
			ListMeta: metav1.ListMeta{
				ResourceVersion: "320850103", // Set a sample resource version
			},
			Items: []RbacV1.Role{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "role1",
						Namespace: "namespace2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "role2",
						Namespace: "namespace2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "role3",
						Namespace: "namespace2",
					},
				},
			},
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
		_, err = w.Write(jsonResponse)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error writing the JSON response: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		// If the namespace is not "namespace2", return an empty list of roles
		roles := &RbacV1.RoleList{
			ListMeta: metav1.ListMeta{
				ResourceVersion: "0",
			},
			Items: []RbacV1.Role{},
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
		_, err = w.Write(jsonResponse)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error writing the JSON response: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

// GetRoleBindings handles HTTP requests to retrieve RoleBinding resources.
func GetRoleBindings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]

	// Check if the namespace is "namespace2"
	if namespace == "namespace2" {
		// Create three sample role bindings for namespace2
		roleBindings := &RbacV1.RoleBindingList{
			TypeMeta: metav1.TypeMeta{
				Kind:       "RoleBindingList",
				APIVersion: "rbac.authorization.k8s.io/v1",
			},
			Items: []RbacV1.RoleBinding{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "rolebinding1",
						Namespace: "namespace2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "rolebinding2",
						Namespace: "namespace2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "rolebinding3",
						Namespace: "namespace2",
					},
				},
			},
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
		_, err = w.Write(jsonResponse)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error writing the JSON response: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
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
		_, err = w.Write(jsonResponse)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error writing the JSON response: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

// GetClusterRoles handles HTTP requests to retrieve ClusterRoles.
func GetClusterRoles(w http.ResponseWriter, r *http.Request) {
	// Create three sample cluster roles
	clusterRoles := &RbacV1.ClusterRoleList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRoleList",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		Items: []RbacV1.ClusterRole{
			{
				// ClusterRole 1
				TypeMeta: metav1.TypeMeta{
					Kind:       "ClusterRole",
					APIVersion: "rbac.authorization.k8s.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:            "clusterroleNAMEWITHNUMBER1",
					ResourceVersion: "320850103", // Same as the previous resource version
				},
				// Add rules, etc. for ClusterRole 1
			},
			{
				// ClusterRole 2
				TypeMeta: metav1.TypeMeta{
					Kind:       "ClusterRole",
					APIVersion: "rbac.authorization.k8s.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:            "clusterroleNAMEWITHNUMBER2",
					ResourceVersion: "320850103", // Same as the previous resource version
				},
				// Add rules, etc. for ClusterRole 2
			},
			{
				// ClusterRole 3
				TypeMeta: metav1.TypeMeta{
					Kind:       "ClusterRole",
					APIVersion: "rbac.authorization.k8s.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:            "clusterroleNAMEWITHNUMBER3",
					ResourceVersion: "320850103", // Same as the previous resource version
				},
				// Add rules, etc. for ClusterRole 3
			},
		},
	}

	// Convert the ClusterRoleList object to JSON
	jsonResponse, err := json.Marshal(clusterRoles)
	if err != nil {
		http.Error(w, "Error marshalling JSON response", http.StatusInternalServerError)
		return
	}

	// Set the response headers and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error writing the JSON response: %v", err), http.StatusInternalServerError)
		return
	}
}

// GetClusterRoleBindings handles HTTP requests to retrieve ClusterRoleBinding resources.
func GetClusterRoleBindings(w http.ResponseWriter, r *http.Request) {
	// Sample response for the ClusterRoleBindingList
	clusterRoleBindings := &RbacV1.ClusterRoleBindingList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRoleBindingList",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		Items: []RbacV1.ClusterRoleBinding{
			{
				// ClusterRoleBinding 1
				TypeMeta: metav1.TypeMeta{
					Kind:       "ClusterRoleBinding",
					APIVersion: "rbac.authorization.k8s.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:            "clusterrolebinding1",
					ResourceVersion: "320850103", // Same as the previous resource version
				},
				// Add roleRef, subjects, etc. for ClusterRoleBinding 1
			},
			{
				// ClusterRoleBinding 2
				TypeMeta: metav1.TypeMeta{
					Kind:       "ClusterRoleBinding",
					APIVersion: "rbac.authorization.k8s.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:            "clusterrolebinding2",
					ResourceVersion: "320850103", // Same as the previous resource version
				},
				// Add roleRef, subjects, etc. for ClusterRoleBinding 2
			},
			{
				// ClusterRoleBinding 3
				TypeMeta: metav1.TypeMeta{
					Kind:       "ClusterRoleBinding",
					APIVersion: "rbac.authorization.k8s.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:            "clusterrolebinding3",
					ResourceVersion: "320850103", // Same as the previous resource version
				},
				// Add roleRef, subjects, etc. for ClusterRoleBinding 3
			},
		},
	}

	// Convert the ClusterRoleBindingList object to JSON
	jsonResponse, err := json.Marshal(clusterRoleBindings)
	if err != nil {
		http.Error(w, "Error marshalling JSON response", http.StatusInternalServerError)
		return
	}

	// Set the response headers and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error writing the JSON response: %v", err), http.StatusInternalServerError)
		return
	}
}
