package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
)

// EnsureDir creates a directory at `path` if it doesn't already exist
func EnsureDir(path string, perm os.FileMode) error {
	logger.Infof("Making sure a directory exists at: %s", path)
	fi, err := os.Stat(path)
	if err == nil {
		if fi.IsDir() {
			return nil
		} else {
			err := errors.New(fmt.Sprintf("Cannot create directory '%s': File exists.", path))
			logger.Fatal(err)
			return err
		}
	} else if os.IsNotExist(err) {
		// nothing found, create a directory
		mkdirErr := os.MkdirAll(path, perm)
		if mkdirErr != nil {
			logger.Errorf("Failed to create a directory at: %s (%s)", path, mkdirErr)
			return mkdirErr
		}
		return nil
	}
	logger.Fatal(err)
	return err
}

// PathExists verifies that the path at path exists
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// Check is a utility function for checking error returns
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

// StartsWith checks to verify if a string starts with a given substring
func StartsWith(list []string, elem string) bool {
	for _, t := range list {
		if strings.HasPrefix(elem, t) {
			return true
		}
	}
	return false
}

// CheckSocket checks to see if a socket `endpoint` is listening.
func CheckSocket(endpoint string) bool {
	logger.Debugf("Checking that socket %s is listening", endpoint)
	_, err := net.Dial("tcp", endpoint)
	if err != nil {
		logger.Fatalf("No socket listening at %s. Giving up: %s", endpoint, err)
		os.Exit(1)
	}
	return true
}

// Which checks to see if a given utility exists on the local machine
func Which(utility string) bool {
	_, err := exec.LookPath(utility)
	if err != nil {
		return false
	}
	return true
}

// ConsulBinaryCall runs an arbitrary command using the consul CLI utility
func ConsulBinaryCall(a, b string) string {
	out, err := exec.Command("consul", a, b).Output()
	if err != nil {
		message := fmt.Sprintf("There was an error querying consul. Giving up: %s", err)
		fmt.Println(message)
		os.Exit(1)
	}
	return string(out)
}

// GitBinaryCall runs an arbitrary command using the git CLI utility
func GitBinaryCall(cmd string) string {
	fullCommand := fmt.Sprintf("git %s", cmd)
	args := strings.Fields(cmd)
	logger.Debugf("Calling git with: `%s`", fullCommand)
	out, err := exec.Command("git", args...).CombinedOutput()
	if err != nil {
		message := fmt.Sprintf("There was an error running git. Giving up: %s", err)
		fmt.Println(message)
		os.Exit(1)
	}
	return string(out)
}

// DeleteEmpty removes the empty strings from a string array
func DeleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
