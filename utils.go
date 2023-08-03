package main

import "log"

func CallRoutines() {
	go func() {
		defer close(dockerContainers)
		var dcontainers []DockerContainer
		out, err := executeCommand("docker", "ps", "-a")
		if err != nil {
			log.Fatal(err)
		}
		containers := splitLines(out, &DockerContainer{})
		for _, container := range containers {
			dcontainers = append(dcontainers, *container.(*DockerContainer))
		}
		dockerContainers <- dcontainers
	}()

	go func() {
		defer close(dockerImages)
		var dimages []DockerImage
		out, err := executeCommand("docker", "images")
		if err != nil {
			log.Fatal(err)
		}
		images := splitLines(out, &DockerImage{})
		for _, image := range images {
			dimages = append(dimages, *image.(*DockerImage))
		}
		dockerImages <- dimages
	}()

	go func() {
		defer close(dockerNetworks)
		var dnetworks []DockerNetwork
		out, err := executeCommand("docker", "network", "ls")
		if err != nil {
			log.Fatal(err)
		}
		networks := splitLines(out, &DockerNetwork{})
		for _, network := range networks {
			dnetworks = append(dnetworks, *network.(*DockerNetwork))
		}
		dockerNetworks <- dnetworks
	}()
}
