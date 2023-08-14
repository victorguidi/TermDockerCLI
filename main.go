package main

// TODO: Change the layout, is to ugly
// TODO: Add the possibility to send commands
// TODO: Get the logs of the containers in ssh too

import (
	"github.com/victorguidi/TermDockerCLI/application"
	"log"
)

func init() {
	if !checkDockerIsInstalled() {
		log.Fatal("Docker is not installed, please install docker")
		return
	}
	// Open SSH Connections in case they exist and close them when the application is closed

}

func main() {

	application := application.NewApplication()
	application.Build()
	application.Run()

}
