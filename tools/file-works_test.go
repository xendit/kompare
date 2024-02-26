package tools

import (
	"log"
	"os"
	"strings"
	"testing"
)

func TestIsValidPath(t *testing.T) {
	// Example usage:
	paths := map[string]struct {
		expectedValid bool
		expectedError string
	}{
		"/path/to/valid/file.txt": {
			expectedValid: false,
			expectedError: "the file or one or more directories in the path do not exist or can not be created",
		},
		"invalid*file.txt": {
			expectedValid: false,
			expectedError: "invalid character '*' found in the path",
		},
		"": {
			expectedValid: false,
			expectedError: "empty path",
		},
		"  ": {
			expectedValid: false,
			expectedError: "empty path",
		},
		"/path/with/invalid|character": {
			expectedValid: false,
			expectedError: "invalid character '|' found in the path",
		},
		"C:\\Program Files\\Example\\file.txt": {
			expectedValid: false,
			expectedError: "invalid character ':' found in the path",
		},
		"D:/Documents/Report.docx": {
			expectedValid: false,
			expectedError: "invalid character ':' found in the path",
		},
		"/home/user/pictures/photo.jpg": {
			expectedValid: false,
			expectedError: "the file or one or more directories in the path do not exist or can not be created",
		},
		"file.txt": {
			expectedValid: true,
			expectedError: "",
		},
		"folder1/file.txt": {
			expectedValid: false,
			expectedError: "the file or one or more directories in the path do not exist or can not be created",
		},
		"../parent/file.txt": {
			expectedValid: false,
			expectedError: "the file or one or more directories in the path do not exist or can not be created",
		},
		"..\\parent\\file.txt": {
			expectedValid: false,
			expectedError: "invalid character '\\' found in the path",
		},
	}

	for path, expected := range paths {
		valid, _, err := IsValidPath(path)
		if err != nil {
			if !strings.Contains(err.Error(), expected.expectedError) {
				t.Errorf("Unexpected error for path %s. Expected: %s, Got: %s", path, expected.expectedError, err.Error())
			}
		} else {
			if valid != expected.expectedValid {
				t.Errorf("Path %s validation failed. Expected valid: %t but got: %t", path, expected.expectedValid, valid)
			}
		}
	}
}

func TestLogOutput(t *testing.T) {
	// Create a temporary file for testing.
	tmpfile, err := os.CreateTemp("", "testlogoutput*.log")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpfile.Name()) // Clean up the temporary file.

	// Call logOutput to redirect output to the temporary file.
	deferredFunc := LogOutput(tmpfile.Name())

	// Send some output to the log.
	log.Print("Test log message")

	// Close the log output.
	deferredFunc()

	// Read the contents of the temporary file.
	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("failed to read temporary file: %v", err)
	}

	// Check if the log message is present in the file content.
	if !strings.Contains(string(content), "Test log message") {
		t.Errorf("log message not found in file content")
	}
}
