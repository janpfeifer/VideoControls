package main

import "fyne.io/fyne"
import "fyne.io/fyne/app"
import "fyne.io/fyne/widget"

type ControlType int

const (
	CTBool ControlType = iota
	CTInt
)

type Control struct {
	name  string
	ctype ControlType

	// maxValue is inclusive, should still be a valid value.
	// These are ignored for CTBool.
	minValue, maxValue int
}

type DevicesMap map[string][]*Control

var testData DevicesMap = DevicesMap{
	"TestDevice": []*Control{
		&Control{"Brightness", CTInt, 0, 100}},
	"FancyWebcam": []*Control{
		&Control{"Brightness", CTInt, 0, 100},
		&Control{"Brightness (Auto)", CTBool, 0, 0},
	},
}

// UI holds all Fyne UI elements.
type UI struct {
	app  fyne.App
	win  fyne.Window
	tabs *widget.TabContainer
}

func makeDeviceLayout(controls []*Control) fyne.Widget {
	var parts []fyne.CanvasObject
	for _, control := range controls {
		parts = append(parts, widget.NewLabel(control.name))
	}
	return widget.NewVBox(parts...)
}

func NewUI(devices DevicesMap) *UI {
	ui := &UI{}
	ui.app = app.New()
	ui.win = ui.app.NewWindow("VideoController")

	tabItems := make([]*widget.TabItem, 0, len(devices))
	for devName, controls := range devices {
		tab := widget.NewTabItem(devName, makeDeviceLayout(controls))
		tabItems = append(tabItems, tab)
	}
	ui.tabs = widget.NewTabContainer(tabItems...)
	ui.win.SetContent(ui.tabs)
	return ui
}

func (ui *UI) Run() {
	ui.win.ShowAndRun()
}

func main() {
	ui := NewUI(testData)
	ui.Run()
}
