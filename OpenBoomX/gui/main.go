package main

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"log"
	"obx/gui/components"
	"obx/gui/controllers"
	"obx/protocol"
	"obx/utils/bluetooth"
	"os"
)

func main() {
	ui := newUI()
	defer func() {
		if ui.SpeakerClient != nil {
			ui.SpeakerClient.CloseConnection()
		}
	}()

	go func() {
		// TODO: retry connection button
		address, err := bluetooth.GetUBoomXAddress()
		if err != nil {
			ui.Error = fmt.Errorf("Is speaker not connected?: %w", err)
			return
		}

		rfcomm, err := protocol.NewRfcommClient(address)
		if err != nil {
			ui.Error = fmt.Errorf("Is device already connected to speaker?: %w", err)
			return
		}

		client := protocol.NewSpeakerClient(rfcomm)
		ui.initialize(client)
	}()

	go func() {
		w := new(app.Window)
		w.Option(
			app.Title("OpenBoomX"),
			app.Size(unit.Dp(300), unit.Dp(750)),
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
	Theme             *material.Theme
	EqButtons         *components.EqButtons
	LightPicker       *components.LightPicker
	BeepSlider        *components.StepSlider
	OffButton         *components.OffButton
	PairingButtons    *components.PairingButtons
	ShutdownSlider    *components.StepSlider
	SpeakerController *controllers.SpeakerController
	SpeakerClient     protocol.ISpeakerClient
	Loaded            bool
	Error             error
}

func newUI() *UI {
	ui := &UI{}
	ui.Theme = material.NewTheme()
	ui.Theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	return ui
}

func (ui *UI) initialize(client protocol.ISpeakerClient) {
	ui.SpeakerClient = client
	ui.SpeakerController = controllers.NewSpeakerController(client)
	ui.EqButtons = components.CreateEQButtons(ui.SpeakerController.OnModeClicked)
	ui.LightPicker = components.CreateLightPicker(ui.SpeakerController.OnActionClicked, ui.SpeakerController.OnColorChanged)
	ui.BeepSlider = components.CreateBeepSlider(5, "Beep Volume", ui.SpeakerController.OnBeepStepChanged)
	ui.OffButton = components.CreateOffButton(ui.SpeakerController.OnOffButtonClicked)
	ui.ShutdownSlider = components.CreateBeepSlider(7, "Shutdown Timeout", ui.SpeakerController.OnShutdownStepChanged)
	ui.PairingButtons = components.CreatePairingButtons(ui.SpeakerController.OnPairingOn, ui.SpeakerController.OnPairingOff)
	ui.Loaded = true
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
	if !ui.Loaded {
		return
	}
	ui.BeepSlider.Update(gtx)
	ui.ShutdownSlider.Update(gtx)
}

func (ui *UI) layout(gtx layout.Context) layout.Dimensions {
	inset := layout.UniformInset(defaultMargin)
	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		if !ui.Loaded {
			return ui.loadingLayout(gtx)
		}
		return ui.appLayout(gtx)
	})
}

func (ui *UI) appLayout(gtx layout.Context) layout.Dimensions {
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

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ui.PairingButtons.Layout(ui.Theme, gtx)
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ui.OffButton.Layout(ui.Theme, gtx)
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ui.ShutdownSlider.Layout(ui.Theme, gtx)
		}),
	)
}

func (ui *UI) loadingLayout(gtx layout.Context) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				text := "Loading..."
				if ui.Error != nil {
					text = ui.Error.Error()
				}
				label := material.H5(ui.Theme, text)
				return label.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if ui.Error != nil {
					return layout.Dimensions{}
				}

				gtx.Constraints.Max.X = gtx.Dp(32)
				gtx.Constraints.Max.Y = gtx.Dp(32)
				return material.Loader(ui.Theme).Layout(gtx)
			}),
		)
	})
}
