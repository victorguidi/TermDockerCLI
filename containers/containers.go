package containers

import ()

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
