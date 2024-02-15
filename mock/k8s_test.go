package mock

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestStartMockCluster(t *testing.T) {
	clusterURL, mux, _ := StartMockCluster()

	// No need to defer a Close call for http.ServeMux

	// Assert that clusterURL is not empty
	if clusterURL == "" {
		t.Error("Expected clusterURL to be non-empty, got empty string")
	}

	// Test the /api/v1/namespaces endpoint
	req := httptest.NewRequest("GET", "/api/v1/namespaces", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	expectedContentType := "application/json"
	if contentType := resp.Header.Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("Expected content type %s, got %s", expectedContentType, contentType)
	}
	// Test the /apis/apps/v1/namespaces/{namespace}/deployments endpoint
	deploymentsReq := httptest.NewRequest("GET", "/apis/apps/v1/namespaces/{namespace}/deployments", nil)
	deploymentsW := httptest.NewRecorder()
	mux.ServeHTTP(deploymentsW, deploymentsReq)
	deploymentsResp := deploymentsW.Result()

	if deploymentsResp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, deploymentsResp.StatusCode)
	}

	if contentType := deploymentsResp.Header.Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("Expected content type %s, got %s", expectedContentType, contentType)
	}
}

func TestSetupTestEnvironment(t *testing.T) {
	sourceClusterURL, targetClusterURL, tempKubeconfig := SetupTestEnvironment()
	defer func() {
		err := tempKubeconfig.Close()
		if err != nil {
			t.Errorf("Error closing tempKubeconfig: %v", err)
		}
	}()

	// Check if the cluster URLs are not empty
	if sourceClusterURL == "" || targetClusterURL == "" {
		t.Error("Cluster URLs should not be empty")
	}

	// Check if tempKubeconfig is not nil
	if tempKubeconfig == nil {
		t.Error("Temp kubeconfig file should not be nil")
	}

	// Read the contents of tempKubeconfig and check if it's not empty
	// kubeconfigData := make([]byte, 1024)
	n, err := os.ReadFile(tempKubeconfig.Name())
	if err != nil {
		t.Errorf("Error reading kubeconfig file: %v", err)
	}

	// if n == 0 {
	// 	t.Error("Kubeconfig file should not be empty")
	// }

	// Ensure that the kubeconfig contains the expected cluster URLs
	expectedKubeconfigContent := fmt.Sprintf(`
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
`, sourceClusterURL, targetClusterURL)

	if string(n) != expectedKubeconfigContent {
		t.Errorf("Expected kubeconfig content does not match\nExpected:\n%s\nGot:\n%s", expectedKubeconfigContent, string(n))
	}
}
