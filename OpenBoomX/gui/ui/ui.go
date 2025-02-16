package ui

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"log"
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
	eqPresetService    *services.EqPresetService
	colorPresetService *services.ColorPresetService
	speakerController  *controllers.SpeakerController
	speakerClient      protocol.ISpeakerClient
	homePage           *pages.HomePage
	loadingPage        *pages.LoadingPage
	loaded             bool
}

func NewUI() *UI {
	ui := &UI{}
	_th := material.NewTheme()
	_th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	th := _th.WithPalette(theme.Palette)
	ui.theme = &th
	btnTheme := _th.WithPalette(theme.ButtonPalette)
	ui.buttonTheme = &btnTheme
	ui.loadingPage = pages.NewLoadingPage(ui.buttonTheme, func() {
		go ui.connectSpeaker()
	})

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

func (ui *UI) connectSpeaker() {
	client, err := bluetooth.ConnectUBoomX()
	if err != nil {
		ui.loadingPage.SetError(err)
		return
	}
	ui.initialize(client)
}

func (ui *UI) initialize(client protocol.ISpeakerClient) {
	ui.speakerClient = client
	ui.speakerController = controllers.NewSpeakerController(client)
	ui.eqPresetService = services.NewEqPresetService()
	ui.colorPresetService = services.NewColorPresetService()

	ui.homePage = pages.NewHomePage(
		ui.theme,
		ui.buttonTheme,
		ui.speakerController,
		ui.eqPresetService,
		ui.colorPresetService,
		func(err error) {
			ui.loadingPage.SetError(err)
			ui.loaded = false
		},
	)

	ui.speakerController.RegisterListener(ui.homePage)

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
					return ui.loadingPage.Layout(gtx)
				}
				return ui.homePage.Layout(gtx)
			})
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
