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
	"obx/gui/pages"
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
	theme              *material.Theme
	buttonTheme        *material.Theme
	navigationBar      *components.NavigationBar
	statusBar          *components.StatusBar
	snackbar           *components.Snackbar
	eqPresetService    *services.EqPresetService
	colorPresetService *services.ColorPresetService
	speakerController  *controllers.SpeakerController
	speakerClient      protocol.ISpeakerClient
	oluvPage           *pages.OluvPage
	eqPage             *pages.EqPage
	presetsPage        *pages.PresetsPage
	lightsPage         *pages.LightsPage
	miscPage           *pages.MiscPage
	loaded             bool
	appError           error
	retryConnection    widget.Clickable
	currentRoute       routes.AppRoute
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
	// add comments to line above and uncomment below
	// to connect a mock speaker for GUI development
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
	ui.snackbar = components.CreateSnackbar()
	ui.speakerController = controllers.NewSpeakerController(client, ui.snackbar)
	ui.eqPresetService = services.NewEqPresetService()
	ui.colorPresetService = services.NewColorPresetService()
	ui.navigationBar = components.CreateNavigationBar(func(route routes.AppRoute) {
		ui.currentRoute = route
	})
	ui.oluvPage = pages.NewOluvPage(ui.buttonTheme, ui.speakerController)
	ui.eqPage = pages.NewEqPage(ui.theme, ui.buttonTheme, ui.eqPresetService, ui.speakerController, ui.snackbar)
	ui.presetsPage = pages.NewPresetsPage(ui.buttonTheme, ui.eqPresetService, ui.snackbar)
	ui.lightsPage = pages.NewLightsPage(ui.theme, ui.buttonTheme, ui.speakerController, ui.colorPresetService, ui.snackbar)
	ui.miscPage = pages.NewMiscPage(ui.theme, ui.buttonTheme, ui.speakerController, ui.getFirmwareName())
	ui.statusBar = components.CreateStatusBar()
	ui.updateBattery()
	ui.currentRoute = routes.Oluv
	ui.loaded = true
}

func (ui *UI) getFirmwareName() string {
	firmware, err := ui.speakerClient.ReadFirmwarePackageName()
	if err != nil {
		log.Println(err)
	}
	return firmware
}

func (ui *UI) updateBattery() {
	updateChannel := make(chan int)

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		batteryLevel, _ := ui.speakerClient.ReadBatteryLevel()
		updateChannel <- batteryLevel

		for range ticker.C {
			batteryLevel, err := ui.speakerClient.ReadBatteryLevel()

			if err == nil {
				updateChannel <- batteryLevel
				continue
			}

			fmt.Println("Error reading battery level:", err)

			// handling for unix and windows if device disconnected
			if protocol.IsSocketDisconnected(err) {
				ui.appError = fmt.Errorf("Is speaker not connected?: %w", err)
				err = ui.speakerClient.CloseConnection()
				if err != nil {
					log.Printf("Error closing speaker connection: %v", err)
				}
				ui.loaded = false
				close(updateChannel)
				break
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
	ui.miscPage.Update(gtx)
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
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			if !ui.loaded {
				return layout.Dimensions{}
			}
			return ui.snackbar.Layout(ui.theme, gtx)
		}),
	)
}

func (ui *UI) Dispose() {
	if ui.speakerClient != nil {
		err := ui.speakerClient.CloseConnection()
		if err != nil {
			log.Printf("Error closing speaker connection: %v", err)
		}
	}
}
