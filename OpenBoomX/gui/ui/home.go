package ui

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"obx/gui/routes"
	"obx/gui/theme"
)

func (ui *UI) constructTopBar() layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Bottom: 16}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
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
						return ui.statusBar.Layout(ui.theme, gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return ui.navigationBar.Layout(ui.buttonTheme, gtx)
					}),
				)
			})
		})
	})
}

func (ui *UI) homeLayout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		ui.constructTopBar(),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			switch ui.currentRoute {
			case routes.Oluv:
				return ui.oluvPage.Layout(gtx)
			case routes.Eq:
				return ui.eqPage.Layout(gtx)
			case routes.EqPresets:
				return ui.presetsPage.Layout(gtx)
			case routes.Lights:
				return ui.lightsPage.Layout(gtx)
			case routes.Misc:
				return ui.miscPage.Layout(gtx)
			default:
				return layout.Dimensions{}
			}
		}),
	)
}
