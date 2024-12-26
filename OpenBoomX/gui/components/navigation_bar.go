package components

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"obx/gui/routes"
	"obx/gui/theme"
)

type RouteButtonData struct {
	Label string
	Route routes.AppRoute
}

var buttons = []RouteButtonData{
	{Label: "Oluv", Route: routes.Oluv},
	{Label: "EQ", Route: routes.Eq},
	{Label: "Profiles", Route: routes.EqProfiles},
	{Label: "Lights", Route: routes.Lights},
	{Label: "Misc", Route: routes.Misc},
}

type NavigationBar struct {
	OnRouteSelected func(route routes.AppRoute)
	clickables      []*widget.Clickable
}

func CreateNavigationBar(onRouteSelected func(route routes.AppRoute)) *NavigationBar {
	clickables := make([]*widget.Clickable, len(buttons))
	for i := range clickables {
		clickables[i] = new(widget.Clickable)
	}
	return &NavigationBar{
		OnRouteSelected: onRouteSelected,
		clickables:      clickables,
	}
}

func (nb *NavigationBar) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	navTheme := *th
	navTheme.ContrastBg = theme.CrustColor

	routeButtons := make([]layout.FlexChild, len(buttons))

	for i, btnData := range buttons {
		clickable := nb.clickables[i]
		route := btnData.Route
		label := btnData.Label

		if clickable.Clicked(gtx) {
			nb.OnRouteSelected(route)
		}

		routeButtons[i] = layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return material.Button(&navTheme, clickable, label).Layout(gtx)
		})
	}

	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceEvenly}.Layout(gtx, routeButtons...)
}
