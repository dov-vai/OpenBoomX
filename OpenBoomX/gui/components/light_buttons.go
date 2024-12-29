package components

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type LightButtons struct {
	defaultButton    widget.Clickable
	offButton        widget.Clickable
	OnDefaultClicked func()
	OnOffClicked     func()
}

func CreateLightButtons(onDefaultClicked func(), onOffClicked func()) *LightButtons {
	return &LightButtons{
		OnDefaultClicked: onDefaultClicked,
		OnOffClicked:     onOffClicked,
	}
}

func (lb *LightButtons) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	if lb.offButton.Clicked(gtx) {
		lb.OnOffClicked()
	}
	if lb.defaultButton.Clicked(gtx) {
		lb.OnDefaultClicked()
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
