package pages

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"obx/gui/components"
	"obx/gui/controllers"
	"obx/protocol"
	"obx/utils"
)

type OluvPage struct {
	buttonTheme       *material.Theme
	eqButtons         *components.EqButtons
	speakerController *controllers.SpeakerController
}

func NewOluvPage(buttonTheme *material.Theme, speakerController *controllers.SpeakerController) *OluvPage {
	page := &OluvPage{}
	page.buttonTheme = buttonTheme
	page.eqButtons = components.CreateEQButtons(
		utils.SortedKeysByValue(protocol.EQModes),
		speakerController.OnModeClicked,
	)
	return page
}

func (o *OluvPage) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return o.eqButtons.Layout(o.buttonTheme, gtx)
		}),
	)
}
