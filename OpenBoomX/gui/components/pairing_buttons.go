package components

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type PairingButtons struct {
	clickableOn  widget.Clickable
	clickableOff widget.Clickable
	OnPairingOn  func()
	OnPairingOff func()
}

func CreatePairingButtons(onPairingOn func(), onPairingOff func()) *PairingButtons {
	return &PairingButtons{
		OnPairingOn:  onPairingOn,
		OnPairingOff: onPairingOff,
	}
}

func (pb *PairingButtons) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	if pb.clickableOn.Clicked(gtx) {
		pb.OnPairingOn()
	}
	if pb.clickableOff.Clicked(gtx) {
		pb.OnPairingOff()
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return material.Button(th, &pb.clickableOn, "Pairing On").Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Spacer{Width: unit.Dp(8)}.Layout(gtx)
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return material.Button(th, &pb.clickableOff, "Pairing Off").Layout(gtx)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(8)}.Layout(gtx)
		}),
	)
}
