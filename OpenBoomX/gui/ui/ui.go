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
	"obx/gui/services"
	"obx/gui/testing"
	"obx/gui/theme"
	"obx/protocol"
	"obx/utils/bluetooth"
)

var defaultMargin = unit.Dp(10)

type UI struct {
	theme              *material.Theme
	buttonTheme        *material.Theme
	snackbar           *components.Snackbar
	eqPresetService    *services.EqPresetService
	colorPresetService *services.ColorPresetService
	speakerController  *controllers.SpeakerController
	speakerClient      protocol.ISpeakerClient
	homePage           *pages.HomePage
	loaded             bool
	appError           error
	retryConnection    widget.Clickable
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
	// go ui.connectTestSpeaker()
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

	ui.homePage = pages.NewHomePage(
		ui.theme,
		ui.buttonTheme,
		ui.speakerController,
		ui.eqPresetService,
		ui.colorPresetService,
		ui.snackbar,
		func(err error) {
			ui.appError = err
			ui.loaded = false
		},
	)

	ui.loaded = true
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
	ui.homePage.Update(gtx)
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
				return ui.homePage.Layout(gtx)
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
