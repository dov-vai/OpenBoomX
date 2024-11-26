package ui

import "gioui.org/layout"

func (ui *UI) homeLayout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ui.EqButtons.Layout(ui.Theme, gtx)
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ui.LightPicker.Layout(ui.Theme, gtx)
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ui.BeepSlider.Layout(ui.Theme, gtx)
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ui.PairingButtons.Layout(ui.Theme, gtx)
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ui.OffButton.Layout(ui.Theme, gtx)
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ui.ShutdownSlider.Layout(ui.Theme, gtx)
		}),
	)
}
