package cli

import (
	"os"
	"reflect"
	"testing"
)

func TestPaserReader(t *testing.T) {
	// Store original os.Args and defer its restoration
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Define test command-line arguments
	testArgs := []string{
		"program_name",
		"-t", "target-context",
		"-i", "deployment,service",
		"-v",
	}

	// Set os.Args to testArgs
	os.Args = testArgs

	// Call the function under test
	args := PaserReader()

	// Add assertions to validate the parsed arguments
	// For example, you can check if the target context and included objects are set correctly
	if args.TargetClusterContext != "target-context" {
		t.Errorf("Expected target context 'target-context', got %s", args.TargetClusterContext)
	}

	if args.Include[0] != "deployment" || args.Include[1] != "service" {
		t.Errorf("Expected included objects 'deployment' and 'service', got %v", args.Include)
	}

	if args.VerboseDiffs != 1 {
		t.Errorf("Expected verbose mode enabled, got %d", args.VerboseDiffs)
	}
}

func TestValidateParametersFromParserArgs(t *testing.T) {
	testCases := []struct {
		name                  string
		argsReceived          ArgumentsReceived
		expectedArgsValidated ArgumentsReceivedValidated
	}{
		{
			name: "Valid Arguments",
			argsReceived: ArgumentsReceived{
				KubeconfigFile:       stringPtr("/path/to/kubeconfig"),
				SourceClusterContext: stringPtr("source-context"),
				TargetClusterContext: stringPtr("target-context"),
				NamespaceName:        stringPtr("default"),
				FiltersForObject:     stringPtr("filter"),
				Include:              stringPtr("deployment"),
				Exclude:              stringPtr("service"),
				VerboseDiffs:         intPtr(1),
				Err:                  nil,
			},
			expectedArgsValidated: ArgumentsReceivedValidated{
				KubeconfigFile:       "/path/to/kubeconfig",
				SourceClusterContext: "source-context",
				TargetClusterContext: "target-context",
				NamespaceName:        "default",
				FiltersForObject:     "filter",
				Include:              []string{"deployment"},
				Exclude:              []string{"service"},
				VerboseDiffs:         1,
				Err:                  nil,
			},
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function under test
			actualArgsValidated := ValidateParametersFromParserArgs(tc.argsReceived)

			// Compare the actual and expected results
			if !argsValidatedEqual(actualArgsValidated, tc.expectedArgsValidated) {
				t.Errorf("Test case '%s': Expected %v, got %v", tc.name, tc.expectedArgsValidated, actualArgsValidated)
			}
		})
	}
}

// Utility functions for creating pointers

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

// Helper function to check equality of ArgumentsReceivedValidated structs
func argsValidatedEqual(a, b ArgumentsReceivedValidated) bool {
	return a.KubeconfigFile == b.KubeconfigFile &&
		a.SourceClusterContext == b.SourceClusterContext &&
		a.TargetClusterContext == b.TargetClusterContext &&
		a.NamespaceName == b.NamespaceName &&
		a.FiltersForObject == b.FiltersForObject &&
		stringSliceEqual(a.Include, b.Include) &&
		stringSliceEqual(a.Exclude, b.Exclude) &&
		a.VerboseDiffs == b.VerboseDiffs &&
		a.Err == b.Err
}

// Helper function to check equality of string slices
func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestValidateKubernetesObjects(t *testing.T) {
	testCases := []struct {
		name            string
		inputObjects    []string
		expectedInvalid []string
		expectedValid   []string
	}{
		{
			name:            "Valid Objects",
			inputObjects:    []string{"deployment", "ing", "sa", "configmaps", "secret", "namespace", "role", "rolebindings", "clusterrolebinding", "crd"},
			expectedInvalid: nil,
			expectedValid:   []string{"deployment", "ingress", "serviceaccount", "configmap", "secret", "namespace", "role", "rolebinding", "clusterrolebinding", "crd"},
		},
		{
			name:            "Mixed Objects",
			inputObjects:    []string{"svc", "invalid", "configmap", "rolebinding", "random"},
			expectedInvalid: []string{"invalid", "random"},
			expectedValid:   []string{"service", "configmap", "rolebinding"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			invalidObjects, validObjects := ValidateKubernetesObjects(tc.inputObjects)

			// Check if invalid objects match
			if !reflect.DeepEqual(invalidObjects, tc.expectedInvalid) {
				t.Errorf("Test case '%s': Expected invalid objects %v, got %v", tc.name, tc.expectedInvalid, invalidObjects)
			}

			// Check if valid objects match
			if !reflect.DeepEqual(validObjects, tc.expectedValid) {
				t.Errorf("Test case '%s': Expected valid objects %v, got %v", tc.name, tc.expectedValid, validObjects)
			}
		})
	}
}
