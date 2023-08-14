package application

import "github.com/rivo/tview"

type Application struct {
	*tview.Application
	Windows       [9]chan any
	CurrentWindow int
}

type FlexBox struct {
	*tview.Flex
	Tabs        [3]*tview.Table
	CurrentPage int
	Data        chan any
}

type DockerInfo struct {
	*tview.TextView
	Data chan any
}
