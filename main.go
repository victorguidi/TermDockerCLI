package main

// TODO: Change the layout, is to ugly
// TODO: Return also the containers from ssh connection
// TODO: Add the possibility to send commands

import (
	"log"
	"os/exec"
	"strings"

	"github.com/victorguidi/TermDockerCLI/containers"
	"github.com/victorguidi/TermDockerCLI/images"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	dockerContainer = containers.NewContainer()
	container       = containers.NewContainerUi()
	image           = images.NewImageUi()
)

func init() {
	if !checkDockerIsInstalled() {
		log.Fatal("Docker is not installed, please install docker")
		return
	}
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

	// CallRoutines()

	containerChannel := make(chan []containers.DockerContainer)
	go dockerContainer.GetAllContainers(containerChannel)
	dcontainers := <-containerChannel // FIX: This might cause a deadlock?

	container.PopulateUi(dcontainers)
	image.PopulateUi()

	leftPanel := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(container.Table, 0, 1, true).
		AddItem(image.Table, 0, 1, true)

	text := tview.NewTextView()
	text.SetBorder(true).SetTitle("Docker TUI")

	go func() {
		for {
			select {
			case logs := <-container.Logs:
				text.SetText(string(logs))
			}
		}
	}() // This will start a loop that will wait for logs to be sent to the channel

	rightPanel := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(text, 0, 1, true)

	// Create the layout
	layout := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(leftPanel, 0, 2, true).
		AddItem(rightPanel, 0, 4, true)

	layout.SetBackgroundColor(tcell.ColorBlack)

	// Create the application
	app := tview.NewApplication()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			if app.GetFocus() == container.Table {
				app.SetFocus(image.Table)
			} else {
				app.SetFocus(container.Table)
			}
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q':
				app.Stop()
			}
		}
		return event
	})

	// Run the application
	if err := app.SetRoot(layout, true).Run(); err != nil {
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
