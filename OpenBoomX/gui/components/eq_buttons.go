package components

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"obx/protocol"
	"obx/utils"
)

type EqButtons struct {
	Buttons       []EQButton
	OnModeClicked func(mode string)
}

type EQButton struct {
	Mode      string
	Clickable widget.Clickable
}

func CreateEQButtons(onModeClicked func(mode string)) *EqButtons {
	buttons := make([]EQButton, 0, len(protocol.EQModes))
	for _, mode := range utils.SortedKeysByValue(protocol.EQModes) {
		buttons = append(buttons, EQButton{Mode: mode})
	}
	return &EqButtons{Buttons: buttons, OnModeClicked: onModeClicked}
}

func (eq *EqButtons) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	var buttons []layout.FlexChild

	for i := range eq.Buttons {
		btn := &eq.Buttons[i]
		if btn.Clickable.Clicked(gtx) {
			eq.OnModeClicked(btn.Mode)
		}

		btnLayout := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Button(th, &btn.Clickable, btn.Mode).Layout(gtx)
		})

		spacerLayout := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(8)}.Layout(gtx)
		})

		buttons = append(buttons, btnLayout, spacerLayout)
	}

	return layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceEvenly,
	}.Layout(gtx, buttons...)
}
