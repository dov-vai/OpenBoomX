package components

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"obx/gui/routes"
	"obx/gui/theme"
)

type RouteButtonData struct {
	label string
	route routes.AppRoute
}

var buttons = []RouteButtonData{
	{label: "Oluv", route: routes.Oluv},
	{label: "EQ", route: routes.Eq},
	{label: "Profiles", route: routes.EqProfiles},
	{label: "Lights", route: routes.Lights},
	{label: "Misc", route: routes.Misc},
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
		route := btnData.route
		label := btnData.label

		if clickable.Clicked(gtx) {
			nb.OnRouteSelected(route)
		}

		routeButtons[i] = layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return material.Button(&navTheme, clickable, label).Layout(gtx)
		})
	}

	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceEvenly}.Layout(gtx, routeButtons...)
}
