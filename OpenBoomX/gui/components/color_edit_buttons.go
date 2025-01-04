package components

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"obx/gui/theme"
)

type ColorEditButtons struct {
	addButton     widget.Clickable
	removeButton  widget.Clickable
	removeMode    bool
	OnAddColor    func()
	OnRemoveColor func(on bool)
}

func CreateColorEditButtons(onAddColor func(), onRemoveColor func(on bool)) *ColorEditButtons {
	return &ColorEditButtons{
		OnAddColor:    onAddColor,
		OnRemoveColor: onRemoveColor,
	}
}

func (ce *ColorEditButtons) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	buttonInset := unit.Dp(4)
	buttonPadding := unit.Dp(8)

	if ce.addButton.Clicked(gtx) {
		ce.OnAddColor()
	}

	if ce.removeButton.Clicked(gtx) {
		ce.removeMode = !ce.removeMode
		ce.OnRemoveColor(ce.removeMode)
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btnStyle := material.IconButton(th, &ce.addButton, theme.AddIcon, "Add color")
			btnStyle.Inset = layout.UniformInset(buttonInset)

			return layout.UniformInset(buttonPadding).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return btnStyle.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btnStyle := material.IconButton(th, &ce.removeButton, theme.DeleteIcon, "Remove mode")
			btnStyle.Inset = layout.UniformInset(buttonInset)
			if ce.removeMode {
				btnStyle.Background = theme.WarningColor
			}

			return layout.UniformInset(buttonPadding).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return btnStyle.Layout(gtx)
			})
		}),
	)
}
