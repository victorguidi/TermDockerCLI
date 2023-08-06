package main

// TODO: Add a terminal UI to interact with the data
// TODO: Change the layout, is to ugly
// TODO: Add a way to refresh the data
// TODO: Format the data inside the boxes
// TODO: Load the logs on separete channels

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"

	"github.com/victorguidi/TermDockerCLI/containers"
	"github.com/victorguidi/TermDockerCLI/images"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

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
		t.SetText(fmt.Sprintf("%v", d))
	}
}

func getContainerLogs(containerID string) string {
	out, err := executeCommand("docker", "logs", containerID)
	if err != nil {
		log.Fatal(err)
	}
	return out
}

func main() {

	CallRoutines()

	container := containers.NewContainerUi()
	container.PopulateUi()

	image := images.NewImageUi()
	image.PopulateUi()

	// containers.SetBorder(true).SetTitle("Docker Containers")
	// containers := tview.NewTextView()
	//
	// images := tview.NewTextView()
	// images.SetBorder(true).SetTitle("Docker Images")
	//
	// networks := tview.NewTextView()
	// networks.SetBorder(true).SetTitle("Docker Networks")

	// wg := sync.WaitGroup{}
	// wg.Add(3)
	// go populate(containers, dockerContainers, &wg)
	// go populate(images, dockerImages, &wg)
	// go populate(networks, dockerNetworks, &wg)
	// wg.Wait()

	// Create the left panel that holds the list view
	leftPanel := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(container.Table, 0, 1, true).
		AddItem(image.Table, 0, 1, true)

	text := tview.NewTextView()
	text.SetBorder(true).SetTitle("Docker TUI")
	rightPanel := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(text, 0, 1, true)

	// Create the layout
	layout := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(leftPanel, 0, 2, true).
		AddItem(rightPanel, 0, 4, true)

	layout.SetBackgroundColor(tcell.ColorBlack)

	// Create the application
	app := tview.NewApplication()

	// with q key we can quit the application
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case 'q':
				app.Stop()
			}
		}
		return event
	})

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab: // When Tab or Down arrow key is pressed
			if app.GetFocus() == container.Table {
				app.SetFocus(image.Table)
			} else {
				app.SetFocus(container.Table)
			}
		}
		return event
	})

	// Run the application
	if err := app.SetRoot(layout, true).Run(); err != nil {
		panic(err)
	}

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
