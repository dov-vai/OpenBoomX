package pages

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/widget/material"
	"log"
	"obx/gui/components"
	"obx/gui/controllers"
	"obx/gui/services"
)

type EqPage struct {
	theme             *material.Theme
	buttonTheme       *material.Theme
	eqSaveButton      *components.EqSaveButton
	eqResetButton     *components.EqResetButton
	eqSlider          *components.EqSlider
	eqPresetService   *services.EqPresetService
	speakerController *controllers.SpeakerController
	snackbar          *components.Snackbar
}

func NewEqPage(
	theme *material.Theme,
	buttonTheme *material.Theme,
	eqPresetService *services.EqPresetService,
	speakerController *controllers.SpeakerController,
	snackbar *components.Snackbar,
) *EqPage {
	page := &EqPage{}
	page.theme = theme
	page.buttonTheme = buttonTheme
	page.eqPresetService = eqPresetService
	page.speakerController = speakerController
	page.snackbar = snackbar

	page.eqSlider = components.CreateEqSlider(page.speakerController.OnEqValuesChanged)
	// set currently active preset if it exists
	activePreset := page.eqPresetService.GetActivePreset()
	if activePreset != "" {
		eqValues, _ := page.eqPresetService.GetPresetValues(activePreset)
		err := page.eqSlider.SetSliderValues(eqValues)
		if err != nil {
			log.Println(err)
		}
	}
	page.eqPresetService.RegisterListener(page.eqSlider)

	page.eqResetButton = components.CreateEqResetButton(func() {
		err := page.eqSlider.ResetValues()
		if err != nil {
			log.Println(err)
		}
	})

	page.eqSaveButton = components.CreateEqSaveButton(func(title string) {
		err := page.eqPresetService.AddPreset(title, page.eqSlider.GetSliderValues())
		if err != nil {
			log.Println(err)
			page.snackbar.ShowMessage(fmt.Sprintf("Error adding preset: %v", err))
			return
		}
		page.snackbar.ShowMessage(fmt.Sprintf("Successfully added (or updated) preset: %s", title))
	})
	page.eqSaveButton.SetText(activePreset)
	page.eqPresetService.RegisterListener(page.eqSaveButton)
	return page
}

func (e *EqPage) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return e.eqSaveButton.Layout(e.buttonTheme, gtx)
				}),
				layout.Rigid(layout.Spacer{Width: 8}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return e.eqResetButton.Layout(e.buttonTheme, gtx)
				}),
			)
		}),
		layout.Rigid(layout.Spacer{Height: 8}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return e.eqSlider.Layout(e.theme, gtx)
		}),
	)
}
