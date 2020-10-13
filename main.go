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
	name string
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

type UI struct {
	app fyne.App
	win fyne.Window
}

func NewUI(devices DevicesMap) *UI {
	ui := &UI{}
	ui.app = app.New()
	ui.win = ui.app.NewWindow("VideoController")

	return ui
}

func main() {

	a := app.New()
	w := a.NewWindow("Hello")

	hello := widget.NewLabel("Hello Fyne!")
	w.SetContent(widget.NewVBox(
		hello,
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome :)")
		}),
	))

	w.ShowAndRun()
}
