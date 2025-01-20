package components

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type VideoModeButtons struct {
	clickableOn    widget.Clickable
	clickableOff   widget.Clickable
	OnModeEnabled  func()
	OnModeDisabled func()
}

func CreateVideoModeButtons(onModeEnabled func(), onModeDisabled func()) *VideoModeButtons {
	return &VideoModeButtons{
		OnModeEnabled:  onModeEnabled,
		OnModeDisabled: onModeDisabled,
	}
}

func (pb *VideoModeButtons) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	if pb.clickableOn.Clicked(gtx) {
		pb.OnModeEnabled()
	}
	if pb.clickableOff.Clicked(gtx) {
		pb.OnModeDisabled()
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return material.Button(th, &pb.clickableOn, "Video Mode On").Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Spacer{Width: unit.Dp(8)}.Layout(gtx)
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return material.Button(th, &pb.clickableOff, "Video Mode Off").Layout(gtx)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(8)}.Layout(gtx)
		}),
	)
}
