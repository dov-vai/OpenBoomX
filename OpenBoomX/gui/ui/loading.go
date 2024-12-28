package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

// TODO: refactor this as a component?
func (ui *UI) loadingLayout(gtx layout.Context) layout.Dimensions {
	if ui.retryConnection.Clicked(gtx) {
		ui.appError = nil
		go ui.connectSpeaker()
	}

	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				text := "Loading..."
				if ui.appError != nil {
					text = ui.appError.Error()
				}
				label := material.H5(ui.theme, text)
				return label.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: unit.Dp(8)}.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if ui.appError != nil {
					return layout.Dimensions{}
				}

				gtx.Constraints.Max.X = gtx.Dp(32)
				gtx.Constraints.Max.Y = gtx.Dp(32)
				return material.Loader(ui.theme).Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if ui.appError == nil {
					return layout.Dimensions{}
				}
				return material.Button(ui.theme, &ui.retryConnection, "Retry").Layout(gtx)
			}),
		)
	})
}
