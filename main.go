package main

import (
	"fmt"
	"github.com/rivo/tview"
	"log"
	"os/exec"
	"strings"
)

// TODO: Check if docker is installed, if not, throw an error
// TODO: Add a terminal UI to interact with the data

var (
	dockerImages     chan []DockerImage
	dockerContainers chan []DockerContainer
	dockerNetworks   chan []DockerNetwork
)

func init() {
	if !checkDockerIsInstalled() {
		log.Fatal("Docker is not installed, please install docker")
		return
	}

	dockerImages = make(chan []DockerImage, 1)
	dockerContainers = make(chan []DockerContainer, 1)
	dockerNetworks = make(chan []DockerNetwork, 1)
}

func executeCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out), err
}

func main() {

	CallRoutines()

	box := tview.NewBox().SetBorder(true).SetTitle("Docker CLI")
	if err := tview.NewApplication().SetRoot(box, true).Run(); err != nil {

		// TODO: Check how to load the data on the terminal UI

		println("Docker Images")
		dockerImages := <-dockerImages
		fmt.Println(dockerImages)

		println("Docker Containers")
		dockerContainers := <-dockerContainers
		fmt.Println(dockerContainers)

		println("Docker Networks")
		dockerNetworks := <-dockerNetworks
		fmt.Println(dockerNetworks)
		panic(err)
	}

}

func splitLines(s string, t Setter) []interface{} {
	var result []interface{}
	lines := strings.Split(s, "\n")
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 1 {
			continue
		}
		t.Set(fields)
		result = append(result, t)
	}
	return result
}

func checkDockerIsInstalled() bool {
	_, err := executeCommand("docker", "version")
	if err != nil {
		return false
	}
	return true
}
