package containers

import (
	"encoding/json"
	"log"
	"os/exec"
	"strings"
	"sync"
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

func GetAllContainers(host string, ssh any, wg *sync.WaitGroup) <-chan []DockerContainer {
	defer wg.Done()
	cc := make(chan []DockerContainer, 100)

	if ssh, ok := ssh.(*SSH); ok {
		cmdToRun := `docker ps -a --format '{"ContainerId":"{{.ID}}", "Image":"{{.Image}}"}' | jq -s .`
		ssh.Channels[host].Command <- cmdToRun

		var containers []DockerContainer
		decoder := json.NewDecoder(strings.NewReader(string(<-ssh.Channels[host].Response)))
		if err := decoder.Decode(&containers); err != nil {
			panic("Error decoding JSON:" + err.Error())
		}
		cc <- containers

	} else {
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
	return cc
}

func GetLogs(cl chan<- []byte, containerId string, wg *sync.WaitGroup) {
	defer wg.Done()

	cmd := exec.Command("docker", "logs", containerId)

	logs, err := cmd.Output()
	if err != nil {
		log.Fatal("Error getting logs:", err)
	}

	cl <- logs

}

func Inspect(cl chan<- []byte, containerId string, wg *sync.WaitGroup) {
	defer wg.Done()

	cmd := exec.Command("docker", "inspect", containerId)

	inspect, err := cmd.Output()
	if err != nil {
		log.Fatal("Error inspecting container:", err)
	}

	cl <- inspect

}
