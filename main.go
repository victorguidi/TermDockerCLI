package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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

// Populate function should accept a generic channel and a generic struct
func populate[T IGeneric](t *tview.TextView, data <-chan []T, wg *sync.WaitGroup) {
	defer wg.Done()
	for d := range data {
		s := fmt.Sprintf("%+v", d)
		t.SetText(s)
	}
}

func main() {

	CallRoutines()

	// TODO: Change the layout, is to ugly
	// TODO: Add a terminal UI to interact with the data
	// TODO: Add a way to refresh the data
	// TODO: Format the data inside the boxes

	// Create the list view to display the Docker images
	containers := tview.NewTextView()
	containers.SetBorder(true).SetTitle("Docker Containers")

	images := tview.NewTextView()
	images.SetBorder(true).SetTitle("Docker Images")

	networks := tview.NewTextView()
	networks.SetBorder(true).SetTitle("Docker Networks")

	wg := sync.WaitGroup{}
	wg.Add(3)
	go populate(containers, dockerContainers, &wg)
	go populate(images, dockerImages, &wg)
	go populate(networks, dockerNetworks, &wg)
	wg.Wait()

	// Create the left panel that holds the list view
	leftPanel := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(containers, 0, 1, true).
		AddItem(images, 0, 1, true).
		AddItem(networks, 0, 1, true)

	// Create the right panel that holds the box (your current code)
	rightPanel := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Docker CLI"), 0, 1, false)

	// Create the layout
	layout := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(leftPanel, 0, 2, false).
		AddItem(rightPanel, 0, 4, false)

	layout.SetBackgroundColor(tcell.ColorBlack)

	// Create the application
	app := tview.NewApplication()

	// Run the application
	if err := app.SetRoot(layout, true).Run(); err != nil {
		panic(err)
	}

	// TODO: Check how to load the data on the terminal UI

	// println("Docker Images")
	// dockerImages := <-dockerImages
	// fmt.Println(dockerImages)
	//
	// println("Docker Containers")
	// dockerContainers := <-dockerContainers
	// fmt.Println(dockerContainers)
	//
	// println("Docker Networks")
	// dockerNetworks := <-dockerNetworks
	// fmt.Println(dockerNetworks)

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
