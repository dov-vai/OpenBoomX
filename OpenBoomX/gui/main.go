package main

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"image/color"
	"log"
	"obx/btutils"
	"obx/gui/components"
	"obx/protocol"
	"obx/utils"
	"os"
	"runtime"
	"strings"
	"time"
	"tinygo.org/x/bluetooth"
)

func main() {
	// TODO: add loading screen & error message if device not connected
	adapter := bluetooth.DefaultAdapter
	utils.Must("enable BLE stack", adapter.Enable())

	var address string
	if runtime.GOOS == "windows" {
		var err error
		address, err = btutils.FindDeviceAddress(adapter, protocol.UBoomXName2, 5*time.Second)
		utils.Must("find device", err)
		// FIXME: a hack for getting the correct MAC address of the device, because scanning on windows doesn't seem to work correctly
		address = strings.Replace(address, protocol.UBoomXOUI2, protocol.UBoomXOUI, 1)
	} else {
		var err error
		address, err = btutils.FindDeviceAddress(adapter, protocol.UBoomXName, 5*time.Second)
		utils.Must("find device", err)
	}

	var rfcomm protocol.RfcommClient = protocol.NewRfcommClient(address)
	client := protocol.NewSpeakerClient(rfcomm)

	ui := newUI(client)

	go func() {
		w := new(app.Window)
		w.Option(
			app.Title("OpenBoomX"),
			app.Size(unit.Dp(300), unit.Dp(700)),
		)
		if err := ui.run(w); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	app.Main()
}

var defaultMargin = unit.Dp(10)

type UI struct {
	Theme       *material.Theme
	EqButtons   components.EqButtons
	LightPicker components.LightPicker
	BeepSlider  components.StepSlider
}

func newUI(client *protocol.SpeakerClient) *UI {
	ui := &UI{}
	ui.Theme = material.NewTheme()
	ui.Theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	// FIXME: probably should be defining these functions somewhere else?
	ui.EqButtons = components.CreateEQButtons(func(mode string) {
		err := client.SetOluvMode(mode)
		if err != nil {
			log.Printf("SetOluvMode failed: %v", err)
		}
	})
	ui.LightPicker = components.CreateLightPicker(func(action string) {
		// FIXME: light changes too frequently for the speaker to handle,
		// also it sets it to dancing white color on gui launch
		err := client.HandleLightAction(action, false)
		if err != nil {
			log.Printf("HandleLightAction failed: %v", err)
		}
	},
		func(color color.NRGBA, solidColor bool) {
			err := client.HandleLightAction(utils.NrgbaToHex(color), solidColor)
			if err != nil {
				log.Printf("HandleLightAction failed: %v", err)
			}
		})
	ui.BeepSlider = components.CreateBeepSlider(5, "Beep Volume", func(step int) {
		err := client.SetBeepVolume(25 * step)
		if err != nil {
			log.Printf("SetBeepVolume failed: %v", err)
		}
	})
	return ui
}

func (ui *UI) run(w *app.Window) error {
	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			ui.update(gtx)
			ui.layout(gtx)
			e.Frame(gtx.Ops)

		case app.DestroyEvent:
			return e.Err
		}
	}
}

func (ui *UI) update(gtx layout.Context) {
	ui.BeepSlider.Update(gtx)
}

func (ui *UI) layout(gtx layout.Context) layout.Dimensions {
	inset := layout.UniformInset(defaultMargin)
	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.EqButtons.Layout(ui.Theme, gtx)
			}),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.LightPicker.Layout(ui.Theme, gtx)
			}),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.BeepSlider.Layout(ui.Theme, gtx)
			}),
		)
	})
}
