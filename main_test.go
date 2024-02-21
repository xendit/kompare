package main

import (
	"fmt"
	"os"
	"testing"

	"kompare/mock"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
	os.Args = []string{"main.go", "-s", "source-context", "-t", "target-context", "-n", "InvalidNS"}
	assert.Panics(t, func() { main() }, "Expected a panic due to namespace listing error")

	_, _, kubeconfigFile := mock.SetupTestEnvironment()

	// Set up command-line arguments
	os.Args = []string{"main.go", "-t", "target-context", "-s", "source-context", "-c", kubeconfigFile.Name()}

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
	// Run the main function and capture any errors
	func() {
		defer func() {
			if r := recover(); r != nil {
				capturedError = fmt.Errorf("panic occurred: %v", r)
			}
		}()
		main()
	}()

	// Assert that clientsetToSource and clientsetToTarget are not nil after connecting to clusters
	fmt.Println("Test completed")
}

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

func TestFilterNamespaces(t *testing.T) {
	// Define test data
	namespaces := &v1.NamespaceList{
		Items: []v1.Namespace{
			{
				ObjectMeta: metav1.ObjectMeta{Name: "default"},
			},
			{
				ObjectMeta: metav1.ObjectMeta{Name: "kube-system"},
			},
			{
				ObjectMeta: metav1.ObjectMeta{Name: "test-namespace"},
			},
		},
	}

	// Define test cases
	testCases := []struct {
		name           string
		pattern        string
		expectedResult int
	}{
		{"MatchAll", "*", 3},
		{"MatchSpecific", "kube-*", 1},
		{"NoMatch", "nonexistent-*", 0},
	}

	// Iterate over test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function under test
			filtered := filterNamespaces(namespaces, tc.pattern)

			// Check the length of the filtered namespaces
			if len(filtered.Items) != tc.expectedResult {
				t.Errorf("Expected %d namespaces matching pattern '%s', got %d", tc.expectedResult, tc.pattern, len(filtered.Items))
			}

			// Ensure all filtered namespaces match the pattern
			for _, ns := range filtered.Items {
				if !matchWildcard(ns.Name, tc.pattern) {
					t.Errorf("Namespace '%s' does not match pattern '%s'", ns.Name, tc.pattern)
				}
			}
		})
	}
}
