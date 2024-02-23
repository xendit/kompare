package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// isValidPath checks if the given path is valid for the current operating system.
func IsValidPath(path string) (isValid bool, absPath string, err error) {
	// Check if the path contains invalid characters for the current OS.
	var invalidChars string
	switch runtime.GOOS {
	case "windows":
		invalidChars = `<>|?*`
	default:
		invalidChars = `<>:"\|?*`
	}

	for _, char := range invalidChars {
		if strings.ContainsRune(path, char) {
			return false, "", fmt.Errorf("invalid character '%c' found in the path", char)
		}
	}

	// Check if the path is empty or consists of only whitespace characters.
	if len(strings.TrimSpace(path)) == 0 {
		return false, "", fmt.Errorf("empty path")
	}

	// Convert the path to absolute path if it's relative.
	if !filepath.IsAbs(path) {
		// Get the current working directory.
		cwd, err := os.Getwd()
		if err != nil {
			// Failed to get the current working directory.
			return false, "", err
		}
		// Join the current working directory with the relative path to get absolute path.
		absPath = filepath.Join(cwd, path)
	} else {
		absPath = path
	}

	// Try to create or open the file.
	file, err := os.OpenFile(absPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		// Check if the error is due to a missing directory in the path.
		if os.IsNotExist(err) {
			return false, "", fmt.Errorf("the file or one or more directories in the path do not exist or can not be created: %v", err)
		}
		// Check if the error is due to permission issues.
		if os.IsPermission(err) {
			return false, "", fmt.Errorf("no permission to write on path: %s", absPath)
		}
		// For other errors, return the error as is.
		return false, "", err
	}
	defer file.Close()

	return true, absPath, nil
}
