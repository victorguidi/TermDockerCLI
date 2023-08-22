package application

import (
	"fmt"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/victorguidi/TermDockerCLI/containers"
	"github.com/victorguidi/TermDockerCLI/utils"
)

// TODO: Load data from dockers
var (
	remoteHosts, yerr = utils.ReadYml("config.yml")

	flexBox = FlexBox{
		Flex:        tview.NewFlex().SetDirection(tview.FlexRow),
		Data:        make(chan any),
		Tabs:        [3]*tview.Table{},
		CurrentPage: 0,
	}

	dcontainers      = containers.NewContainerUi()
	assh             = containers.NewSSH(remoteHosts)
	dockerContainers = []containers.DockerContainer{}

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
		CurrentWindow: 0,
	}
}

func (a *Application) PopulateWindows() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	if a.CurrentWindow > 0 {
		dcontainers.PopulateUi(<-containers.GetAllContainers(assh[a.CurrentWindow].Host, assh[a.CurrentWindow], wg), assh[a.CurrentWindow])
	} else {
		dcontainers.PopulateUi(<-containers.GetAllContainers("local", nil, wg), nil)
	}
	wg.Wait()
}

func (a *Application) Build() {
	a.AddInputCommands()
	a.PopulateWindows()

	leftPanel.AddItem(dcontainers.Table, 0, 1, true)
	leftPanel.AddItem(flexBox, 0, 1, true)

	rightPanel.AddItem(dockerInfo.SetTitle("Docker Info").SetBorder(true), 0, 1, true)

	body.AddItem(leftPanel, 0, 2, true)
	body.AddItem(rightPanel, 0, 3, true)

	// Layout will is a flexbox with a header and a body
	layout := tview.NewFlex().SetDirection(tview.FlexRow)
	layout.AddItem(body, 0, 1, true)

	a.SetRoot(layout, true)
}

// Add Input Commands to the application
func (a *Application) AddInputCommands() {
	a.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case 'J': // shift+j move down
				a.SetFocus(flexBox)
			case 'K': // shift+k move up
				a.SetFocus(dcontainers.Table)
			}
		}

		if event.Key() == tcell.KeyTab {
			// Switch between the left panel and the right panel with Tab
			if a.GetFocus() == dcontainers.Table || a.GetFocus() == flexBox {
				a.SetFocus(dockerInfo)
			} else {
				a.SetFocus(dcontainers.Table)
			}
		}

		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case '0':
				a.CurrentWindow = 0
			case '1':
				if len(remoteHosts.Hosts) < 1 {
					break
				}
				a.CurrentWindow = 1
			case '2':
				if len(remoteHosts.Hosts) < 2 {
					break
				}
				a.CurrentWindow = 2
			case '3':
				if len(remoteHosts.Hosts) < 3 {
					break
				}
				a.CurrentWindow = 3
			case '4':
				if len(remoteHosts.Hosts) < 4 {
					break
				}
				a.CurrentWindow = 4
			case '5':
				if len(remoteHosts.Hosts) < 5 {
					break
				}
				a.CurrentWindow = 5
			case '6':
				if len(remoteHosts.Hosts) < 6 {
					break
				}
				a.CurrentWindow = 6
			case '7':
				if len(remoteHosts.Hosts) < 7 {
					break
				}
				a.CurrentWindow = 7
			case '8':
				if len(remoteHosts.Hosts) < 8 {
					break
				}
				a.CurrentWindow = 8
			case '9':
				if len(remoteHosts.Hosts) < 9 {
					break
				}
				a.CurrentWindow = 9
			}

			if a.CurrentWindow == 0 {
				dcontainers.Table.SetTitle(fmt.Sprintf(" local: Containers-[10] "))
				a.PopulateWindows()
			} else {
				dcontainers.Table.SetTitle(fmt.Sprintf(" %s: Containers-[10] ", remoteHosts.Hosts[a.CurrentWindow-1].IP))
				a.PopulateWindows()
			}
		}
		return event
	})
}
