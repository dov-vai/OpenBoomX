package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"obx/gui/theme"
)
import "obx/gui/routes"

func (ui *UI) constructTopBar() []layout.FlexChild {
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
					return ui.statusBar.Layout(ui.theme, gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return ui.navigationBar.Layout(ui.buttonTheme, gtx)
				}),
			)
		})
	})

	children = append(children,
		topBar,
		layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
	)

	return children
}

func (ui *UI) constructOluvPage() []layout.FlexChild {
	var children []layout.FlexChild

	children = append(children,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ui.eqButtons.Layout(ui.buttonTheme, gtx)
		}),
	)

	return children
}

func (ui *UI) constructEqPage() []layout.FlexChild {
	var children []layout.FlexChild

	children = append(children,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return ui.eqSaveButton.Layout(ui.buttonTheme, gtx)
				}),
				layout.Rigid(layout.Spacer{Width: 8}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return ui.eqResetButton.Layout(ui.buttonTheme, gtx)
				}),
			)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ui.eqSlider.Layout(ui.theme, gtx)
		}),
	)

	return children
}

func (ui *UI) constructProfilesPage() []layout.FlexChild {
	var children []layout.FlexChild

	children = append(children,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ui.presetButtons.Layout(ui.buttonTheme, gtx)
		}),
	)

	return children
}

func (ui *UI) constructLightsPage() []layout.FlexChild {
	var children []layout.FlexChild

	children = append(children,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ui.lightButtons.Layout(ui.buttonTheme, gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ui.lightPicker.Layout(ui.theme, gtx)
		}),
	)

	return children
}

func (ui *UI) constructMiscPage() []layout.FlexChild {
	var children []layout.FlexChild

	children = append(children,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ui.beepSlider.Layout(ui.theme, gtx)
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ui.pairingButtons.Layout(ui.buttonTheme, gtx)
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ui.shutdownSlider.Layout(ui.theme, gtx)
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ui.offButton.Layout(ui.buttonTheme, gtx)
		}),
	)

	return children
}

func (ui *UI) homeLayout(gtx layout.Context) layout.Dimensions {
	var children []layout.FlexChild

	children = append(children, ui.constructTopBar()...)

	switch ui.currentRoute {
	case routes.Oluv:
		children = append(children, ui.constructOluvPage()...)
	case routes.Eq:
		children = append(children, ui.constructEqPage()...)
	case routes.EqProfiles:
		children = append(children, ui.constructProfilesPage()...)
	case routes.Lights:
		children = append(children, ui.constructLightsPage()...)
	case routes.Misc:
		children = append(children, ui.constructMiscPage()...)
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
}
