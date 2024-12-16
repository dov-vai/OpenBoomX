package components

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type StatusBar struct {
	BatteryLevel int
}

func CreateStatusBar() *StatusBar {
	return &StatusBar{}
}

func (sb *StatusBar) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.H6(th, fmt.Sprintf("Battery: %d%%", sb.BatteryLevel)).Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(32)}.Layout(gtx)
		}),
	)
}
