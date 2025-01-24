package pages

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"obx/gui/components"
	"obx/gui/controllers"
	"obx/protocol"
	"obx/utils"
)

type MiscPage struct {
	theme             *material.Theme
	buttonTheme       *material.Theme
	beepSlider        *components.StepSlider
	videoModeButtons  *components.VideoModeButtons
	shutdownSlider    *components.StepSlider
	offButton         *components.OffButton
	speakerController *controllers.SpeakerController
	firmwareName      widget.Editor
}

func NewMiscPage(
	theme *material.Theme,
	buttonTheme *material.Theme,
	speakerController *controllers.SpeakerController,
	firmwareName string,
) *MiscPage {
	page := &MiscPage{}
	page.theme = theme
	page.buttonTheme = buttonTheme
	page.speakerController = speakerController

	page.firmwareName.ReadOnly = true
	page.firmwareName.SingleLine = true
	page.firmwareName.SetText(firmwareName)

	page.beepSlider = components.CreateBeepSlider(5, "Beep Volume", utils.SortedKeysByValueInt(protocol.BeepVolumes), page.speakerController.OnBeepStepChanged)
	page.offButton = components.CreateOffButton(page.speakerController.OnOffButtonClicked)
	page.shutdownSlider = components.CreateBeepSlider(7, "Shutdown Timeout", utils.SortedKeysByValue(protocol.ShutdownTimeouts), page.speakerController.OnShutdownStepChanged)
	page.videoModeButtons = components.CreateVideoModeButtons(page.speakerController.OnVideoModeEnabled, page.speakerController.OnVideoModeDisabled)
	return page
}

func (m *MiscPage) Update(gtx layout.Context) {
	m.beepSlider.Update(gtx)
	m.shutdownSlider.Update(gtx)
}

func (m *MiscPage) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return m.beepSlider.Layout(m.theme, gtx)
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return m.videoModeButtons.Layout(m.buttonTheme, gtx)
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return m.shutdownSlider.Layout(m.theme, gtx)
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return m.offButton.Layout(m.buttonTheme, gtx)
		}),

		layout.Rigid(layout.Spacer{Height: 8}.Layout),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.H6(m.theme, "Firmware:").Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.Editor(m.theme, &m.firmwareName, "").Layout(gtx)
				}),
			)
		}),
	)
}
