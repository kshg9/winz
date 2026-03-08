package main

import (
	"os/exec"
	"strings"
)

// runCmd runs a command and returns combined stdout as a string.
func runCmd(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).CombinedOutput()
	return strings.TrimSpace(string(out)), err
}
