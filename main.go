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

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	dockerImages     chan []DockerImage
	dockerContainers chan []DockerContainer
	dockerNetworks   chan []DockerNetwork
	logs             string
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

	container := containers.ContainerUi{
		Table:   tview.NewTable(),
		Options: []string{"start", "stop"},
	}

	container.Table.SetBorder(true).SetTitle("Docker Containers")
	container.Table.SetCell(0, 0, tview.NewTableCell("ID").SetTextColor(tcell.ColorYellow).SetSelectable(false))
	container.Table.SetCell(0, 1, tview.NewTableCell("Image").SetTextColor(tcell.ColorYellow).SetSelectable(false))

	container.Table.SetCell(1, 0, tview.NewTableCell("1b323bb1j").SetTextColor(tcell.ColorGreen))
	container.Table.SetCell(1, 1, tview.NewTableCell("ubuntu").SetTextColor(tcell.ColorGreen))
	container.Table.SetCell(2, 0, tview.NewTableCell("18kkasd12").SetTextColor(tcell.ColorGreen))
	container.Table.SetCell(2, 1, tview.NewTableCell("arch").SetTextColor(tcell.ColorGreen))

	container.Table.SetFixed(1, 1)
	container.Table.SetSelectable(true, false)
	container.Table.Select(1, 1)

	container.Table.SetSelectedFunc(func(row, column int) {
		containerID := container.Table.GetCell(row, 0).Text
		logs = "Container ID: " + containerID + "\n"
	})

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
		AddItem(container.Table, 0, 1, true)

	logs = "oi"
	text := tview.NewTextView()
	text.SetBorder(true).SetTitle("Docker Logs")
	text.SetText(logs)
	rightPanel := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(text, 0, 1, true)

	// Create the layout
	layout := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(leftPanel, 0, 2, true).
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
