package containers

import (
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ContainerUi struct {
	Table   *tview.Table
	Logs    chan []byte
	Options []string
}

func NewContainerUi() *ContainerUi {
	return &ContainerUi{
		Table:   tview.NewTable(),
		Logs:    make(chan []byte, 100),
		Options: []string{"start", "stop"},
	}
}

func (c *ContainerUi) PopulateUi(containers []DockerContainer) {

	c.Table.SetBorder(true).SetTitle("Docker Containers")
	c.Table.SetCell(0, 0, tview.NewTableCell("ID").SetTextColor(tcell.ColorYellow).SetSelectable(false))
	c.Table.SetCell(0, 1, tview.NewTableCell("Image").SetTextColor(tcell.ColorYellow).SetSelectable(false))

	for i, container := range containers {
		c.Table.SetCell(i+1, 0, tview.NewTableCell(container.ContainerId).SetTextColor(tcell.ColorGreen))
		c.Table.SetCell(i+1, 1, tview.NewTableCell(container.Image).SetTextColor(tcell.ColorGreen))
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go GetLogs(c.Logs, containers[0].ContainerId, wg)
	wg.Wait()

	c.Table.SetSelectionChangedFunc(func(row, column int) {
		containerId := c.Table.GetCell(row, 0).Text
		wg.Add(1)
		go GetLogs(c.Logs, containerId, wg)
		wg.Wait()
	})

	c.Table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case 'i':
				containerId := c.GetSelectedContainer()
				wg.Add(1)
				go Inspect(c.Logs, containerId, wg)
				wg.Wait()
			}
		}
		return event
	})

	c.Table.SetFixed(1, 1)
	c.Table.SetSelectable(true, false)
	c.Table.Select(1, 1)
}

func (c *ContainerUi) GetSelectedContainer() string {
	row, _ := c.Table.GetSelection()
	return c.Table.GetCell(row, 0).Text
}
