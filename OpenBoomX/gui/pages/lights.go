package pages

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/widget/material"
	"image/color"
	"log"
	"obx/gui/components"
	"obx/gui/controllers"
	"obx/gui/services"
)

type LightsPage struct {
	theme              *material.Theme
	buttonTheme        *material.Theme
	lightButtons       *components.LightButtons
	colorButtons       *components.ColorButtons
	colorEditButtons   *components.ColorEditButtons
	lightPicker        *components.LightPicker
	colorWheel         *components.ColorWheel
	gradientSelector   *components.GradientSelector
	snackbar           *components.Snackbar
	speakerController  *controllers.SpeakerController
	colorPresetService *services.ColorPresetService
	colorRemoveMode    bool
}

func NewLightsPage(
	theme *material.Theme,
	buttonTheme *material.Theme,
	speakerController *controllers.SpeakerController,
	colorPresetService *services.ColorPresetService,
	snackbar *components.Snackbar,
) *LightsPage {
	page := &LightsPage{}
	page.theme = theme
	page.buttonTheme = buttonTheme
	page.speakerController = speakerController
	page.colorPresetService = colorPresetService
	page.snackbar = snackbar

	page.lightButtons = components.CreateLightButtons(page.speakerController.OnLightDefaultClicked, page.speakerController.OnLightOffClicked)
	page.gradientSelector = components.CreateGradientSelector(func(color color.NRGBA) {
		page.speakerController.OnColorChanged(color, true)
	})
	page.lightPicker = components.CreateLightPicker(func(color color.NRGBA, solid bool) {
		page.speakerController.OnColorChangedDebounced(color, solid)
		page.gradientSelector.OnColorSelected(color)
	})
	page.colorButtons = components.CreateColorButtons(page.colorPresetService.ListColors(), 10, func(color color.NRGBA) {
		if page.colorRemoveMode {
			err := page.colorPresetService.DeleteColor(color)
			if err != nil {
				log.Printf("Error deleting color: %v", err)
			}
			return
		}
		page.lightPicker.SetColor(color)
	})
	page.colorPresetService.RegisterListener(page.colorButtons)
	page.colorEditButtons = components.CreateColorEditButtons(
		func() {
			err := page.colorPresetService.AddColor(page.lightPicker.GetColor())
			if err != nil {
				log.Printf("Error adding color: %v", err)
				page.snackbar.ShowMessage(fmt.Sprintf("Error adding color: %v", err))
			}
		},
		func(on bool) {
			page.colorRemoveMode = on
		})

	page.colorWheel = components.CreateColorWheel(page.lightPicker.SetColor)
	return page
}

func (l *LightsPage) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return l.lightButtons.Layout(l.buttonTheme, gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceEvenly, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return l.colorButtons.Layout(l.theme, gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return l.colorEditButtons.Layout(l.buttonTheme, gtx)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceEvenly}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return l.lightPicker.Layout(l.theme, gtx)
				}),
				layout.Rigid(layout.Spacer{Width: 16}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return l.colorWheel.Layout(gtx, float32(gtx.Constraints.Max.X)/8)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return l.gradientSelector.Layout(l.buttonTheme, gtx)
		}),
	)
}
