package application

import "github.com/rivo/tview"

type Application struct {
	*tview.Application
	Windows [9]chan interface{}
}

type FlexBox struct {
	*tview.Flex
	Tabs        [3]*tview.Table
	CurrentPage int
	Data        chan any
}
