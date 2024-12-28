package components

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type EqResetButton struct {
	clickable widget.Clickable
	OnClicked func()
}

func CreateEqResetButton(onClicked func()) *EqResetButton {
	return &EqResetButton{OnClicked: onClicked}
}

func (e *EqResetButton) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	if e.clickable.Clicked(gtx) {
		e.OnClicked()
	}

	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Button(th, &e.clickable, "Reset").Layout(gtx)
		}),
	)
}
