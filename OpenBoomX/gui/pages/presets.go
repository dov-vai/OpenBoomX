package pages

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"obx/gui/components"
	"obx/gui/services"
)

type PresetsPage struct {
	buttonTheme     *material.Theme
	presetButtons   *components.PresetButtons
	eqPresetService *services.EqPresetService
	snackbar        *components.Snackbar
}

func NewPresetsPage(
	buttonTheme *material.Theme,
	eqPresetService *services.EqPresetService,
	snackbar *components.Snackbar,
) *PresetsPage {
	page := &PresetsPage{}
	page.buttonTheme = buttonTheme
	page.eqPresetService = eqPresetService
	page.snackbar = snackbar
	page.presetButtons = components.CreatePresetButtons(page.eqPresetService, page.snackbar)
	page.eqPresetService.RegisterListener(page.presetButtons)
	return page
}

func (p *PresetsPage) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return p.presetButtons.Layout(p.buttonTheme, gtx)
		}),
	)
}
