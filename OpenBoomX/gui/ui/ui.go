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
	"obx/utils"
	"obx/utils/bluetooth"
	"time"
)

var defaultMargin = unit.Dp(10)

type UI struct {
	theme             *material.Theme
	buttonTheme       *material.Theme
	eqButtons         *components.EqButtons
	lightButtons      *components.LightButtons
	lightPicker       *components.LightPicker
	beepSlider        *components.StepSlider
	offButton         *components.OffButton
	pairingButtons    *components.PairingButtons
	shutdownSlider    *components.StepSlider
	eqSlider          *components.EqSlider
	navigationBar     *components.NavigationBar
	statusBar         *components.StatusBar
	eqSaveButton      *components.EqSaveButton
	eqResetButton     *components.EqResetButton
	presetButtons     *components.PresetButtons
	eqPresetService   *services.EqPresetService
	speakerController *controllers.SpeakerController
	speakerClient     protocol.ISpeakerClient
	loaded            bool
	appError          error
	retryConnection   widget.Clickable
	currentRoute      routes.AppRoute
}

func NewUI() *UI {
	ui := &UI{}
	_th := material.NewTheme()
	_th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	th := _th.WithPalette(theme.Palette)
	ui.theme = &th
	btnTheme := _th.WithPalette(theme.ButtonPalette)
	ui.buttonTheme = &btnTheme
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
		ui.appError = fmt.Errorf("Is speaker not connected?: %w", err)
		return
	}

	rfcomm, err := protocol.NewRfcommClient(address)
	if err != nil {
		ui.appError = fmt.Errorf("Is device already connected to speaker?: %w", err)
		return
	}

	client := protocol.NewSpeakerClient(rfcomm)
	ui.initialize(client)
}

func (ui *UI) initialize(client protocol.ISpeakerClient) {
	ui.speakerClient = client
	ui.speakerController = controllers.NewSpeakerController(client)
	ui.eqPresetService = services.NewEqPresetService()
	ui.eqButtons = components.CreateEQButtons(ui.speakerController.OnModeClicked)
	ui.lightButtons = components.CreateLightButtons(ui.speakerController.OnActionClicked)
	ui.lightPicker = components.CreateLightPicker(ui.speakerController.OnColorChanged)
	ui.beepSlider = components.CreateBeepSlider(5, "Beep Volume", utils.SortedKeysByValueInt(protocol.BeepVolumes), ui.speakerController.OnBeepStepChanged)
	ui.offButton = components.CreateOffButton(ui.speakerController.OnOffButtonClicked)
	ui.shutdownSlider = components.CreateBeepSlider(7, "Shutdown Timeout", utils.SortedKeysByValue(protocol.ShutdownTimeouts), ui.speakerController.OnShutdownStepChanged)
	ui.pairingButtons = components.CreatePairingButtons(ui.speakerController.OnPairingOn, ui.speakerController.OnPairingOff)
	ui.navigationBar = components.CreateNavigationBar(func(route routes.AppRoute) {
		ui.currentRoute = route
	})
	ui.statusBar = components.CreateStatusBar()
	ui.updateBattery()

	ui.eqSlider = components.CreateEqSlider(ui.speakerController.OnEqValuesChanged)
	// set currently active preset if it exists
	activePreset := ui.eqPresetService.GetActivePreset()
	if activePreset != "" {
		eqValues, _ := ui.eqPresetService.GetPresetValues(activePreset)
		err := ui.eqSlider.SetSliderValues(eqValues)
		if err != nil {
			log.Println(err)
		}
	}
	ui.eqPresetService.RegisterListener(ui.eqSlider)

	ui.eqResetButton = components.CreateEqResetButton(func() {
		err := ui.eqSlider.ResetValues()
		if err != nil {
			log.Println(err)
		}
	})

	ui.eqSaveButton = components.CreateEqSaveButton(func(title string) {
		err := ui.eqPresetService.AddPreset(title, ui.eqSlider.GetSliderValues())
		if err != nil {
			log.Println(err)
		}
	})
	ui.eqSaveButton.SetText(activePreset)
	ui.eqPresetService.RegisterListener(ui.eqSaveButton)

	ui.presetButtons = components.CreatePresetButtons(ui.eqPresetService)
	ui.eqPresetService.RegisterListener(ui.presetButtons)
	ui.currentRoute = routes.Oluv
	ui.loaded = true
}

func (ui *UI) updateBattery() {
	updateChannel := make(chan int)

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		batteryLevel, _ := ui.speakerClient.ReadBatteryLevel()
		updateChannel <- batteryLevel

		for range ticker.C {
			batteryLevel, err := ui.speakerClient.ReadBatteryLevel()
			if err != nil {
				fmt.Println("Error reading battery level:", err)
			} else {
				updateChannel <- batteryLevel
			}
		}
	}()

	go func() {
		for batteryLevel := range updateChannel {
			ui.statusBar.BatteryLevel = batteryLevel
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
	if !ui.loaded {
		return
	}
	ui.beepSlider.Update(gtx)
	ui.shutdownSlider.Update(gtx)
}

func (ui *UI) layout(gtx layout.Context) layout.Dimensions {
	inset := layout.UniformInset(defaultMargin)

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			surfaceStyle := component.Surface(
				&material.Theme{
					Palette: material.Palette{
						Bg: ui.theme.Bg,
					},
				})

			surfaceStyle.CornerRadius = 0

			return surfaceStyle.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Dimensions{Size: gtx.Constraints.Max}
			})
		}),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				if !ui.loaded {

					return ui.loadingLayout(gtx)
				}
				return ui.homeLayout(gtx)
			})
		}),
	)
}

func (ui *UI) Dispose() {
	if ui.speakerClient != nil {
		ui.speakerClient.CloseConnection()
	}
}
