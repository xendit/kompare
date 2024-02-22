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

func TestGetNamespace(t *testing.T) {
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

	namespace, err := GetNamespace(config, "namespace2")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if namespace.Name != "namespace2" {
		t.Errorf("Expected namespace2, got: %v", namespace)
	}
}

func TestListDeployments(t *testing.T) {
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

	// Perform the test
	deployments, err := ListDeployments(config, "namespace2")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	expectedLength := 3 // Update this value with the expected number of deployments
	if len(deployments.Items) != expectedLength {
		t.Errorf("Expected %d deployments, got: %d", expectedLength, len(deployments.Items))
	}
}

func TestGetHPA(t *testing.T) {
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

	// Perform the test
	hpas, err := ListHPAs(config, "namespace2")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	expectedLength := 3 // Update this value with the expected number of HPAs
	if len(hpas.Items) != expectedLength {
		t.Errorf("Expected %d HPAs, got: %d", expectedLength, len(hpas.Items))
	}
}

func TestListCronJobs(t *testing.T) {
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

	// Perform the test
	cronJobs, err := ListCronJobs(config, "namespace2")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	expectedLength := 3 // Update this value with the expected number of CronJobs
	if len(cronJobs.Items) != expectedLength {
		t.Errorf("Expected %d CronJobs, got: %d", expectedLength, len(cronJobs.Items))
	}
}

func TestListCRDs(t *testing.T) {
	// Set up test environment and get the temporary kubeconfig file
	_, _, tempKubeconfig := mock.SetupTestEnvironment()
	defer tempKubeconfig.Close() // Close the file after the test completes

	// Load the kubeconfig data
	tempKubeconfigByte, err := os.ReadFile(tempKubeconfig.Name())
	if err != nil {
		t.Fatalf("Error reading kubeconfig file: %v", err)
	}
	kubeconfig, err := clientcmd.Load(tempKubeconfigByte)
	if err != nil {
		t.Fatalf("Error loading kubeconfig: %v", err)
	}

	// Get the current context from the kubeconfig
	currentContext := kubeconfig.CurrentContext

	// Query the CRDs using the context and the kubeconfig file path
	crdList, err := ListCRDs(currentContext, tempKubeconfig.Name())
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	expectedLength := 3 // Adjust this based on your test case
	if len(crdList.Items) != expectedLength {
		t.Errorf("Expected %d CRDs, got: %d", expectedLength, len(crdList.Items))
	}
}

func TestListIngresses(t *testing.T) {
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

	// Connect to the Kubernetes cluster using the test context and kubeconfig file path
	x := tempKubeconfig.Name()
	config, err := connect.ConnectToSource(testContext, &x)
	if err != nil {
		t.Fatalf("Error creating config: %v", err)
	}

	// Perform the test
	ingresses, err := ListIngresses(config, "namespace2")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	expectedLength := 3 // Update this value with the expected number of ingresses
	if len(ingresses.Items) != expectedLength {
		t.Errorf("Expected %d ingresses, got: %d", expectedLength, len(ingresses.Items))
	}
}

func TestListServices(t *testing.T) {
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
	// Connect to the Kubernetes cluster using the test context and kubeconfig file path
	x := tempKubeconfig.Name()
	config, err := connect.ConnectToSource(testContext, &x)
	if err != nil {
		t.Fatalf("Error creating config: %v", err)
	}
	// Perform the test
	services, err := ListServices(config, "namespace2")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	expectedLength := 3 // Update this value with the expected number of services
	if len(services.Items) != expectedLength {
		t.Errorf("Expected %d services, got: %d", expectedLength, len(services.Items))
	}
}

// TestListConfigMaps tests the ListConfigMaps function.
func TestListConfigMaps(t *testing.T) {
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
	// Connect to the Kubernetes cluster using the test context and kubeconfig file path
	x := tempKubeconfig.Name()
	config, err := connect.ConnectToSource(testContext, &x)
	if err != nil {
		t.Fatalf("Error creating config: %v", err)
	}
	// Perform the test
	configMaps, err := ListConfigMaps(config, "namespace2")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	expectedLength := 3 // Update this value with the expected number of config maps
	if len(configMaps.Items) != expectedLength {
		t.Errorf("Expected %d config maps, got: %d", expectedLength, len(configMaps.Items))
	}
}

// TestListSecrets tests the ListSecrets function.
func TestListSecrets(t *testing.T) {
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

	// Connect to the Kubernetes cluster using the test context and kubeconfig file path
	x := tempKubeconfig.Name()
	config, err := connect.ConnectToSource(testContext, &x)
	if err != nil {
		t.Fatalf("Error creating config: %v", err)
	}

	// Perform the test
	secrets, err := ListSecrets(config, "namespace2")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	expectedLength := 3 // Update this value with the expected number of secrets
	if len(secrets.Items) != expectedLength {
		t.Errorf("Expected %d secrets, got: %d", expectedLength, len(secrets.Items))
	}
}

