package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"obx/gui/theme"
)
import "obx/gui/routes"

func (ui *UI) homeLayout(gtx layout.Context) layout.Dimensions {
	var children []layout.FlexChild

	topBar := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		surfaceStyle := component.Surface(
			&material.Theme{
				Palette: material.Palette{
					Bg: theme.CrustColor,
				},
			})

		surfaceStyle.CornerRadius = 16

		return surfaceStyle.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return ui.StatusBar.Layout(ui.Theme, gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return ui.NavigationBar.Layout(ui.ButtonTheme, gtx)
				}),
			)
		})
	})

	children = append(children,
		topBar,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(16)}.Layout(gtx)
		}),
	)

	// TODO: better to have separate pages?
	switch ui.CurrentRoute {
	case routes.Oluv:
		children = append(children,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.EqButtons.Layout(ui.ButtonTheme, gtx)
			}),
		)
	case routes.Eq:
		children = append(children,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.EqSaveButton.Layout(ui.ButtonTheme, gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.EqSlider.Layout(ui.Theme, gtx)
			}),
		)
	case routes.EqProfiles:
		children = append(children,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.PresetButtons.Layout(ui.ButtonTheme, gtx)
			}),
		)
	case routes.Lights:
		children = append(children,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.LightButtons.Layout(ui.ButtonTheme, gtx)
			}),
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
				return ui.PairingButtons.Layout(ui.ButtonTheme, gtx)
			}),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.ShutdownSlider.Layout(ui.Theme, gtx)
			}),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.OffButton.Layout(ui.ButtonTheme, gtx)
			}),
		)
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
}
