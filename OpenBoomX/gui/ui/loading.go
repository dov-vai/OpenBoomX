package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

// TODO: refactor this as a component?
func (ui *UI) loadingLayout(gtx layout.Context) layout.Dimensions {
	if ui.RetryConnection.Clicked(gtx) {
		ui.Error = nil
		go ui.connectSpeaker()
	}

	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				text := "Loading..."
				if ui.Error != nil {
					text = ui.Error.Error()
				}
				label := material.H5(ui.Theme, text)
				return label.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: unit.Dp(8)}.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if ui.Error != nil {
					return layout.Dimensions{}
				}

				gtx.Constraints.Max.X = gtx.Dp(32)
				gtx.Constraints.Max.Y = gtx.Dp(32)
				return material.Loader(ui.Theme).Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if ui.Error == nil {
					return layout.Dimensions{}
				}
				return material.Button(ui.Theme, &ui.RetryConnection, "Retry").Layout(gtx)
			}),
		)
	})
}
