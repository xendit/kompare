package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchWildcard(t *testing.T) {
	// Test case when the string matches the wildcard pattern
	result := matchWildcard("hello", "h*llo")
	if !result {
		t.Error("Expected true, but got false")
	}

	// Test case when the string does not match the wildcard pattern
	result = matchWildcard("world", "h*llo")
	if result {
		t.Error("Expected false, but got true")
	}

	// Test case when there is an error in pattern matching
	result = matchWildcard("hello", "[")
	if result {
		t.Error("Expected false, but got true")
	}
}

func TestDetectNamespacePattern(t *testing.T) {
	// Test case for empty pattern
	result := DetectNamespacePattern("")
	if result != "empty" {
		t.Errorf("Expected 'empty' but got '%s'", result)
	}

	// Test case for wildcard pattern
	result = DetectNamespacePattern("*.example.com")
	if result != "wildcard" {
		t.Errorf("Expected 'wildcard' but got '%s'", result)
	}

	// Test case for specific pattern
	result = DetectNamespacePattern("example.com")
	if result != "specific" {
		t.Errorf("Expected 'specific' but got '%s'", result)
	}
}

func TestMain(t *testing.T) {
	// Redirect standard output to /dev/null
	os.Stdout, _ = os.Open(os.DevNull)
	defer func() {
		os.Stdout = os.NewFile(1, "/dev/stdout")
	}()

	// Simulate a parsing error
	os.Args = []string{"main.go", "--invalid-flag"}
	assert.Panics(t, func() { main() }, "Expected a panic due to parsing error")

	// Simulate connection errors
	os.Args = []string{"main.go"}
	assert.Panics(t, func() { main() }, "Expected a panic due to connection error")

	// Simulate namespace listing errors
	// os.Args = []string{"main.go", "-s", "source-context", "-t", "target-context", "-n", "test-namespace", "-c", "invalid-path"}
	// assert.Panics(t, func() { main() }, "Expected a panic due to namespace listing error")

	_, _, kubeconfigFile := setupTestEnvironment()

	// Set up command-line arguments
	os.Args = []string{"main.go", "-t", "target-context", "-s", "source-context", "-n", "test-namespace", "-c", kubeconfigFile}

	// Capture errors from main function
	var capturedError error

	// Run the main function and capture any errors
	func() {
		defer func() {
			if r := recover(); r != nil {
				capturedError = fmt.Errorf("panic occurred: %v", r)
			}
		}()
		main()
	}()

	assert.NoError(t, capturedError, "Expected no errors during main function execution")

	// Assert that clientsetToSource and clientsetToTarget are not nil after connecting to clusters
	// assert.NotNil(t, clientsetToSource, "Expected clientsetToSource to be not nil")
	// assert.NotNil(t, clientsetToTarget, "Expected clientsetToTarget to be not nil")
	fmt.Println("Test completed")
}

func startMockCluster() (string, *http.ServeMux) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/namespaces", func(w http.ResponseWriter, r *http.Request) {
		// Return a sample namespace list or an empty list if needed
		w.Header().Set("Content-Type", "application/json") // Set content type to JSON
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"kind":"NamespaceList","items":[]}`))
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

func setupTestEnvironment() (string, string, string) {
	// Set up temporary kubeconfig file with mock cluster URLs
	sourceClusterURL, _ := startMockCluster()
	targetClusterURL, _ := startMockCluster()

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
	defer os.Remove(tempKubeconfig.Name())

	_, err = tempKubeconfig.Write(kubeconfigData)
	if err != nil {
		panic(fmt.Sprintf("Error writing kubeconfig data: %v", err))
	}

	return sourceClusterURL, targetClusterURL, tempKubeconfig.Name()
}
