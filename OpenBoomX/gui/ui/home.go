package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
)
import "obx/gui/routes"

func (ui *UI) homeLayout(gtx layout.Context) layout.Dimensions {
	var children []layout.FlexChild

	children = append(children,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ui.StatusBar.Layout(ui.Theme, gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ui.NavigationBar.Layout(ui.Theme, gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(32)}.Layout(gtx)
		}),
	)

	// TODO: better to have separate pages?
	switch ui.CurrentRoute {
	case routes.Oluv:
		children = append(children,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.EqButtons.Layout(ui.Theme, gtx)
			}),
		)
	case routes.Eq:
		children = append(children,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.EqSaveButton.Layout(ui.Theme, gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.EqSlider.Layout(ui.Theme, gtx)
			}),
		)
	case routes.Lights:
		children = append(children,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.LightPicker.Layout(ui.Theme, gtx)
			}),
		)
	case routes.Misc:
		children = append(children,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.BeepSlider.Layout(ui.Theme, gtx)
			}),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.PairingButtons.Layout(ui.Theme, gtx)
			}),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.ShutdownSlider.Layout(ui.Theme, gtx)
			}),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.OffButton.Layout(ui.Theme, gtx)
			}),
		)
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
}
