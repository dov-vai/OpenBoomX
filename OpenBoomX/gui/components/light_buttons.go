package components

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type LightButtons struct {
	defaultButton   widget.Clickable
	offButton       widget.Clickable
	OnActionClicked func(action string)
}

func CreateLightButtons(onActionClicked func(action string)) *LightButtons {
	return &LightButtons{
		OnActionClicked: onActionClicked,
	}
}

func (lb *LightButtons) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	if lb.offButton.Clicked(gtx) {
		// TODO: fix hardcodes everywhere..
		lb.OnActionClicked("off")
	}
	if lb.defaultButton.Clicked(gtx) {
		lb.OnActionClicked("default")
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return material.Button(th, &lb.defaultButton, "Default").Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Spacer{Width: unit.Dp(8)}.Layout(gtx)
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return material.Button(th, &lb.offButton, "Off").Layout(gtx)
				}),
			)
		}),
	)
}
