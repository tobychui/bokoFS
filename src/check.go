package main

import (
	"fmt"
	"os/exec"
)

// commandExists checks if a given command exists on the system
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// checkRuntimeEnvironment checks if the required commands are available in the runtime environment
func checkRuntimeEnvironment() bool {
	packageMissing := false
	commands := []string{"ffmpeg", "smartctl", "lsblk", "blkid", "df"}
	for _, cmd := range commands {
		if commandExists(cmd) {
			fmt.Printf("\033[32m✔\033[0m '%s' exists\n", cmd)
		} else {
			packageMissing = true
			fmt.Printf("\033[31m✘\033[0m '%s' does not exist\n", cmd)
		}
	}
	return !packageMissing
}