// TestListServiceAccounts tests the ListServiceAccounts function.
func TestListServiceAccounts(t *testing.T) {
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

	// Connect to the Kubernetes cluster using the test context and kubeconfig file path
	x := tempKubeconfig.Name()
	config, err := connect.ConnectToSource(testContext, &x)
	if err != nil {
		t.Fatalf("Error creating config: %v", err)
	}

	// Perform the test
	serviceAccounts, err := ListServiceAccounts(config, "namespace2")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	expectedLength := 3 // Update this value with the expected number of service accounts
	if len(serviceAccounts.Items) != expectedLength {
		t.Errorf("Expected %d service accounts, got: %d", expectedLength, len(serviceAccounts.Items))
	}
}

// TestListRoles tests the ListRoles function.
func TestListRoles(t *testing.T) {
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

	// Connect to the Kubernetes cluster using the test context and kubeconfig file path
	x := tempKubeconfig.Name()
	config, err := connect.ConnectToSource(testContext, &x)
	if err != nil {
		t.Fatalf("Error creating config: %v", err)
	}

	// Perform the test
	roles, err := ListRoles(config, "namespace2")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	expectedLength := 3 // Update this value with the expected number of roles
	if len(roles.Items) != expectedLength {
		t.Errorf("Expected %d roles, got: %d", expectedLength, len(roles.Items))
	}
}

// TestListRoleBindings tests the ListRoleBindings function.
func TestListRoleBindings(t *testing.T) {
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

	// Connect to the Kubernetes cluster using the test context and kubeconfig file path
	x := tempKubeconfig.Name()
	config, err := connect.ConnectToSource(testContext, &x)
	if err != nil {
		t.Fatalf("Error creating config: %v", err)
	}

	// Perform the test
	roleBindings, err := ListRoleBindings(config, "namespace2")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	expectedLength := 3 // Update this value with the expected number of role bindings
	if len(roleBindings.Items) != expectedLength {
		t.Errorf("Expected %d role bindings, got: %d", expectedLength, len(roleBindings.Items))
	}
}

// TestListClusterRoles tests the ListClusterRoles function.
func TestListClusterRoles(t *testing.T) {
	// Set up test environment and get the temporary kubeconfig file
	_, _, tempKubeconfig := mock.SetupTestEnvironment()
	defer tempKubeconfig.Close() // Close the file after the test completes

	// Load the kubeconfig data
	tempKubeconfigByte, err := os.ReadFile(tempKubeconfig.Name())
	if err != nil {
		t.Fatalf("Error reading kubeconfig file: %v", err)
	}
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
	// Perform the test
	clusterRoles, err := ListClusterRoles(config)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	expectedLength := 3 // Update this value with the expected number of cluster roles
	if len(clusterRoles.Items) != expectedLength {
		t.Errorf("Expected %d cluster roles, got: %d", expectedLength, len(clusterRoles.Items))
	}
}

// TestListRoleBindings tests the ListRoleBindings function.
func TestListRoleBindings2(t *testing.T) {
	// Set up test environment and get the temporary kubeconfig file
	_, _, tempKubeconfig := mock.SetupTestEnvironment()
	defer tempKubeconfig.Close() // Close the file after the test completes

	// Load the kubeconfig data
	tempKubeconfigByte, err := os.ReadFile(tempKubeconfig.Name())
	if err != nil {
		t.Fatalf("Error reading kubeconfig file: %v", err)
	}
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

	// Perform the test
	roleBindings, err := ListRoleBindings(config, "namespace2")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	expectedLength := 3 // Update this value with the expected number of role bindings
	if len(roleBindings.Items) != expectedLength {
		t.Errorf("Expected %d role bindings, got: %d", expectedLength, len(roleBindings.Items))
	}

}

// TestNetworkPolicies tests the ListNetworkPolicies function.
func TestListNetworkPolicies(t *testing.T) {
	// Set up test environment and get the temporary kubeconfig file
	_, _, tempKubeconfig := mock.SetupTestEnvironment()
	defer tempKubeconfig.Close() // Close the file after the test completes

	// Load the kubeconfig data
	tempKubeconfigByte, err := os.ReadFile(tempKubeconfig.Name())
	if err != nil {
		t.Fatalf("Error reading kubeconfig file: %v", err)
	}
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

	// Connect to the Kubernetes cluster using the test context and kubeconfig file path
	x := tempKubeconfig.Name()
	config, err := connect.ConnectToSource(testContext, &x)
	if err != nil {
		t.Fatalf("Error creating config: %v", err)
	}

	// Perform the test
	networkPolicies, err := ListNetworkPolicies(config, "namespace2")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	expectedLength := 3 // Update this value with the expected number of network policies
	if len(networkPolicies.Items) != expectedLength {
		t.Errorf("Expected %d network policies, got: %d", expectedLength, len(networkPolicies.Items))
	}
}
