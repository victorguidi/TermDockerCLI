package main

// TODO: Change the layout, is to ugly
// TODO: Add the possibility to send commands
// TODO: Get the logs of the containers in ssh too

import (
	"log"

	"github.com/victorguidi/TermDockerCLI/application"
)

func init() {
	if !checkDockerIsInstalled() {
		log.Fatal("Docker is not installed, please install docker")
		return
	}
}

func main() {

	application := application.NewApplication()
	application.Build()
	application.Run()

}
