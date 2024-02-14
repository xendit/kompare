package mock

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"

	Corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func StartMockCluster() (string, *http.ServeMux) {
	mux := http.NewServeMux()
	// Handle the /api/v1/namespaces endpoint
	mux.HandleFunc("/api/v1/namespaces", func(w http.ResponseWriter, r *http.Request) {
		// Define a sample list of namespaces
		namespaces := []Corev1.Namespace{
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
		}

		// Create a NamespaceList response structure
		namespaceList := struct {
			Kind  string             `json:"kind"`
			Items []Corev1.Namespace `json:"items"`
		}{
			Kind:  "NamespaceList",
			Items: namespaces,
		}
		// Convert the response structure to JSON
		responseJSON, err := json.Marshal(namespaceList)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error marshalling JSON response: %v", err), http.StatusInternalServerError)
			return
		}

		// Set the response headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Write the JSON response to the client
		w.Write(responseJSON)
	})

	// Handle customresourcedefinitions endpoint
	mux.HandleFunc("/apis/apiextensions.k8s.io/v1/customresourcedefinitions", func(w http.ResponseWriter, r *http.Request) {
		// Return a sample CRD list or an empty list if needed
		w.Header().Set("Content-Type", "application/json") // Set content type to JSON
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"kind":"CustomResourceDefinitionList","items":[]}`))
	})

	// Handle Deployment resource
	mux.HandleFunc("/apis/apps/v1/namespaces/test-namespace/deployments", func(w http.ResponseWriter, r *http.Request) {
		// Return a sample deployment list or an empty list if needed
		w.Header().Set("Content-Type", "application/json") // Set content type to JSON
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"kind":"DeploymentList","items":[]}`))
	})

	// Handle Ingress resource
	mux.HandleFunc("/apis/networking.k8s.io/v1/namespaces/test-namespace/ingresses", func(w http.ResponseWriter, r *http.Request) {
		// Return a sample ingress list or an empty list if needed
		w.Header().Set("Content-Type", "application/json") // Set content type to JSON
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"kind":"IngressList","items":[]}`))
	})

	// Handle Service resource
	mux.HandleFunc("/api/v1/namespaces/test-namespace/services", func(w http.ResponseWriter, r *http.Request) {
		// Return a sample service list or an empty list if needed
		w.Header().Set("Content-Type", "application/json") // Set content type to JSON
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"kind":"ServiceList","items":[]}`))
	})

	// Handle ServiceAccount resource
	mux.HandleFunc("/api/v1/namespaces/test-namespace/serviceaccounts", func(w http.ResponseWriter, r *http.Request) {
		// Return a sample service account list or an empty list if needed
		w.Header().Set("Content-Type", "application/json") // Set content type to JSON
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"kind":"ServiceAccountList","items":[]}`))
	})

	// Handle ConfigMap resource
	mux.HandleFunc("/api/v1/namespaces/test-namespace/configmaps", func(w http.ResponseWriter, r *http.Request) {
		// Return a sample config map list or an empty list if needed
		w.Header().Set("Content-Type", "application/json") // Set content type to JSON
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"kind":"ConfigMapList","items":[]}`))
	})

	// Handle Secret resource
	mux.HandleFunc("/api/v1/namespaces/test-namespace/secrets", func(w http.ResponseWriter, r *http.Request) {
		// Return a sample secret list or an empty list if needed
		w.Header().Set("Content-Type", "application/json") // Set content type to JSON
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"kind":"SecretList","items":[]}`))
	})

	// Handle Role resource
	mux.HandleFunc("/apis/rbac.authorization.k8s.io/v1/namespaces/test-namespace/roles", func(w http.ResponseWriter, r *http.Request) {
		// Return a sample role list or an empty list if needed
		w.Header().Set("Content-Type", "application/json") // Set content type to JSON
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"kind":"RoleList","items":[]}`))
	})

	// Handle RoleBinding resource
	mux.HandleFunc("/apis/rbac.authorization.k8s.io/v1/namespaces/test-namespace/rolebindings", func(w http.ResponseWriter, r *http.Request) {
		// Return a sample role binding list or an empty list if needed
		w.Header().Set("Content-Type", "application/json") // Set content type to JSON
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"kind":"RoleBindingList","items":[]}`))
	})

	// Handle HorizontalPodAutoscaler resource
	mux.HandleFunc("/apis/autoscaling/v1/namespaces/test-namespace/horizontalpodautoscalers", func(w http.ResponseWriter, r *http.Request) {
		// Return a sample HPA list or an empty list if needed
		w.Header().Set("Content-Type", "application/json") // Set content type to JSON
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"kind":"HorizontalPodAutoscalerList","items":[]}`))
	})

	// Handle CronJob resource
	mux.HandleFunc("/apis/batch/v1/namespaces/test-namespace/cronjobs", func(w http.ResponseWriter, r *http.Request) {
		// Return a sample cron job list or an empty list if needed
		w.Header().Set("Content-Type", "application/json") // Set content type to JSON
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"kind":"CronJobList","items":[]}`))
	})

	server := &http.Server{Handler: mux}

	listener, err := net.Listen("tcp", "localhost:0") // Listen on any available port
	if err != nil {
		fmt.Printf("Error starting mock cluster: %v\n", err)
		os.Exit(1)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	clusterURL := fmt.Sprintf("http://localhost:%d", port)

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error serving mock cluster: %v\n", err)
			os.Exit(1)
		}
	}()

	return clusterURL, mux
}

func SetupTestEnvironment() (string, string, *os.File) {
	// Set up temporary kubeconfig file with mock cluster URLs
	sourceClusterURL, _ := StartMockCluster()
	targetClusterURL, _ := StartMockCluster()

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
