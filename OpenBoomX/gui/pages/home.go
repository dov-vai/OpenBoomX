package pages

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"obx/gui/components"
	"obx/gui/controllers"
	"obx/gui/routes"
	"obx/gui/services"
)

type HomePage struct {
	theme              *material.Theme
	buttonTheme        *material.Theme
	topBar             *components.TopBar
	snackbar           *components.Snackbar
	speakerController  *controllers.SpeakerController
	eqPresetService    *services.EqPresetService
	colorPresetService *services.ColorPresetService
	oluvPage           *OluvPage
	eqPage             *EqPage
	presetsPage        *PresetsPage
	lightsPage         *LightsPage
	miscPage           *MiscPage
	currentRoute       routes.AppRoute
}

func NewHomePage(
	theme *material.Theme,
	buttonTheme *material.Theme,
	speakerController *controllers.SpeakerController,
	eqPresetService *services.EqPresetService,
	colorPresetService *services.ColorPresetService,
	snackbar *components.Snackbar,
	onUnload func(err error),
) *HomePage {
	page := &HomePage{}
	page.theme = theme
	page.buttonTheme = buttonTheme
	page.speakerController = speakerController
	page.eqPresetService = eqPresetService
	page.colorPresetService = colorPresetService
	page.snackbar = snackbar
	page.currentRoute = routes.Oluv

	page.topBar = components.CreateTopBar(page.theme, page.buttonTheme, func(route routes.AppRoute) {
		page.currentRoute = route
	})

	page.oluvPage = NewOluvPage(page.buttonTheme, page.speakerController)
	page.eqPage = NewEqPage(page.theme, page.buttonTheme, page.eqPresetService, page.speakerController, page.snackbar)
	page.presetsPage = NewPresetsPage(page.buttonTheme, page.eqPresetService, page.snackbar)
	page.lightsPage = NewLightsPage(page.theme, page.buttonTheme, page.speakerController, page.colorPresetService, page.snackbar)
	page.miscPage = NewMiscPage(page.theme, page.buttonTheme, page.speakerController, page.speakerController.GetFirmwareName())

	go page.speakerController.UpdateBattery(func(value int, err error) {
		page.topBar.UpdateBatteryLevel(value)
		if err != nil {
			onUnload(err)
		}
	})

	return page
}

func (h *HomePage) Update(gtx layout.Context) {
	h.miscPage.Update(gtx)
}

func (h *HomePage) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return h.topBar.Layout(gtx)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			switch h.currentRoute {
			case routes.Oluv:
				return h.oluvPage.Layout(gtx)
			case routes.Eq:
				return h.eqPage.Layout(gtx)
			case routes.EqPresets:
				return h.presetsPage.Layout(gtx)
			case routes.Lights:
				return h.lightsPage.Layout(gtx)
			case routes.Misc:
				return h.miscPage.Layout(gtx)
			default:
				return layout.Dimensions{}
			}
		}),
	)
}
