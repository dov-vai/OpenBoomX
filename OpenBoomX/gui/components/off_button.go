package components

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type OffButton struct {
	clickable       widget.Clickable
	OnButtonClicked func()
}

func CreateOffButton(OnButtonClicked func()) *OffButton {
	return &OffButton{OnButtonClicked: OnButtonClicked}
}

func (btn *OffButton) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	if btn.clickable.Clicked(gtx) {
		btn.OnButtonClicked()
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Button(th, &btn.clickable, "Power Off").Layout(gtx)
		}))
}
