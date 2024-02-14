package tools

import (
	"testing"
)

func TestIsJSONCompatible(t *testing.T) {
	tests := []struct {
		name     string
		jsonStr  string
		expected bool
	}{
		{
			name:     "Valid JSON string",
			jsonStr:  `{"key": "value"}`,
			expected: true,
		},
		{
			name:     "Invalid JSON string",
			jsonStr:  `{"key": "value"`,
			expected: false,
		},
		{
			name:     "Empty string",
			jsonStr:  "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isJSONCompatible(tt.jsonStr)
			if result != tt.expected {
				t.Errorf("Expected isJSONCompatible(%s) to be %t, got %t", tt.jsonStr, tt.expected, result)
			}
		})
	}
}

func TestPrettifyJSON(t *testing.T) {
	tests := []struct {
		name          string
		jsonStr       string
		expectedJSON  string
		expectedError bool
	}{
		{
			name:          "Valid JSON string",
			jsonStr:       `{"key": "value"}`,
			expectedJSON:  "{\n  \"key\": \"value\"\n}",
			expectedError: false,
		},
		{
			name:          "Invalid JSON string",
			jsonStr:       `{"key": "value"`,
			expectedJSON:  "",
			expectedError: true,
		},
		{
			name:          "Empty string",
			jsonStr:       "",
			expectedJSON:  "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := prettifyJSON(tt.jsonStr)
			if (err != nil) != tt.expectedError {
				t.Errorf("Expected error status does not match for prettifyJSON(%s)", tt.jsonStr)
			}
			if tt.expectedError {
				return
			}
			if result != tt.expectedJSON {
				t.Errorf("Expected prettifyJSON(%s) to be %s, got %s", tt.jsonStr, tt.expectedJSON, result)
			}
		})
	}
}
