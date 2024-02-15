package query

import (
	"os"
	"testing"

	"kompare/connect"
	"kompare/mock"

	"k8s.io/client-go/tools/clientcmd"
)

func TestListNameSpaces(t *testing.T) {
	// Set up test environment and get the temporary kubeconfig file
	_, _, tempKubeconfig := mock.SetupTestEnvironment()
	defer tempKubeconfig.Close() // Close the file after the test completes
	// Read the content of the kubeconfig file
	tempKubeconfigByte, err := os.ReadFile(tempKubeconfig.Name())
	if err != nil {
		t.Fatalf("Error reading kubeconfig file: %v", err)
	}
	// Load the kubeconfig data
	kubeconfig, err := clientcmd.Load(tempKubeconfigByte)
	if err != nil {
		t.Fatalf("Error loading kubeconfig: %v", err)
	}

	// Choose one of the contexts from the kubeconfig
	var testContext string
	for context := range kubeconfig.Contexts {
		testContext = context
		break // Choose the first context, you can modify this logic as needed
	}
	x := tempKubeconfig.Name()
	config, err := connect.ConnectToSource(testContext, &x)
	if err != nil {
		t.Fatalf("Error creating config: %v", err)
	}

	namespaces, err := ListNameSpaces(config)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	expectedLength := 3
	if len(namespaces.Items) != expectedLength {
		t.Errorf("Expected %d namespaces, got: %d", expectedLength, len(namespaces.Items))
	}
}

// func TestGetNamespace(t *testing.T) {
// 	// Set up test environment and get the temporary kubeconfig file
// 	_, _, tempKubeconfig := mock.SetupTestEnvironment()
// 	defer tempKubeconfig.Close() // Close the file after the test completes
// 	// Read the content of the kubeconfig file
// 	tempKubeconfigByte, err := os.ReadFile(tempKubeconfig.Name())
// 	if err != nil {
// 		t.Fatalf("Error reading kubeconfig file: %v", err)
// 	}
// 	// Load the kubeconfig data
// 	kubeconfig, err := clientcmd.Load(tempKubeconfigByte)
// 	if err != nil {
// 		t.Fatalf("Error loading kubeconfig: %v", err)
// 	}

// 	// Choose one of the contexts from the kubeconfig
// 	var testContext string
// 	for context := range kubeconfig.Contexts {
// 		testContext = context
// 		break // Choose the first context, you can modify this logic as needed
// 	}
// 	x := tempKubeconfig.Name()
// 	config, err := connect.ConnectToSource(testContext, &x)
// 	if err != nil {
// 		t.Fatalf("Error creating config: %v", err)
// 	}

// 	namespaces, err := GetNamespace(config, "namespace2")
// 	if err != nil {
// 		t.Errorf("Expected no error, got: %v", err)
// 	}

// 	if namespaces.Name != "namespace2" {
// 		t.Errorf("Expected namespace2, got: %v", namespaces)
// 	}
// }
