package application

import (
	"github.com/rivo/tview"
	"github.com/victorguidi/TermDockerCLI/utils"
)

var (
	remoteHosts, yerr = utils.ReadYml("config.yml")
)

func init() {
	if yerr != nil {
		panic(yerr)
	}
}

type Application struct {
	*tview.Application
	Windows [9]chan interface{}
}

func NewApplication() *Application {
	return &Application{
		Application: tview.NewApplication(),
	}
}

func (a *Application) Build() {
	for i := 0; i < len(remoteHosts.Hosts) && i < 9; i++ {
		a.Windows[i] = make(chan interface{})
	}
}
