// connect/connect_test.go

package connect

import (
	"os"
	"testing"

	"kompare/mock" // Import the mock package for testing

	"k8s.io/client-go/tools/clientcmd"
)

func TestCreateConfig(t *testing.T) {
	// Set up test environment and get the temporary kubeconfig file
	_, _, tempKubeconfig := mock.SetupTestEnvironment()
	defer tempKubeconfig.Close() // Close the file after the test completes

	// Call CreateConfig with the path of the temporary kubeconfig file
	x := tempKubeconfig.Name()
	config, err := CreateConfig(&x)
	if err != nil {
		t.Fatalf("Error creating config: %v", err)
	}

	// Validate the created config
	if config == nil {
		t.Error("Expected non-nil config, got nil")
	}
	// Add more validation as needed, such as checking server URLs, etc.
}

func TestNewK8sConnectionConfig(t *testing.T) {
	// Set up test environment and get the temporary kubeconfig file
	_, _, tempKubeconfig := mock.SetupTestEnvironment()
	defer tempKubeconfig.Close() // Close the file after the test completes

	// Create a mock configuration for testing
	x := tempKubeconfig.Name()
	config, err := CreateConfig(&x)
	if err != nil {
		t.Fatalf("Error creating config: %v", err)
	}

	// Call NewK8sConnectionConfig with the mock configuration
	clientset, err := NewK8sConnectionConfig(config)
	if err != nil {
		t.Fatalf("Error creating clientset: %v", err)
	}

	// Validate the created clientset
	if clientset == nil {
		t.Error("Expected non-nil clientset, got nil")
	}
}

func TestConnectNow(t *testing.T) {
	// Set up test environment and get the temporary kubeconfig file
	_, _, tempKubeconfig := mock.SetupTestEnvironment()
	defer tempKubeconfig.Close() // Close the file after the test completes

	// Call CreateConfig with the path of the temporary kubeconfig file
	x := tempKubeconfig.Name()
	config, err := ConnectNow(&x)
	if err != nil {
		t.Fatalf("Error creating config: %v", err)
	}

	// Validate the created config
	if config == nil {
		t.Error("Expected non-nil config, got nil")
	}
	// Add more validation as needed, such as checking server URLs, etc.
}
func TestBuildConfigWithContextFromFlags(t *testing.T) {
	// Set up test environment and get the temporary kubeconfig file
	_, _, tempKubeconfig := mock.SetupTestEnvironment()
	defer tempKubeconfig.Close() // Close the file after the test completes

	// Read the content of the kubeconfig file
	x, err := os.ReadFile(tempKubeconfig.Name())
	if err != nil {
		t.Fatalf("Error reading kubeconfig file: %v", err)
	}

	// Load the kubeconfig data
	kubeconfig, err := clientcmd.Load(x)
	if err != nil {
		t.Fatalf("Error loading kubeconfig: %v", err)
	}

	// Choose one of the contexts from the kubeconfig
	var testContext string
	for context := range kubeconfig.Contexts {
		testContext = context
		break // Choose the first context, you can modify this logic as needed
	}

	// Call BuildConfigWithContextFromFlags with the test context and temp kubeconfig path
	config, err := BuildConfigWithContextFromFlags(testContext, tempKubeconfig.Name())
	if err != nil {
		t.Fatalf("Error building config with context: %v", err)
	}

	// Validate the created config
	if config == nil {
		t.Error("Expected non-nil config, got nil")
	}

	// Add more validation if needed
}

func TestConnectToSource(t *testing.T) {
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
	config, err := ConnectToSource(testContext, &x)
	if err != nil {
		t.Fatalf("Error creating config: %v", err)
	}

	// Validate the created config
	if config == nil {
		t.Error("Expected non-nil config, got nil")
	}
	// Add more validation as needed, such as checking server URLs, etc.
}

func TestContextSwitch(t *testing.T) {
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
	config, err := ContextSwitch(testContext, &x)
	if err != nil {
		t.Fatalf("Error creating config: %v", err)
	}

	// Validate the created config
	if config == nil {
		t.Error("Expected non-nil config, got nil")
	}
	// Add more validation as needed, such as checking server URLs, etc.
}
