package tools

import (
	"fmt"
	"io"
	"log"
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
	// Defer the closure of the file and its removal.
	defer func() {
		file.Close()
		os.Remove(absPath)
	}()

	return true, absPath, nil
}

func LogOutput(logfile string) func() {
	// open file read/write | create if not exist | clear file at open if exists
	f, _ := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

	// save existing stdout | MultiWriter writes to saved stdout and file
	out := os.Stdout
	mw := io.MultiWriter(out, f)

	// get pipe reader and writer | writes to pipe writer come out pipe reader
	r, w, _ := os.Pipe()

	// replace stdout,stderr with pipe writer | all writes to stdout, stderr will go through pipe instead (fmt.print, log)
	os.Stdout = w
	os.Stderr = w

	// writes with log.Print should also write to mw
	log.SetOutput(mw)

	//create channel to control exit | will block until all copies are finished
	exit := make(chan bool)

	go func() {
		// copy all reads from pipe to multiwriter, which writes to stdout and file
		_, _ = io.Copy(mw, r)
		// when r or w is closed copy will finish and true will be sent to channel
		exit <- true
	}()

	// function to be deferred in main until program exits
	return func() {
		// close writer then block on exit channel | this will let mw finish writing before the program exits
		_ = w.Close()
		<-exit
		// close file after all writes have finished
		_ = f.Close()
	}

}
