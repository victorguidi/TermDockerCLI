package application

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/victorguidi/TermDockerCLI/utils"
)

var (
	remoteHosts, yerr = utils.ReadYml("config.yml")
	flexBox           = FlexBox{
		Flex:        tview.NewFlex().SetDirection(tview.FlexRow),
		Data:        make(chan any),
		Tabs:        [3]*tview.Table{},
		CurrentPage: 0,
	}
	containers = tview.NewBox().SetTitle("Containers").SetBorder(true).SetTitleAlign(tview.AlignLeft)
)

func init() {
	if yerr != nil {
		panic(yerr)
	}
	flexBox.build()
}

func NewApplication() *Application {
	return &Application{
		Application: tview.NewApplication(),
	}
}

func (a *Application) Build() {
	a.AddInputCommands()
	for i := 0; i < len(remoteHosts.Hosts) && i < 9; i++ {
		a.Windows[i] = make(chan interface{})
	}

	// Header is a flexbox with tabs for each remote host
	header := tview.NewFlex().SetDirection(tview.FlexColumn)
	for i := 0; i < len(remoteHosts.Hosts) && i < 9; i++ {
		button := tview.NewButton(remoteHosts.Hosts[i].IP)
		button.SetBorder(true)
		header.AddItem(button, 0, 1, false)
	}

	leftPanel := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(containers, 0, 1, false).
		AddItem(flexBox, 0, 1, true)

	rightPanel := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewBox().SetTitle("Container Details").SetBorder(true), 0, 1, false)

	body := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(leftPanel, 0, 2, true).
		AddItem(rightPanel, 0, 3, false)

	// Layout will is a flexbox with a header and a body
	layout := tview.NewFlex().SetDirection(tview.FlexRow)
	layout.AddItem(header, 1, 1, false)
	layout.AddItem(body, 0, 1, true)

	a.SetRoot(layout, true)
}

func (a *Application) AddInputCommands() {
	// Function to handle the keyboard events, h go to the previous page, l go to the next page
	a.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case 'j': // j move down
				a.SetFocus(flexBox)
			case 'k': // k move up
				a.SetFocus(containers)
			}
		}
		return event
	})

	// TODO: Add TAB to switch between the right panel and the left panel
}
