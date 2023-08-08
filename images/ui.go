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

func (i *ImageUi) PopulateUi(images []DockerImage) {
	i.Table.SetBorder(true).SetTitle("Docker Images")
	i.Table.SetCell(0, 0, tview.NewTableCell("Repository").SetTextColor(tcell.ColorYellow).SetSelectable(false))
	i.Table.SetCell(0, 1, tview.NewTableCell("Tag").SetTextColor(tcell.ColorYellow).SetSelectable(false))
	i.Table.SetCell(0, 2, tview.NewTableCell("ID").SetTextColor(tcell.ColorYellow).SetSelectable(false))
	i.Table.SetCell(0, 3, tview.NewTableCell("Size").SetTextColor(tcell.ColorYellow).SetSelectable(false))

	for index, image := range images {
		i.Table.SetCell(index+1, 0, tview.NewTableCell(image.Repository).SetTextColor(tcell.ColorGreen))
		i.Table.SetCell(index+1, 1, tview.NewTableCell(image.Tag).SetTextColor(tcell.ColorGreen))
		i.Table.SetCell(index+1, 2, tview.NewTableCell(image.ImageId).SetTextColor(tcell.ColorGreen))
		i.Table.SetCell(index+1, 3, tview.NewTableCell(image.Size).SetTextColor(tcell.ColorGreen))
	}

	i.Table.SetFixed(1, 1)
	i.Table.SetSelectable(true, false)
	i.Table.Select(1, 1)
}
