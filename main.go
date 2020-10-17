package main

import (
	"flag"
	"fmt"
	"fyne.io/fyne"
	 "fyne.io/fyne/app"
	 "fyne.io/fyne/widget"
	"log"
	"math"
	"sort"
	"strconv"

	"github.com/blackjack/webcam"
)

var uiTestFlag = flag.Bool("test_ui", false, "Fill UI with test values, used for UI development.")

type ControlType int

const (
	CTBool ControlType = iota
	CTInt
)

type Device struct{
	name, devPath string
	webcam *webcam.Webcam
	controls map[webcam.ControlID]webcam.Control
}

// List available devices and populate it.
func populateDevices() []*Device {
	// Populate devices.
	devicePaths, err := webcam.ListDevices()
	if err != nil {
		log.Fatalf("Failed list devices: %v", err)
	}
	devices := make([]*Device, 0, len(devicePaths))
	for devPath, name := range devicePaths {
		fmt.Printf("Device %q at %s\n", name, devPath)
		w, err := webcam.Open(devPath)
		if err != nil {
			log.Fatalf("Failed to open device %q at %s: %v", name, devPath, err)
		}
		dev := &Device{
			name: name,
			devPath: devPath,
			webcam: w,
			controls: w.GetControls(),
		}
		devices = append(devices, dev)
	}
	return devices
}

// Closes open devices.
func closeDevices(devices []*Device) {
	for _, dev := range devices {
		if dev.webcam != nil {
			dev.webcam.Close()
		}
	}
}

var testData []*Device = []*Device{
	&Device{
		name: "TestDevice",
		devPath: "/dev/null",
		controls: map[webcam.ControlID]webcam.Control{
			1: webcam.Control{"Brightness", 0, 100},
		},
	},
	&Device{
		name:    "FancyWebcam",
		devPath: "/dev/null",
		controls: map[webcam.ControlID]webcam.Control{
			1: webcam.Control{"Brightness", 0, 100},
			2: webcam.Control{"Exposure, Auto Priority", 0, 1},
			3: webcam.Control{"Some categorical control", 0, 2},
		},
	},
}

// UI holds all Fyne UI elements.
type UI struct {
	app  fyne.App
	win  fyne.Window
	tabs *widget.TabContainer
}

func makeDeviceLayout(dev *Device) (fyne.Widget, error) {
	var parts []fyne.CanvasObject

	// Sort controls in alphabetical order.
	type ControlAndId struct {
		id webcam.ControlID
		control webcam.Control
	}
	orderedControls := make([]ControlAndId, 0, len(dev.controls))
	for id, control := range dev.controls {
		orderedControls = append(orderedControls, ControlAndId{id, control})
	}
	sort.Slice(orderedControls, func(i, j int) bool {
		return orderedControls[i].control.Name < orderedControls[j].control.Name
	})

	// Add controls.
	for _, cid := range orderedControls {
		part, err := makeControlLayout(dev, cid.id, cid.control)
		if err != nil {
			return nil, err
		}
		parts = append(parts, part)
	}
	return widget.NewVBox(parts...), nil
}

func makeControlLayout(dev *Device, id webcam.ControlID, control webcam.Control) (fyne.CanvasObject, error) {
	numChoices := control.Max - control.Min + 1
	currentValue, err := dev.webcam.GetControl(id)
	if err != nil {
		return nil, fmt.Errorf("failed to set control %q (id=%d): %v", control.Name, id, err)
	}
	if numChoices == 2 {
		// Boolean checkbox.
		checkBox := widget.NewCheck(control.Name,
			func(checked bool) {
				value := control.Min
				if checked {
					value = control.Max
				}
				err := dev.webcam.SetControl(id, value)
				if err != nil {
					log.Printf("Failed to set control %q (id=%d) to %d: %v", control.Name, id, value, err)
				}
				fmt.Printf("Set %q to %v", control.Name, checked)
			})
		checkBox.SetChecked(currentValue == control.Max)
		return checkBox, nil

	} else if numChoices <= 4 {
		// Assume categorical, use radio buttons.
		var parts []fyne.CanvasObject
		parts = append(parts, widget.NewLabelWithStyle(control.Name, fyne.TextAlignTrailing,
			fyne.TextStyle{}))
		options := make([]string, 0, numChoices)
		for ii := control.Min; ii <= control.Max; ii++ {
			options = append(options, strconv.Itoa(int(ii)))
		}
		radioButtons := widget.NewRadio(options, func(s string) {
			value, err := strconv.Atoi(s)
			if err != nil {
				log.Printf("Invalid value for control %q (id=%d): %q", control.Name, id, s)
				return
			}
			if err := dev.webcam.SetControl(id, int32(value)); err != nil {
				log.Printf("Failed to set control %q (id=%d) to %d: %v", control.Name, id, value, err)
			}
			fmt.Printf("Set %q to %v\n", control.Name, value)
		})
		radioButtons.SetSelected(strconv.Itoa(int(currentValue)))
		parts = append(parts, radioButtons)
		return widget.NewHBox(parts...), nil

	}

	// Use slider to select value.
	label := widget.NewLabel(
		fmt.Sprintf("%s (%d to %d): ", control.Name, control.Min, control.Max))
	labelValue := widget.NewLabelWithStyle(strconv.Itoa(int(currentValue)), fyne.TextAlignLeading,
		fyne.TextStyle{Monospace: true})
	labels := widget.NewHBox(label, labelValue)
	slider := widget.NewSlider(float64(control.Min), float64(control.Max))
	slider.Value = float64(currentValue)
	slider.OnChanged = func(x float64) {
		value := int32(math.Round(x))
		if value <= control.Min {
			value = control.Min
		} else if value >= control.Max {
			value = control.Max
		}
		if err := dev.webcam.SetControl(id, int32(value)); err != nil {
			log.Printf("Failed to set control %q (id=%d) to %d: %v", control.Name, id, value, err)
		}
		fmt.Printf("Set %q to %v\n", control.Name, value)
		labelValue.SetText(strconv.Itoa(int(value)))
	}
	container := widget.NewVBox(labels, slider)
	return container, nil
}


func NewUI(devices []*Device) (*UI, error) {
	ui := &UI{}
	ui.app = app.New()
	ui.win = ui.app.NewWindow("VideoController")

	tabItems := make([]*widget.TabItem, 0, len(devices))
	for _, dev := range devices {
		tabLayout, err := makeDeviceLayout(dev)
		if err != nil {
			return nil, fmt.Errorf("failed to create layout for device %q: %v", dev.name, err)
		}
		tab := widget.NewTabItem(dev.name, tabLayout)
		tabItems = append(tabItems, tab)
	}
	ui.tabs = widget.NewTabContainer(tabItems...)
	ui.win.SetContent(ui.tabs)
	return ui, nil
}

func (ui *UI) Run() {
	ui.win.ShowAndRun()
}

func main() {
	var devices []*Device
	if *uiTestFlag {
		devices = testData
	} else {
		devices = populateDevices()
	}
	defer closeDevices(devices)

	// Run UI
	ui, err := NewUI(devices)
	if err != nil {
		log.Fatal(err)
	}
	ui.Run()
}
