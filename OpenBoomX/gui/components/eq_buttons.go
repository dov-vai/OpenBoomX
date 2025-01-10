package components

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type EqButtons struct {
	Buttons       []EQButton
	OnModeClicked func(mode string)
}

type EQButton struct {
	mode      string
	clickable widget.Clickable
}

func CreateEQButtons(modes []string, onModeClicked func(mode string)) *EqButtons {
	buttons := make([]EQButton, 0, len(modes))
	for _, mode := range modes {
		buttons = append(buttons, EQButton{mode: mode})
	}
	return &EqButtons{Buttons: buttons, OnModeClicked: onModeClicked}
}

func (eq *EqButtons) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	var buttons []layout.FlexChild

	caser := cases.Title(language.English)

	for i := range eq.Buttons {
		btn := &eq.Buttons[i]
		if btn.clickable.Clicked(gtx) {
			eq.OnModeClicked(btn.mode)
		}

		btnLayout := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Button(th, &btn.clickable, caser.String(btn.mode)).Layout(gtx)
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
