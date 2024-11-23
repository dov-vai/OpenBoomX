package components

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type OffButton struct {
	Clickable       widget.Clickable
	OnButtonClicked func()
}

func CreateOffButton(OnButtonClicked func()) *OffButton {
	return &OffButton{OnButtonClicked: OnButtonClicked}
}

func (btn *OffButton) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	if btn.Clickable.Clicked(gtx) {
		btn.OnButtonClicked()
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Button(th, &btn.Clickable, "Power Off").Layout(gtx)
		}))
}
