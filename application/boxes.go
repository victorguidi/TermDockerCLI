package application

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (b *FlexBox) build() {

	page := fmt.Sprintf(" Extras - [%d/%d] ", b.CurrentPage, len(b.Tabs)-1)
	b.SetBorder(true)
	b.SetTitle(page)
	b.SetTitleAlign(tview.AlignLeft)
	b.AddInputCommands()
	b.Clear()

	switch b.CurrentPage {
	case 0:
		if b.Tabs[0] == nil {
			b.ImagesPage()
		} else {
			b.AddItem(b.Tabs[0], 0, 1, false)
		}
	case 1:
		if b.Tabs[1] == nil {
			b.NetworksPage()
		} else {
			b.AddItem(b.Tabs[1], 0, 1, false)
		}
	case 2:
		if b.Tabs[2] == nil {
			b.VolumesPage()
		} else {
			b.AddItem(b.Tabs[2], 0, 1, false)
		}
	default:
		b.AddItem(b.Tabs[b.CurrentPage], 0, 1, false)
	}
}

func (b *FlexBox) AddInputCommands() {
	// Function to handle the keyboard events, h go to the previous page, l go to the next page
	b.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case 'h':
				if b.CurrentPage > 0 {
					b.CurrentPage--
				}
			case 'l':
				if b.CurrentPage < len(b.Tabs)-1 {
					b.CurrentPage++
				}
			}
		}
		b.build()
		return event
	})
}

func (b *FlexBox) ImagesPage() {
	b.CurrentPage = 0

	b.Tabs[0] = tview.NewTable()
	b.Tabs[0].SetCell(0, 0, tview.NewTableCell("ID").SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
	b.Tabs[0].SetCell(0, 1, tview.NewTableCell("Repository").SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
	b.Tabs[0].SetCell(0, 2, tview.NewTableCell("Tag").SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
	b.Tabs[0].SetCell(0, 3, tview.NewTableCell("Size").SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))

	b.build()
}

func (b *FlexBox) NetworksPage() {

	b.CurrentPage = 1

	b.Tabs[1] = tview.NewTable()
	b.Tabs[1].SetCell(0, 0, tview.NewTableCell("ID").SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
	b.Tabs[1].SetCell(0, 1, tview.NewTableCell("Name").SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
	b.Tabs[1].SetCell(0, 2, tview.NewTableCell("Driver").SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
	b.Tabs[1].SetCell(0, 3, tview.NewTableCell("Scope").SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
	b.Tabs[1].SetCell(0, 4, tview.NewTableCell("Internal").SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))

	b.build()
}

func (b *FlexBox) VolumesPage() {

	b.CurrentPage = 2

	b.Tabs[2] = tview.NewTable()
	b.Tabs[2].SetCell(0, 0, tview.NewTableCell("Driver").SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
	b.Tabs[2].SetCell(0, 1, tview.NewTableCell("Name").SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
	b.Tabs[2].SetCell(0, 2, tview.NewTableCell("Mountpoint").SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))

	b.build()

}
