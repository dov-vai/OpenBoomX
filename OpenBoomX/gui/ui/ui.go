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
	"gioui.org/x/component"
	"log"
	"obx/gui/components"
	"obx/gui/controllers"
	"obx/gui/routes"
	"obx/gui/services"
	"obx/gui/testing"
	"obx/gui/theme"
	"obx/protocol"
	"obx/utils/bluetooth"
	"time"
)

var defaultMargin = unit.Dp(10)

type UI struct {
	Theme             *material.Theme
	ButtonTheme       *material.Theme
	EqButtons         *components.EqButtons
	LightButtons      *components.LightButtons
	LightPicker       *components.LightPicker
	BeepSlider        *components.StepSlider
	OffButton         *components.OffButton
	PairingButtons    *components.PairingButtons
	ShutdownSlider    *components.StepSlider
	EqSlider          *components.EqSlider
	NavigationBar     *components.NavigationBar
	StatusBar         *components.StatusBar
	EqSaveButton      *components.EqSaveButton
	PresetButtons     *components.PresetButtons
	EqPresetService   *services.EqPresetService
	SpeakerController *controllers.SpeakerController
	SpeakerClient     protocol.ISpeakerClient
	Loaded            bool
	Error             error
	RetryConnection   widget.Clickable
	CurrentRoute      routes.AppRoute
}

func NewUI() *UI {
	ui := &UI{}
	_th := material.NewTheme()
	_th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	th := _th.WithPalette(theme.Palette)
	ui.Theme = &th
	btnTheme := _th.WithPalette(theme.ButtonPalette)
	ui.ButtonTheme = &btnTheme
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
	ui.EqPresetService = services.NewEqPresetService()
	ui.EqButtons = components.CreateEQButtons(ui.SpeakerController.OnModeClicked)
	ui.LightButtons = components.CreateLightButtons(ui.SpeakerController.OnActionClicked)
	ui.LightPicker = components.CreateLightPicker(ui.SpeakerController.OnColorChanged)
	ui.BeepSlider = components.CreateBeepSlider(5, "Beep Volume", ui.SpeakerController.OnBeepStepChanged)
	ui.OffButton = components.CreateOffButton(ui.SpeakerController.OnOffButtonClicked)
	ui.ShutdownSlider = components.CreateBeepSlider(7, "Shutdown Timeout", ui.SpeakerController.OnShutdownStepChanged)
	ui.PairingButtons = components.CreatePairingButtons(ui.SpeakerController.OnPairingOn, ui.SpeakerController.OnPairingOff)
	ui.NavigationBar = components.CreateNavigationBar(func(route routes.AppRoute) {
		ui.CurrentRoute = route
	})
	ui.StatusBar = components.CreateStatusBar()
	ui.updateBattery()

	ui.EqSlider = components.CreateEqSlider(ui.SpeakerController.OnEqValuesChanged)
	// set currently active preset if it exists
	activePreset := ui.EqPresetService.GetActivePreset()
	if activePreset != "" {
		eqValues, _ := ui.EqPresetService.GetPresetValues(activePreset)
		err := ui.EqSlider.SetSliderValues(eqValues)
		if err != nil {
			log.Println(err)
		}
	}
	ui.EqPresetService.RegisterListener(ui.EqSlider)

	ui.EqSaveButton = components.CreateEqSaveButton(func(title string) {
		err := ui.EqPresetService.AddPreset(title, ui.EqSlider.GetSliderValues())
		if err != nil {
			log.Println(err)
		}
	})
	ui.EqSaveButton.SetText(activePreset)
	ui.EqPresetService.RegisterListener(ui.EqSaveButton)

	ui.PresetButtons = components.CreatePresetButtons(ui.EqPresetService)
	ui.EqPresetService.RegisterListener(ui.PresetButtons)
	ui.CurrentRoute = routes.Oluv
	ui.Loaded = true
}

func (ui *UI) updateBattery() {
	updateChannel := make(chan int)

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		batteryLevel, _ := ui.SpeakerClient.ReadBatteryLevel()
		updateChannel <- batteryLevel

		for range ticker.C {
			batteryLevel, err := ui.SpeakerClient.ReadBatteryLevel()
			if err != nil {
				fmt.Println("Error reading battery level:", err)
			} else {
				updateChannel <- batteryLevel
			}
		}
	}()

	go func() {
		for batteryLevel := range updateChannel {
			ui.StatusBar.BatteryLevel = batteryLevel
		}
	}()
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

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			surfaceStyle := component.Surface(
				&material.Theme{
					Palette: material.Palette{
						Bg: ui.Theme.Bg,
					},
				})

			surfaceStyle.CornerRadius = 0

			return surfaceStyle.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Dimensions{Size: gtx.Constraints.Max}
			})
		}),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				if !ui.Loaded {

					return ui.loadingLayout(gtx)
				}
				return ui.homeLayout(gtx)
			})
		}),
	)
}

func (ui *UI) Dispose() {
	if ui.SpeakerClient != nil {
		ui.SpeakerClient.CloseConnection()
	}
}
