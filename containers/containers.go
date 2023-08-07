package containers

import (
	"encoding/json"
	"log"
	"os/exec"
)

type DockerContainer struct {
	ContainerId string
	Image       string
	Command     string
	Created     string
	Status      string
	Ports       string
	Names       string
}

func NewContainer() *DockerContainer {
	return &DockerContainer{}
}

func (c *DockerContainer) GetAllContainers(cc chan []DockerContainer) {
	defer close(cc)

	cmd := exec.Command("docker", "ps", "-a", "--format", `{"ContainerId":"{{.ID}}", "Image":"{{.Image}}"}`)
	jqCmd := exec.Command("jq", "-s", ".")

	jqCmd.Stdin, _ = cmd.StdoutPipe()

	jqOutput, err := jqCmd.StdoutPipe()
	if err != nil {
		log.Fatal("Error creating pipe:", err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal("Error starting docker command:", err)
	}
	if err := jqCmd.Start(); err != nil {
		log.Fatal("Error starting jq command:", err)
	}

	var containers []DockerContainer
	decoder := json.NewDecoder(jqOutput)
	if err := decoder.Decode(&containers); err != nil {
		log.Fatal("Error decoding JSON:", err)
	}

	if err := jqCmd.Wait(); err != nil {
		log.Fatal("Error waiting for jq command:", err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal("Error waiting for docker command:", err)
	}

	cc <- containers
}

func GetLogs(cl chan []byte, containerId string) {

	cmd := exec.Command("docker", "logs", containerId)

	logs, err := cmd.Output()
	if err != nil {
		log.Fatal("Error getting logs:", err)
	}

	cl <- logs
}
