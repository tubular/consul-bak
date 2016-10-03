package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
)

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
