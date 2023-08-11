package main

import (
	"log"
	"os/exec"
)

func CallRoutines() {}

func executeCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out), err
}

func checkDockerIsInstalled() bool {
	_, err := executeCommand("docker", "version")
	if err != nil {
		return false
	}
	return true
}
