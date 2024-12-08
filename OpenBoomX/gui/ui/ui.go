package ui

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"obx/gui/components"
	"obx/gui/controllers"
	"obx/gui/testing"
	"obx/protocol"
	"obx/utils/bluetooth"
)

var defaultMargin = unit.Dp(10)

type UI struct {
	Theme             *material.Theme
	EqButtons         *components.EqButtons
	LightPicker       *components.LightPicker
	BeepSlider        *components.StepSlider
	OffButton         *components.OffButton
	PairingButtons    *components.PairingButtons
	ShutdownSlider    *components.StepSlider
	EqSlider          *components.EqSlider
	SpeakerController *controllers.SpeakerController
	SpeakerClient     protocol.ISpeakerClient
	Loaded            bool
	Error             error
	RetryConnection   widget.Clickable
}

func NewUI() *UI {
	ui := &UI{}
	ui.Theme = material.NewTheme()
	ui.Theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	go ui.connectSpeaker()
	//go ui.connectTestSpeaker()
	return ui
}

func (ui *UI) connectTestSpeaker() {
	client := &testing.MockSpeakerClient{}
	ui.initialize(client)
}

// TODO: should this be in ui? maybe speaker client could handle it
func (ui *UI) connectSpeaker() {
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
	ui.EqSlider = components.CreateEqSlider(ui.SpeakerController.OnEqValuesChanged)
	ui.Loaded = true
}

func (ui *UI) Run(w *app.Window) error {
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
		return ui.homeLayout(gtx)
	})
}

func (ui *UI) Dispose() {
	if ui.SpeakerClient != nil {
		ui.SpeakerClient.CloseConnection()
	}
}
