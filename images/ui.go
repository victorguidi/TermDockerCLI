package images

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ImageUi struct {
	Table   *tview.Table
	Options []string
}

func NewImageUi() *ImageUi {
	return &ImageUi{
		Table:   tview.NewTable(),
		Options: []string{"build", "rm"},
	}
}

func (i *ImageUi) PopulateUi() {
	i.Table.SetBorder(true).SetTitle("Docker Images")
	i.Table.SetCell(0, 0, tview.NewTableCell("Name").SetTextColor(tcell.ColorYellow).SetSelectable(false))
	i.Table.SetCell(0, 1, tview.NewTableCell("Tag").SetTextColor(tcell.ColorYellow).SetSelectable(false))

	// i.Table.SetCell(1, 0, tview.NewTableCell("1b323bb1j").SetTextColor(tcell.ColorGreen))
	// i.Table.SetCell(1, 1, tview.NewTableCell("ubuntu").SetTextColor(tcell.ColorGreen))
	// i.Table.SetCell(2, 0, tview.NewTableCell("18kkasd12").SetTextColor(tcell.ColorGreen))
	// i.Table.SetCell(2, 1, tview.NewTableCell("arch").SetTextColor(tcell.ColorGreen))

	i.Table.SetFixed(1, 1)
	i.Table.SetSelectable(true, false)
	i.Table.Select(1, 1)
}
