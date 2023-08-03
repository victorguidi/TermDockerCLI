package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// TODO: Abstract all of that to run in go routines so it loads all of the data at once
// TODO: Check if docker is installed, if not, throw an error

type DockerContainer struct {
	ContainerId string
	Image       string
	Command     string
	Created     string
	Status      string
	Ports       string
	Names       string
}

type DockerImage struct {
	Repository string
	Tag        string
	ImageId    string
	Created    string
	Size       string
}

func init() {
	if !checkDockerIsInstalled() {
		log.Fatal("Docker is not installed, please install docker")
		return
	}
}

func main() {
	ps := exec.Command("docker", "ps", "-a")
	cout, err := ps.Output()
	if err != nil {
		println(err.Error())
		return
	}
	// println(string(out))

	images := exec.Command("docker", "images")
	iout, err := images.Output()
	if err != nil {
		println(err.Error())
		return
	}

	dockerContainer := []DockerContainer{}
	dockerImages := []DockerImage{}

	lines := strings.Split(string(cout), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if fields[0] == "CONTAINER" {
			continue
		}
		container := DockerContainer{
			ContainerId: fields[0],
			Image:       fields[1],
			Command:     fields[2],
			Created:     fields[3],
			Status:      fields[4],
			Ports:       fields[5],
			Names:       fields[6],
		}
		dockerContainer = append(dockerContainer, container)
	}

	lines = strings.Split(string(iout), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if fields[0] == "REPOSITORY" {
			continue
		}
		image := DockerImage{
			Repository: fields[0],
			Tag:        fields[1],
			ImageId:    fields[2],
			Created:    fields[3],
			Size:       fields[4],
		}
		dockerImages = append(dockerImages, image)
	}

	println("Docker Containers")
	fmt.Println(dockerContainer)

	println("Docker Images")
	fmt.Println(dockerImages)

}

func checkDockerIsInstalled() bool {
	return true
}
