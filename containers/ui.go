package containers

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ContainerUi struct {
	Table   *tview.Table
	Options []string
}

func NewContainerUi() *ContainerUi {
	return &ContainerUi{
		Table:   tview.NewTable(),
		Options: []string{"start", "stop"},
	}
}

func (c *ContainerUi) PopulateUi() {
	c.Table.SetBorder(true).SetTitle("Docker Containers")
	c.Table.SetCell(0, 0, tview.NewTableCell("ID").SetTextColor(tcell.ColorYellow).SetSelectable(false))
	c.Table.SetCell(0, 1, tview.NewTableCell("Image").SetTextColor(tcell.ColorYellow).SetSelectable(false))

	// c.Table.SetCell(1, 0, tview.NewTableCell("1b323bb1j").SetTextColor(tcell.ColorGreen))
	// c.Table.SetCell(1, 1, tview.NewTableCell("ubuntu").SetTextColor(tcell.ColorGreen))
	// c.Table.SetCell(2, 0, tview.NewTableCell("18kkasd12").SetTextColor(tcell.ColorGreen))
	// c.Table.SetCell(2, 1, tview.NewTableCell("arch").SetTextColor(tcell.ColorGreen))

	c.Table.SetFixed(1, 1)
	c.Table.SetSelectable(true, false)
	c.Table.Select(1, 1)
}
