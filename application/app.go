package application

import (
	"fmt"

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
	dockerInfo = DockerInfo{
		TextView: tview.NewTextView().SetDynamicColors(true).SetRegions(true).SetScrollable(true),
		Data:     make(chan any),
	}
	body       = tview.NewFlex().SetDirection(tview.FlexColumn)
	leftPanel  = tview.NewFlex().SetDirection(tview.FlexRow)
	rightPanel = tview.NewFlex().SetDirection(tview.FlexRow)
)

func init() {
	if yerr != nil {
		panic(yerr)
	}
	flexBox.build()
}

func NewApplication() *Application {
	return &Application{
		Application:   tview.NewApplication(),
		Windows:       [9]chan any{},
		CurrentWindow: 0,
	}
}

func (a *Application) Build() {
	a.AddInputCommands()
	for i := 0; i < len(remoteHosts.Hosts) && i < 9; i++ {
		a.Windows[i] = make(chan any)
	}

	// Header is a flexbox with tabs for each remote host
	header := tview.NewFlex().SetDirection(tview.FlexColumn)
	for i := 0; i < len(remoteHosts.Hosts) && i < 9; i++ {
		button := tview.NewButton(remoteHosts.Hosts[i].IP)
		button.SetBorder(true)
		header.AddItem(button, 0, 1, false)
	}

	leftPanel.AddItem(containers.SetTitle(fmt.Sprintf("Containers - [10] / tab - [%d/%d]", a.CurrentWindow, len(a.Windows)-1)).SetBorder(true), 0, 1, true)
	leftPanel.AddItem(flexBox, 0, 1, true)

	rightPanel.AddItem(dockerInfo.SetTitle("Docker Info").SetBorder(true), 0, 1, true)

	body.AddItem(leftPanel, 0, 2, true)
	body.AddItem(rightPanel, 0, 3, true)

	// Layout will is a flexbox with a header and a body
	layout := tview.NewFlex().SetDirection(tview.FlexRow)
	layout.AddItem(header, 1, 1, false)
	layout.AddItem(body, 0, 1, true)

	a.SetRoot(layout, true)
}

func (a *Application) AddInputCommands() {
	a.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case 'J': // shift+j move down
				a.SetFocus(flexBox)
			case 'K': // shift+k move up
				a.SetFocus(containers)
			}
		}
		if event.Key() == tcell.KeyTab {
			// Switch between the left panel and the right panel with Tab
			if a.GetFocus() == containers || a.GetFocus() == flexBox {
				a.SetFocus(dockerInfo)
			} else {
				a.SetFocus(containers)
			}
		}
		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case '0':
				a.CurrentWindow = 0
			case '1':
				a.CurrentWindow = 1
			case '2':
				a.CurrentWindow = 2
			case '3':
				a.CurrentWindow = 3
			case '4':
				a.CurrentWindow = 4
			case '5':
				a.CurrentWindow = 5
			case '6':
				a.CurrentWindow = 6
			case '7':
				a.CurrentWindow = 7
			case '8':
				a.CurrentWindow = 8
			}
			containers.SetTitle(fmt.Sprintf("Containers - [10] / tab - [%d/%d]", a.CurrentWindow, len(a.Windows)-1))
		}
		return event
	})
}
