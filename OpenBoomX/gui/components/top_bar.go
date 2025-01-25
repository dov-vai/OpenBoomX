package components

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"obx/gui/routes"
	"obx/gui/theme"
)

type TopBar struct {
	theme         *material.Theme
	buttonTheme   *material.Theme
	navigationBar *NavigationBar
	statusBar     *StatusBar
}

func CreateTopBar(theme *material.Theme, buttonTheme *material.Theme, onRouteSelected func(route routes.AppRoute)) *TopBar {
	bar := &TopBar{}
	bar.theme = theme
	bar.buttonTheme = buttonTheme
	bar.navigationBar = CreateNavigationBar(onRouteSelected)
	bar.statusBar = CreateStatusBar()
	return bar
}

func (t *TopBar) Layout(gtx layout.Context) layout.Dimensions {
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
					return t.statusBar.Layout(t.theme, gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return t.navigationBar.Layout(t.buttonTheme, gtx)
				}),
			)
		})
	})
}

func (t *TopBar) UpdateBatteryLevel(value int) {
	t.statusBar.BatteryLevel = value
}
