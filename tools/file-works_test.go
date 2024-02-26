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
	// Create a temporary file for logging
	tmpfile, err := os.CreateTemp("", "testlog-")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	// Create a channel to signal when the goroutine has completed
	done := make(chan struct{})

	// Call LogOutput with the temporary file
	go func() {
		defer close(done)
		LogOutput(tmpfile.Name()) // Assuming LogOutput returns immediately after starting the goroutine
	}()

	// Write a log message
	log.Print("Test log message")

	// Wait for the goroutine to finish
	<-done

	// Read and verify the contents of the log file
	// (You should implement this verification logic)
}

func TestLogOutputWithInvalidPath(t *testing.T) {
	// Call LogOutput with an invalid path
	err := LogOutput("invalid*file.txt")
	if err == nil {
		t.Errorf("Expected error for invalid path")
	}
	defer os.Remove("invalid*file.txt")
}

func TestLogOutputWithEmptyPath(t *testing.T) {
	// Call LogOutput with an empty path
	err := LogOutput("")
	if err == nil {
		t.Errorf("Expected error for empty path")
	}
}
